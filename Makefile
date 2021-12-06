BIN_DIR ?= ./bin

proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc -I=./api/core/v1alpha2 ./api/core/v1alpha2/core.proto --go_out=.
	protoc -I=./api/functionblock/v1alpha1 ./api/functionblock/v1alpha1/functionblock.proto --go_out=.
	protoc -I=./api/iou01/v1alpha1 ./api/iou01/v1alpha1/iou01.proto --go_out=.

.PHONY: all build clean test proto
