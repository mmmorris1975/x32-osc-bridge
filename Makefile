EXE := x32-osc-bridge
VER := $(shell git describe --tags)
#GOOS ?= $(shell go env GOOS)
#GOARCH ?= $(shell go env GOARCH)

.PHONY: clean dist-clean darwin-amd64 darwin-arm linux-amd64 linux-arm windows-amd64 windows-arm release test

$(EXE): go.mod *.go **/*.go
	go build -v -ldflags '-X main.Version=$(VER)'

clean:
	rm -f $(EXE) $(EXE)-*-*-*

dist-clean: clean
	rm -f go.sum

release: $(EXE) darwin-amd64 linux-amd64 linux-arm windows-amd64 windows-arm

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w -X main.Version=$(VER)' -o $(EXE)-$(VER)-$@
	upx -v $(EXE)-$(VER)-$@

darwin-arm:
	GOOS=darwin GOARCH=arm64 go build -ldflags '-s -w -X main.Version=$(VER)' -o $(EXE)-$(VER)-$@
	upx -v $(EXE)-$(VER)-$@

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags '-s -w -X main.Version=$(VER)' -o $(EXE)-$(VER)-$@
	upx -v $(EXE)-$(VER)-$@

linux-arm:
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags '-s -w -X main.Version=$(VER)' -o $(EXE)-$(VER)-$@
	upx -v $(EXE)-$(VER)-$@

windows-amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags '-s -w -X main.Version=$(VER)' -o $(EXE)-$(VER)-$@.exe
	upx -v $(EXE)-$(VER)-$@.exe

windows-arm:
	GOOS=windows GOARCH=arm go build -ldflags '-s -w -X main.Version=$(VER)' -o $(EXE)-$(VER)-$@.exe
	#upx -v $(EXE)-$(VER)-$@.exe