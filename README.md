About:

The codesweep program analyzes a collection of Go packages to identify
objects that are never used.

Install:

    $ go get github.com/mdempsky/codesweep

Usage:

    $ codesweep cmd/compile/...

Caveats:

Codesweep's primary use case is for cleaning up the Go toolchain's
machine-translated and non-idiomatic Go source code.

At this time, I suspect https://github.com/dominikh/go-unused with its
-exported flag does everything I want with fewer false positives.
