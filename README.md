# gi [![](https://github.com/shihanng/gi/workflows/main/badge.svg?branch=develop)](https://github.com/shihanng/gi/actions?query=workflow%3Amain) [![](https://github.com/shihanng/gi/workflows/release/badge.svg?branch=develop)](https://github.com/shihanng/gi/actions?query=workflow%3Arelease) [![GitHub](https://img.shields.io/github/license/shihanng/gi)](https://github.com/shihanng/gi/blob/develop/LICENSE) [![GitHub release (latest by date)](https://img.shields.io/github/v/release/shihanng/gi)](https://github.com/shihanng/gi/releases)

`gi` is a command line tool to help you create useful `.gitignore` files for your project.
It is inspired by [gitignore.io](https://www.gitignore.io/) and make use of
the large collection of useful [`.gitignore` templates](https://github.com/toptal/gitignore) of the web service.
This also means that `gi` supports the are four file types that gitignore.io recognizes.

# Install

## Binaries

The [release page contains](https://github.com/shihanng/gi/releases) binaries built
for various platforms. Download then extract the binary with `tar -xf`.
Place the binary in the `$PATH` e.g. `/usr/local/bin`.

## With `go get`

```
go get github.com/shihanng/gi
```

# Usage

Use the supported language as input arguments.

```
gi gen Go Elm
```

At the very first run the program will clone the templates repository <https://github.com/toptal/gitignore.git>
into `$XDG_CACHE_HOME/gi`.
This means that internet connection is not required after the first successful run.
