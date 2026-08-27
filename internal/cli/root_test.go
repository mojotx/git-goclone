package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCmd_RequiresArgs(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd(&bytes.Buffer{})
	cmd.SetArgs([]string{})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	require.Error(t, err)
}

func TestRootCmd_Help(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	cmd := NewRootCmd(&bytes.Buffer{})
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--help"})
	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), "git-goclone")
	assert.Contains(t, out.String(), "--depth")
	assert.Contains(t, out.String(), "--timeout")
}

func TestVersionFallback(t *testing.T) {
	t.Parallel()
	// In the test binary Main.Version is typically "(devel)" or empty, so we
	// expect the fallback.
	assert.NotEmpty(t, version())
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	quiet := newLogger(io.Discard, true)
	regular := newLogger(io.Discard, false)

	assert.NotNil(t, quiet)
	assert.NotNil(t, regular)
	quiet.Info().Msg("quiet logger should be usable")
	regular.Info().Msg("regular logger should be usable")
}

func TestAsExit(t *testing.T) {
	t.Parallel()

	var direct exitCodeError
	assert.True(t, asExit(exitCodeError(7), &direct))
	assert.Equal(t, exitCodeError(7), direct)

	var wrapped exitCodeError
	assert.True(t, asExit(fmt.Errorf("wrapped: %w", exitCodeError(9)), &wrapped))
	assert.Equal(t, exitCodeError(9), wrapped)

	var none exitCodeError
	assert.False(t, asExit(errors.New("plain error"), &none))
	assert.Equal(t, exitCodeError(0), none)
}

func TestRootCmd_Flags(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd(io.Discard)
	require.NoError(t, cmd.ParseFlags([]string{"--depth=2", "--timeout=45s", "--quiet"}))
	assert.Equal(t, "2", cmd.Flags().Lookup("depth").Value.String())
	assert.Equal(t, "45s", cmd.Flags().Lookup("timeout").Value.String())
	assert.Equal(t, "true", cmd.Flags().Lookup("quiet").Value.String())
	assert.Equal(t, "", cmd.Flags().Lookup("version").Shorthand)
}

func TestExitCodeError(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "exit code 3", exitCodeError(3).Error())
}

func TestExecute_NoArgsReturnsOne(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	code := execute(&stderr, []string{})
	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "Error:")
}
