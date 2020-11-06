# diff-cov

calculates test coverage of the changed go files

## Install
```shell script
$ go install github.com/nim4/cmd/diff-cov
```

## Usage
```shell script
$ go test ./... -coverprofile cover.out
$ diff-cov
```