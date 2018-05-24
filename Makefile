_BINARY="unity-cli"
_VERSION="v0.0.3"
_GOPATH="/Users/wilelb/Code/work/clarin/git/infrastructure2/golang"

_NAME_LINUX="${_BINARY}_linux_${_VERSION}"
_NAME_OSX="${_BINARY}_osx_${_VERSION}"

all: linux osx compress

linux:
	@echo "Building ${_NAME_LINUX}"
	@GOPATH=${_GOPATH} CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o "${_NAME_LINUX}" *.go
osx:
	@echo "Building ${_NAME_OSX}"
	@GOPATH=${_GOPATH} CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -a -installsuffix cgo -ldflags="-s -w" -o "${_NAME_OSX}" *.go

compress:
	@echo Compressing binaries
	@tar -pczvf "${_NAME_LINUX}.tar.gz" "${_NAME_LINUX}" && rm "${_NAME_LINUX}"
	@tar -pczvf "${_NAME_OSX}.tar.gz" "${_NAME_OSX}" && rm "${_NAME_OSX}"
