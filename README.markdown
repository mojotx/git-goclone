# git go-clone

![Coverage](https://img.shields.io/badge/Coverage-72.7%25-brightgreen)
[![CodeQL](https://github.com/mojotx/git-goclone/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/mojotx/git-goclone/actions/workflows/codeql-analysis.yml)
[![golangci-lint](https://github.com/mojotx/git-goclone/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/mojotx/git-goclone/actions/workflows/golangci-lint.yml)

This is fancy wrapper around `git clone` that preserves
directory structures.

For example, if you have some complex organization, and you want to
clone the repository `github.com/someorg/foo/bar/baz/project.git`, it
will create the directory structure based on the url path.

I used to use a shell script version of this, but rewrote it in Go
to make it more robust.

## Installation

To install this, use the following command:

```shell
    go install github.com/mojotx/git-goclone@latest
```

## Usage

Once the utility is installed, you can use it by giving one or more
Git URIs on a command-line.

For example:

```text
    $ git goclone https://github.com/mojotx/git-goclone.git
    processing https://github.com/mojotx/git-goclone.git...
    Cloning repo https://github.com/mojotx/git-goclone.git into mojotx/git-goclone...
    Enumerating objects: 28, done.
    Counting objects: 100% (28/28), done.
    Compressing objects: 100% (19/19), done.
    Total 28 (delta 7), reused 28 (delta 7), pack-reused 0
```
