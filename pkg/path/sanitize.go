package path

import (
	"os"
	"regexp"

	"github.com/mojotx/git-goclone/pkg/msg"
)

// Sanitize checks the path string to determine if the path ends with the substring ".git",
// and if so, removes that from the end of the string. The intent is that we do not want our
// Git repositories sitting in filesystem directories containing ".git".
//
// For example, if you clone this repository as https://github.com/mojotx/git-goclone.git, you
// would want the directory name containing the repository to be simply "git-goclone", not
// "git-goclone.git". This mimics the behavior of the Git command-line client.
func Sanitize(path string) string {

	pattern := ".git$"

	if m, err := regexp.MatchString(pattern, path); m {
		return path[0 : len(path)-4]
	} else if err != nil {
		msg.Err("weird error trying to match '%s' against regexp '%s': %s", path, pattern, err.Error())
		os.Exit(-1)
	}
	return path
}
