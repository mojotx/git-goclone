# git go-clone

![Coverage](https://img.shields.io/badge/Coverage-82.6%25-brightgreen)
[![CodeQL](https://github.com/mojotx/git-goclone/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/mojotx/git-goclone/actions/workflows/codeql-analysis.yml)
[![golangci-lint](https://github.com/mojotx/git-goclone/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/mojotx/git-goclone/actions/workflows/golangci-lint.yml)

This is a small Git wrapper that clones repositories into a directory tree
that mirrors the URL path — so cloning
`https://gitlab.com/company/department/project.git` lands in
`./company/department/project` rather than `./project`.

I use this to keep hundreds of repositories organised across GitHub, GitLab
and internal hosts without having to `mkdir -p` beforehand.

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

### Flags

| Flag          | Default | Description                             |
| ------------- | ------- | --------------------------------------- |
| `--depth`     | `1`     | Clone depth. `0` requests full history. |
| `--timeout`   | `5m`    | Per-URL clone timeout.                  |
| `-q, --quiet` | `false` | Suppress logs and git progress output.  |
| `--version`   |         | Print version and exit.                 |
| `-h, --help`  |         | Print help and exit.                    |

Passwords embedded in URLs are redacted from all log output.

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

