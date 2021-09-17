BIN_DIR ?= ./bin

proto:
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc -I=./api ./api/io4edge_core_api.proto --go_out=.

.PHONY: all build clean test proto
