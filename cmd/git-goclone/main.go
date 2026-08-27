// Command git-goclone clones one or more repositories into a directory
// structure that mirrors the URL path.
package main

import (
	"os"

	"github.com/mojotx/git-goclone/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
