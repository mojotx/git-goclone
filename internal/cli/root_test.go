package cli

import (
	"bytes"
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
