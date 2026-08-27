package clonepath

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		want      string
		wantError bool
	}{
		{name: "leading slash", path: "/example", want: "example"},
		{name: "leading backslash", path: "\\example", want: "example"},
		{name: "no leading separator", path: "example", want: "example"},
		{name: "trailing .git", path: "example.git", want: "example"},
		{name: "no .git suffix", path: "example", want: "example"},
		{name: "trailing 'git' is not '.git'", path: "examplegit", want: "examplegit"},
		{name: "empty string", path: "", wantError: true},
		{name: "root only", path: "/", wantError: true},
		{name: "nested forward slashes", path: "/org/team/project.git", want: "org/team/project"},
		{name: "path traversal", path: "../../etc/passwd", wantError: true},
		{name: "traversal disguised with .git", path: "../../../.ssh/authorized_keys.git", wantError: true},
		{name: "traversal in middle", path: "org/../../etc", wantError: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Sanitize(tt.path)
			if tt.wantError {
				assert.Error(t, err, "Sanitize(%q) should return error", tt.path)
				return
			}
			assert.NoError(t, err, "Sanitize(%q) should not return error", tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSanitize_RejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	actual := t.TempDir()
	require.NoError(t, os.Symlink(actual, filepath.Join(root, "escape")))

	origWD, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(root))
	t.Cleanup(func() { _ = os.Chdir(origWD) })

	_, err = Sanitize("escape/repo.git")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "outside")
}

func TestSanitize_ResolvesExistingPathWithinCWD(t *testing.T) {
	root := t.TempDir()
	origWD, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(root))
	t.Cleanup(func() { _ = os.Chdir(origWD) })

	got, err := Sanitize("org/team/proj.git")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join("org", "team", "proj"), got)
}

func TestValidateWithinCWD_RejectsRuntimeSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	origWD, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(root))
	t.Cleanup(func() { _ = os.Chdir(origWD) })

	require.NoError(t, os.MkdirAll(filepath.Join("org"), 0o755))
	require.NoError(t, ValidateWithinCWD(filepath.Join("org", "proj")))

	require.NoError(t, os.RemoveAll(filepath.Join("org")))
	require.NoError(t, os.Symlink(outside, filepath.Join("org")))

	err = ValidateWithinCWD(filepath.Join("org", "proj"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "outside working directory")
}

func TestValidateWithinCWD_RejectsUnsafeInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dest string
	}{
		{name: "empty", dest: ""},
		{name: "absolute", dest: "/tmp/escape"},
		{name: "traversal", dest: "../../etc/passwd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWithinCWD(tt.dest)
			require.Error(t, err)
		})
	}
}

func TestResolveCanonicalPathAndWithinDir(t *testing.T) {
	root := t.TempDir()
	resolved, err := resolveCanonicalPath(root, filepath.Join(root, "org", "repo"))
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "org", "repo"), resolved)

	_, err = resolveCanonicalPath(root, filepath.Join(root, "..", "outside"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "outside")

	assert.True(t, isWithinDir(root, filepath.Join(root, "org", "repo")))
	assert.False(t, isWithinDir(root, filepath.Join(root, "..", "outside")))
}

func TestTrimGitSuffix(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "example", trimGitSuffix("example.git"))
	assert.Equal(t, "example", trimGitSuffix("example"))
	assert.Equal(t, "examplegit", trimGitSuffix("examplegit"), "must not trim without leading dot")
	assert.Equal(t, "", trimGitSuffix(""))
}

func TestTrimLeadingSeparators(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "example", trimLeadingSeparators("/example"))
	assert.Equal(t, "example", trimLeadingSeparators("\\example"))
	assert.Equal(t, "example", trimLeadingSeparators("//\\/example"))
	assert.Equal(t, "example", trimLeadingSeparators("example"))
	assert.Equal(t, "", trimLeadingSeparators(""))
}
