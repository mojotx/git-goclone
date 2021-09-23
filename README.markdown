# git go-clone

This is fancy wrapper around `git clone` that preserves
directory structures.

For example, if you have some complex organization, and you want to
clone the repository github.com/someorg/foo/bar/baz/project.git, it
will create the directory structure based on the url path.

I used to use a shell script version of this, but rewrote it in Go
to make it more robust.


