package clonepath

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
