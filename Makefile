BIN_DIR ?= ./bin

proto:
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc -I=./proto ./proto/io4edge_base_function.proto --go_out=.

.PHONY: all build clean test proto
