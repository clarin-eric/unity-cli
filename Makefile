_BINARY="unity-cli"
_VERSION="v0.0.2"
_GOPATH="/Users/wilelb/Code/work/clarin/git/infrastructure2/golang"

_NAME_LINUX="${_BINARY}_linux_${_VERSION}"
_NAME_OSX="${_BINARY}_osx_${_VERSION}"

all: linux osx

linux:
	@echo "Building ${_NAME_LINUX}"
	@GOPATH=${_GOPATH} CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o "${_NAME_LINUX}" *.go

osx:
	@echo "Building ${_NAME_OSX}"
	@GOPATH=${_GOPATH} CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -a -installsuffix cgo -ldflags="-s -w" -o "${_NAME_OSX}" *.go
