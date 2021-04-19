VERSION := $(shell git describe --always --long --dirty)
DATE := $(shell date +%s)
LDFLAGS := -ldflags="-X github.com/ambientsound/visp/version.Version=${VERSION} -X github.com/ambientsound/visp/version.buildDate=${DATE}"

.PHONY: visp test linux-amd64 darwin-amd64 windows-amd64

visp:
	go build ${LDFLAGS} -o bin/visp cmd/visp/visp.go

visp-authproxy:
	go build ${LDFLAGS} -o bin/visp-authproxy cmd/visp-authproxy/main.go

test:
	go test ./...

test-coverage:
	go test -coverprofile=cover.out ./...

linux-amd64:
	GOOS=linux GOARCH=amd64 \
	go build ${LDFLAGS} -o bin/visp-linux-amd64 cmd/visp/visp.go

linux-arm64:
	GOOS=linux GOARCH=arm64 \
	go build ${LDFLAGS} -o bin/visp-linux-arm64 cmd/visp/visp.go

linux-arm:
	GOOS=linux GOARCH=arm \
	go build ${LDFLAGS} -o bin/visp-linux-arm cmd/visp/visp.go

darwin-amd64:
	GOOS=darwin GOARCH=amd64 \
	go build ${LDFLAGS} -o bin/visp-darwin-amd64 cmd/visp/visp.go

darwin-arm64:
	GOOS=darwin GOARCH=arm64 \
	go build ${LDFLAGS} -o bin/visp-darwin-arm64 cmd/visp/visp.go

windows-amd64.exe:
	GOOS=windows GOARCH=amd64 \
	go build ${LDFLAGS} -o bin/visp-windows-amd64.exe cmd/visp/visp.go
