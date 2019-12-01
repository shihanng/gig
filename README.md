# gi

`gi` is a command line tool to help you create useful `.gitignore` files for your project.
It is inspired by [gitignore.io](https://www.gitignore.io/) and make use of
the large collection of useful [`.gitignore` templates](https://github.com/toptal/gitignore) of the web service.
This also means that `gi` supports the are four file types that gitignore.io recognizes.

# Install

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
