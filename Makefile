

NAME = io4edge-cli
BIN_DIR ?= bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/io4edge-client-go/internal/version.Version=$(VERSION)"

all: test build-io4edge-cli

build: build-io4edge-cli

build-io4edge-cli:
	go mod tidy
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME} cmd/cli/main.go

test:
	go test ./...

clean:
	rm -f ${BIN_DIR}/${NAME}

proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc -I=./api/io4edge_core_api/v1alpha2 ./api/io4edge_core_api/v1alpha2/io4edge_core_api.proto --go_out=.
	protoc -I=./api/functionblock/v1alpha1 ./api/functionblock/v1alpha1/functionblock.proto --go_out=.
	protoc -I=./api/analogInTypeA/v1alpha1 ./api/analogInTypeA/v1alpha1/analogInTypeA.proto --go_out=.
	protoc -I=./api/binaryIoTypeA/v1alpha1 ./api/binaryIoTypeA/v1alpha1/binaryIoTypeA.proto --go_out=.

.PHONY: all clean proto build build-io4edge-cli
