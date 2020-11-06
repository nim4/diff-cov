# diff-cov

calculates test coverage of the changed go files

## Install
```shell script
$ go get github.com/nim4/diff-cov/cmd/diff-cov
```

## Usage
```shell script
$ go test ./... -coverprofile cover.out
$ diff-cov
```