

NAME = io4edge-cli
BIN_DIR ?= bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -tags 'netgo osusergo static_build' -ldflags "-X github.com/ci4rail/io4edge-client-go/v2/internal/version.Version=$(VERSION)"

all: test build-io4edge-cli

build: build-io4edge-cli

build-io4edge-cli:
	go mod tidy
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME} cmd/cli/main.go
	GOOS=linux GOARCH=arm go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME}-arm cmd/cli/main.go
	GOOS=darwin GOARCH=arm64 go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME}-darwin-arm64 cmd/cli/main.go

test:
	go test ./...

clean:
	rm -f ${BIN_DIR}/${NAME}
	rm -rf ${BIN_DIR}/examples

.PHONY: all clean build build-io4edge-cli
