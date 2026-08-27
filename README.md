# git go-clone

[![CI](https://github.com/mojotx/git-goclone/actions/workflows/ci.yml/badge.svg)](https://github.com/mojotx/git-goclone/actions/workflows/ci.yml)
[![CodeQL](https://github.com/mojotx/git-goclone/actions/workflows/codeql.yml/badge.svg)](https://github.com/mojotx/git-goclone/actions/workflows/codeql.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mojotx/git-goclone.svg)](https://pkg.go.dev/github.com/mojotx/git-goclone)

This is a small Git wrapper that clones repositories into a directory tree
that mirrors the URL path — so cloning
`https://gitlab.com/company/department/project.git` lands in
`./company/department/project` rather than `./project`.

I have used this for years to keep all of my repositories organised across
GitHub, GitLab and internal hosts without having to `mkdir -p` beforehand.

## Installation

```shell
go install github.com/mojotx/git-goclone/cmd/git-goclone@latest
```

Because the binary is named `git-goclone`, Git will pick it up as a
subcommand and you can invoke it as `git goclone …`.

## Usage

```text
$ git goclone https://github.com/mojotx/git-goclone.git
2026-08-27T12:00:00Z INF processing url=https://github.com/mojotx/git-goclone.git
2026-08-27T12:00:00Z INF cloning depth=1 dest=mojotx/git-goclone url=https://github.com/mojotx/git-goclone.git
```

Multiple URLs may be given; each is cloned independently, and the process
exits with a non-zero status if any of them failed.

### Supported URL forms

| Form             | Example                                       |
| ---------------- | --------------------------------------------- |
| HTTPS            | `https://github.com/mojotx/git-goclone.git`   |
| HTTPS with creds | `https://user:token@github.com/org/repo.git`  |
| SSH (`ssh://`)   | `ssh://git@github.com/mojotx/git-goclone.git` |
| SSH (SCP-style)  | `git@github.com:mojotx/git-goclone.git`       |
| Git protocol     | `git://github.com/mojotx/git-goclone.git`     |

SSH authentication uses your running SSH agent (`$SSH_AUTH_SOCK`); the host
must already be present in `~/.ssh/known_hosts`. There is no
`--ssh-key` / `-i` flag yet — PRs welcome.

### Flags

| Flag          | Default | Description                                     |
| ------------- | ------- | ----------------------------------------------- |
| `--depth`     | `0`     | Clone depth. `0` = full history; `1` = shallow. |
| `--timeout`   | `5m`    | Per-URL clone timeout.                          |
| `-q, --quiet` | `false` | Suppress logs and git progress output.          |
| `--version`   |         | Print version and exit.                         |
| `-h, --help`  |         | Print help and exit.                            |

Passwords embedded in URLs are redacted from all log output.

## Security notes

The tool derives a checkout path from the Git URL and validates that the final destination stays inside the current working directory before cloning. It rejects obvious traversal attempts, rejects absolute paths, and resolves symlinks/canonical paths to block directory escape via link-based tricks.

This is intended to protect the local filesystem when cloning from untrusted or malformed repository URLs. The code also performs a final destination re-check immediately before the clone call to reduce the remaining TOCTOU window.

## Project layout

```text
cmd/git-goclone/    # binary entry point
internal/cli/       # cobra command, flags, exit codes, logging
internal/clone/     # clone orchestration; uses go-git PlainCloneContext
internal/clonepath/ # URL-path -> filesystem-path sanitizer
```

## Development

```shell
go test ./...
go vet ./...
go build ./cmd/git-goclone
```
