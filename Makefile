
clean:
	rm -Rfv bin
	mkdir bin

build: clean
	go build -o bin/lightroom2aftershot cmd/main.go

build-all: clean
	GOOS="linux"   GOARCH="amd64"       go build -o bin/lightroom2aftershot__linux-amd64 cmd/main.go
	GOOS="linux"   GOARCH="arm" GOARM=6 go build -o bin/lightroom2aftershot__linux-armv6 cmd/main.go
	GOOS="linux"   GOARCH="arm" GOARM=7 go build -o bin/lightroom2aftershot__linux-armv7 cmd/main.go
	GOOS="linux"   GOARCH="arm"         go build -o bin/lightroom2aftershot__linux-arm   cmd/main.go
	GOOS="darwin"  GOARCH="amd64"       go build -o bin/lightroom2aftershot__macos-amd64 cmd/main.go
	GOOS="windows" GOARCH="amd64"       go build -o bin/lightroom2aftershot__win-amd64 cmd/main.go
