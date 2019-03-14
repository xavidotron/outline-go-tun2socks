# outline-go-tun2socks

Go package for building and using `go-tun2socks` for macOS, iOS and Android.

## Prerequisites

- macOS host (iOS)
- Xcode (iOS)
- make
- Go >= 1.12
- A C compiler (e.g.: clang, gcc)
- gomobile (https://github.com/golang/go/wiki/Mobile)
- Other common utilities (e.g.: git)

## iOS Golang Runtime

We use a custom Golang runtime to build the iOS framework built of Go 1.12. The [patch](https://go-review.googlesource.com/c/go/+/159117) improves memory reporting to the OS. This should not be necessary after Go 1.13 is released (scheduled for August 2019).

```bash
# We assume that Go is installed in /usr/local; this may vary on your system.
cd /usr/local/
# Temporarily move the current Go version, so as not to clobber $PATH.
mv go go1.12
# Download the Go source.
git clone https://go.googlesource.com/go
cd go
git checkout release-branch.go1.12
# Apply the patch.
git fetch https://go.googlesource.com/go refs/changes/17/159117/5 && git cherry-pick FETCH_HEAD
# Update the version.
echo "go1.12-runtime-patch" > VERSION
# Build the runtime.
cd src
GOROOT_BOOTSTRAP=/usr/local/go1.12/ ./make.bash
# Verify that the installed binary matches the custom version (i.e. 'go version go1.12-dev-runtime darwin/amd64').
go version
```

After building the framework, you can delete the custom runtime and revert to Go 1.12.

## Build
```bash
go get -d ./...
./build_[ios|android].sh
```