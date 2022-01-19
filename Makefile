

NAME = io4edge-cli
BIN_DIR ?= bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -tags 'netgo osusergo static_build' -ldflags "-X github.com/ci4rail/io4edge-client-go/internal/version.Version=$(VERSION)"

all: test build-io4edge-cli

build: build-io4edge-cli

build-io4edge-cli:
	go mod tidy
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME} cmd/cli/main.go
	GOOS=linux GOARCH=arm go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME}-arm cmd/cli/main.go
	scp ${BIN_DIR}/${NAME}-arm 192.168.24.32:~/bin/io4edge-cli

test:
	go test ./...

clean:
	rm -f ${BIN_DIR}/${NAME}
	rm -rf ${BIN_DIR}/examples

examples:
	mkdir -p ${BIN_DIR}/examples
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/examples/binaryIoTypeA_configurationControl examples/binaryIoTypeA_configurationControl/main.go
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/examples/binaryIoTypeA_functionControl examples/binaryIoTypeA_functionControl/main.go

proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc -I=./api/io4edge_core_api/v1alpha2 ./api/io4edge_core_api/v1alpha2/io4edge_core_api.proto --go_out=.
	protoc -I=./api/functionblock/v1alpha1 ./api/functionblock/v1alpha1/functionblock.proto --go_out=.
	protoc -I=./api/analogInTypeA/v1alpha1 ./api/analogInTypeA/v1alpha1/analogInTypeA.proto --go_out=.
	protoc -I=./api/binaryIoTypeA/v1alpha1 ./api/binaryIoTypeA/v1alpha1/binaryIoTypeA.proto --go_out=.

.PHONY: all clean proto build build-io4edge-cli examples
