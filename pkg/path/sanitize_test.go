package path

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Remove leading slash",
			path: "/example",
			want: "example",
		},
		{
			name: "Remove leading backslash",
			path: "\\example",
			want: "example",
		},
		{
			name: "No leading separator",
			path: "example",
			want: "example",
		},
		{
			name: "Remove .git suffix",
			path: "example.git",
			want: "example",
		},
		{
			name: "No .git suffix",
			path: "example",
			want: "example",
		},
		{
			name: "Empty string",
			path: "",
			want: "",
		},
		{
			name: "forward slash multi",
			path: "/mojotx/git-goclone.git",
			want: "mojotx/git-goclone",
		},
		{
			name: "backslash multi",
			path: "\\mojotx\\git-goclone.git",
			want: "mojotx\\git-goclone",
		},
		{
			name: "mixed format",
			path: "\\/mojotx\\/git-goclone.git",
			want: "mojotx\\/git-goclone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sanitize(tt.path)
			assert.Equalf(t, got, tt.want, "Sanitize: expected '%s', got '%s'", tt.want, got)
		})
	}
}

func TestPreTrim(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Remove leading slash",
			path: "/example",
			want: "example",
		},
		{
			name: "Remove leading backslash",
			path: "\\example",
			want: "example",
		},
		{
			name: "No leading separator",
			path: "example",
			want: "example",
		},
		{
			name: "Empty string",
			path: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PreTrim(tt.path)
			assert.Equalf(t, got, tt.want, "PreTrim: expected '%s', got '%s'", tt.want, got)
		})
	}
}

func TestPostTrim(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Remove .git suffix",
			path: "example.git",
			want: "example",
		},
		{
			name: "No .git suffix",
			path: "example",
			want: "example",
		},
		{
			name: "Empty string",
			path: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PostTrim(tt.path)
			assert.Equalf(t, got, tt.want, "PostTrim: expected '%s', got '%s'", tt.want, got)
		})
	}
}
