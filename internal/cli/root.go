// Package cli wires the cobra command and process-level concerns (logging,
// exit codes) to the clone package.
package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/mojotx/git-goclone/internal/clone"
)

// Execute builds and runs the root command, returning the process exit code.
func Execute() int {
	cmd := NewRootCmd(os.Stderr)
	err := cmd.Execute()
	if err == nil {
		return 0
	}

	var exit exitCodeError
	if asExit(err, &exit) {
		// Per-URL failures were already logged; just propagate the code.
		return int(exit)
	}
	// Usage / flag errors: cobra silenced them, so surface here.
	fmt.Fprintln(os.Stderr, "Error:", err)
	return 1
}

// NewRootCmd returns the root cobra command. logOut receives structured logs
// and any progress output from git.
func NewRootCmd(logOut io.Writer) *cobra.Command {
	var (
		depth   int
		timeout time.Duration
		quiet   bool
	)

	cmd := &cobra.Command{
		Use:   "git-goclone <url> [url...]",
		Short: "Clone git repositories into a directory tree that mirrors the URL path",
		Long: `git-goclone wraps "git clone" so that a URL like
  https://gitlab.com/org/team/project.git
is cloned into ./org/team/project instead of ./project.

Multiple URLs may be given; each is cloned independently and errors on one
do not stop the others.`,
		Args:          cobra.MinimumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version(),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := newLogger(logOut, quiet)

			progress := logOut
			if quiet {
				progress = io.Discard
			}

			var errCount int
			for _, raw := range args {
				ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
				err := clone.URL(ctx, logger, raw, clone.Options{
					Depth:    depth,
					Progress: progress,
				})
				cancel()
				if err != nil {
					errCount++
				}
			}
			if errCount > 0 {
				return exitCodeError(1)
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&depth, "depth", 0, "clone depth; 0 for full history, e.g. 1 for a shallow clone")
	cmd.Flags().DurationVar(&timeout, "timeout", 5*time.Minute, "per-URL clone timeout")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress logs and git progress output")

	// Cobra assigns `-v` to --version by default; drop the shorthand so
	// users can bind `-v` to something more conventional later.
	cmd.InitDefaultVersionFlag()
	if f := cmd.Flags().Lookup("version"); f != nil {
		f.Shorthand = ""
	}

	return cmd
}

func newLogger(w io.Writer, quiet bool) zerolog.Logger {
	if quiet {
		return zerolog.New(io.Discard)
	}
	return zerolog.New(zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()
}

// version resolves a build-info version at runtime so `go install` users get
// the module version and locally-built binaries still work.
func version() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}

// exitCodeError lets RunE signal a specific exit code without printing an
// extra "Error:" line — Execute unwraps it silently.
type exitCodeError int

func (e exitCodeError) Error() string { return fmt.Sprintf("exit code %d", int(e)) }

func asExit(err error, target *exitCodeError) bool {
	for e := err; e != nil; {
		if ec, ok := e.(exitCodeError); ok {
			*target = ec
			return true
		}
		type unwrapper interface{ Unwrap() error }
		u, ok := e.(unwrapper)
		if !ok {
			return false
		}
		e = u.Unwrap()
	}
	return false
}
