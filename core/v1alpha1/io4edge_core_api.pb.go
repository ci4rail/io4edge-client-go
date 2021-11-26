//
//Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: io4edge_core_api.proto

package v1alpha1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CommandId int32

const (
	CommandId_IDENTIFY_HARDWARE               CommandId = 0
	CommandId_IDENTIFY_FIRMWARE               CommandId = 1
	CommandId_LOAD_FIRMWARE_CHUNK             CommandId = 2
	CommandId_PROGRAM_HARDWARE_IDENTIFICATION CommandId = 3
	CommandId_RESTART                         CommandId = 4
	CommandId_PROGRAM_DEVICE_IDENTIFICATION   CommandId = 5
)

// Enum value maps for CommandId.
var (
	CommandId_name = map[int32]string{
		0: "IDENTIFY_HARDWARE",
		1: "IDENTIFY_FIRMWARE",
		2: "LOAD_FIRMWARE_CHUNK",
		3: "PROGRAM_HARDWARE_IDENTIFICATION",
		4: "RESTART",
		5: "PROGRAM_DEVICE_IDENTIFICATION",
	}
	CommandId_value = map[string]int32{
		"IDENTIFY_HARDWARE":               0,
		"IDENTIFY_FIRMWARE":               1,
		"LOAD_FIRMWARE_CHUNK":             2,
		"PROGRAM_HARDWARE_IDENTIFICATION": 3,
		"RESTART":                         4,
		"PROGRAM_DEVICE_IDENTIFICATION":   5,
	}
)

func (x CommandId) Enum() *CommandId {
	p := new(CommandId)
	*p = x
	return p
}

func (x CommandId) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CommandId) Descriptor() protoreflect.EnumDescriptor {
	return file_io4edge_core_api_proto_enumTypes[0].Descriptor()
}

func (CommandId) Type() protoreflect.EnumType {
	return &file_io4edge_core_api_proto_enumTypes[0]
}

func (x CommandId) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CommandId.Descriptor instead.
func (CommandId) EnumDescriptor() ([]byte, []int) {
	return file_io4edge_core_api_proto_rawDescGZIP(), []int{0}
}

type Status int32

const (
	Status_OK                          Status = 0
	Status_UNKNOWN_COMMAND             Status = 1
	Status_ILLEGAL_PARAMETER           Status = 2
	Status_BAD_CHUNK_SEQ               Status = 3
	Status_BAD_CHUNK_SIZE              Status = 4
	Status_NOT_COMPATIBLE              Status = 5
	Status_INTERNAL_ERROR              Status = 6
	Status_PROGRAMMING_ERROR           Status = 7
	Status_NO_HW_INVENTORY             Status = 8
	Status_THIS_VERSION_FAILED_ALREADY Status = 9
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "OK",
		1: "UNKNOWN_COMMAND",
		2: "ILLEGAL_PARAMETER",
		3: "BAD_CHUNK_SEQ",
		4: "BAD_CHUNK_SIZE",
		5: "NOT_COMPATIBLE",
		6: "INTERNAL_ERROR",
		7: "PROGRAMMING_ERROR",
		8: "NO_HW_INVENTORY",
		9: "THIS_VERSION_FAILED_ALREADY",
	}
	Status_value = map[string]int32{
		"OK":                          0,
		"UNKNOWN_COMMAND":             1,
		"ILLEGAL_PARAMETER":           2,
		"BAD_CHUNK_SEQ":               3,
		"BAD_CHUNK_SIZE":              4,
		"NOT_COMPATIBLE":              5,
		"INTERNAL_ERROR":              6,
		"PROGRAMMING_ERROR":           7,
		"NO_HW_INVENTORY":             8,
		"THIS_VERSION_FAILED_ALREADY": 9,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_io4edge_core_api_proto_enumTypes[1].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_io4edge_core_api_proto_enumTypes[1]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_io4edge_core_api_proto_rawDescGZIP(), []int{1}
}

type SerialNumber struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hi uint64 `protobuf:"fixed64,1,opt,name=hi,proto3" json:"hi,omitempty"`
	Lo uint64 `protobuf:"fixed64,2,opt,name=lo,proto3" json:"lo,omitempty"`
}

func (x *SerialNumber) Reset() {
	*x = SerialNumber{}
	if protoimpl.UnsafeEnabled {
		mi := &file_io4edge_core_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SerialNumber) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SerialNumber) ProtoMessage() {}

func (x *SerialNumber) ProtoReflect() protoreflect.Message {
	mi := &file_io4edge_core_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SerialNumber.ProtoReflect.Descriptor instead.
func (*SerialNumber) Descriptor() ([]byte, []int) {
	return file_io4edge_core_api_proto_rawDescGZIP(), []int{0}
}

func (x *SerialNumber) GetHi() uint64 {
	if x != nil {
		return x.Hi
	}
	return 0
}

func (x *SerialNumber) GetLo() uint64 {
	if x != nil {
		return x.Lo
	}
	return 0
}

// LoadFirmware
// Client sends sequence of CmdLoadFirmwareChunk commands, with increasing
// chunk numbers. Clients defines chunk size.
// Server must acknowledge each chunk with Response.
// Last chunk has is_last_chunk set to True, so server knows that programming
// has finished
type LoadFirmwareChunkCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkNumber uint32 `protobuf:"varint,1,opt,name=chunk_number,json=chunkNumber,proto3" json:"chunk_number,omitempty"`
	IsLastChunk bool   `protobuf:"varint,2,opt,name=is_last_chunk,json=isLastChunk,proto3" json:"is_last_chunk,omitempty"`
	Data        []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *LoadFirmwareChunkCommand) Reset() {
	*x = LoadFirmwareChunkCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_io4edge_core_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadFirmwareChunkCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadFirmwareChunkCommand) ProtoMessage() {}

func (x *LoadFirmwareChunkCommand) ProtoReflect() protoreflect.Message {
	mi := &file_io4edge_core_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadFirmwareChunkCommand.ProtoReflect.Descriptor instead.
func (*LoadFirmwareChunkCommand) Descriptor() ([]byte, []int) {
	return file_io4edge_core_api_proto_rawDescGZIP(), []int{1}
}

func (x *LoadFirmwareChunkCommand) GetChunkNumber() uint32 {
	if x != nil {
		return x.ChunkNumber
	}
	return 0
}

func (x *LoadFirmwareChunkCommand) GetIsLastChunk() bool {
	if x != nil {
		return x.IsLastChunk
	}
	return false
}

func (x *LoadFirmwareChunkCommand) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type ProgramHardwareIdentificationCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signature    string        `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	RootArticle  string        `protobuf:"bytes,2,opt,name=root_article,json=rootArticle,proto3" json:"root_article,omitempty"`
	MajorVersion uint32        `protobuf:"varint,3,opt,name=major_version,json=majorVersion,proto3" json:"major_version,omitempty"`
	SerialNumber *SerialNumber `protobuf:"bytes,4,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
}

func (x *ProgramHardwareIdentificationCommand) Reset() {
	*x = ProgramHardwareIdentificationCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_io4edge_core_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProgramHardwareIdentificationCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProgramHardwareIdentificationCommand) ProtoMessage() {}

func (x *ProgramHardwareIdentificationCommand) ProtoReflect() protoreflect.Message {
	mi := &file_io4edge_core_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProgramHardwareIdentificationCommand.ProtoReflect.Descriptor instead.
func (*ProgramHardwareIdentificationCommand) Descriptor() ([]byte, []int) {
	return file_io4edge_core_api_proto_rawDescGZIP(), []int{2}
}

func (x *ProgramHardwareIdentificationCommand) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

func (x *ProgramHardwareIdentificationCommand) GetRootArticle() string {
	if x != nil {
		return x.RootArticle
	}
	return ""
}

func (x *ProgramHardwareIdentificationCommand) GetMajorVersion() uint32 {
	if x != nil {
		return x.MajorVersion
	}
	return 0
}

func (x *ProgramHardwareIdentificationCommand) GetSerialNumber() *SerialNumber {
	if x != nil {
		return x.SerialNumber
	}
	return nil
}

type ProgramDeviceIdentificationCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InstanceName string `protobuf:"bytes,1,opt,name=instance_name,json=instanceName,proto3" json:"instance_name,omitempty"`
}

func (x *ProgramDeviceIdentificationCommand) Reset() {
	*x = ProgramDeviceIdentificationCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_io4edge_core_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProgramDeviceIdentificationCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProgramDeviceIdentificationCommand) ProtoMessage() {}

func (x *ProgramDeviceIdentificationCommand) ProtoReflect() protoreflect.Message {
	mi := &file_io4edge_core_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProgramDeviceIdentificationCommand.ProtoReflect.Descriptor instead.
func (*ProgramDeviceIdentificationCommand) Descriptor() ([]byte, []int) {
	return file_io4edge_core_api_proto_rawDescGZIP(), []int{3}
}

func (x *ProgramDeviceIdentificationCommand) GetInstanceName() string {
	if x != nil {
		return x.InstanceName
	}
	return ""
}

type IdentifyHardwareResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RootArticle  string        `protobuf:"bytes,1,opt,name=root_article,json=rootArticle,proto3" json:"root_article,omitempty"`
	MajorVersion uint32        `protobuf:"varint,2,opt,name=major_version,json=majorVersion,proto3" json:"major_version,omitempty"`
	SerialNumber *SerialNumber `protobuf:"bytes,3,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
}

func (x *IdentifyHardwareResponse) Reset() {
	*x = IdentifyHardwareResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_io4edge_core_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdentifyHardwareResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdentifyHardwareResponse) ProtoMessage() {}

func (x *IdentifyHardwareResponse) ProtoReflect() protoreflect.Message {
	mi := &file_io4edge_core_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdentifyHardwareResponse.ProtoReflect.Descriptor instead.
func (*IdentifyHardwareResponse) Descriptor() ([]byte, []int) {
	return file_io4edge_core_api_proto_rawDescGZIP(), []int{4}
}

func (x *IdentifyHardwareResponse) GetRootArticle() string {
	if x != nil {
		return x.RootArticle
	}
	return ""
}

func (x *IdentifyHardwareResponse) GetMajorVersion() uint32 {
	if x != nil {
		return x.MajorVersion
	}
	return 0
}

func (x *IdentifyHardwareResponse) GetSerialNumber() *SerialNumber {
	if x != nil {
		return x.SerialNumber
	}
	return nil
}

type IdentifyFirmwareResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *IdentifyFirmwareResponse) Reset() {
	*x = IdentifyFirmwareResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_io4edge_core_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdentifyFirmwareResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdentifyFirmwareResponse) ProtoMessage() {}

func (x *IdentifyFirmwareResponse) ProtoReflect() protoreflect.Message {
	mi := &file_io4edge_core_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdentifyFirmwareResponse.ProtoReflect.Descriptor instead.
func (*IdentifyFirmwareResponse) Descriptor() ([]byte, []int) {
	return file_io4edge_core_api_proto_rawDescGZIP(), []int{5}
}

func (x *IdentifyFirmwareResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *IdentifyFirmwareResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// The common messages
type CoreCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id CommandId `protobuf:"varint,1,opt,name=id,proto3,enum=io4edgeCoreApi.CommandId" json:"id,omitempty"`
	// Types that are assignable to Data:
	//	*CoreCommand_LoadFirmwareChunk
	//	*CoreCommand_ProgramHardwareIdentification
	//	*CoreCommand_ProgramDeviceIdentification
	Data isCoreCommand_Data `protobuf_oneof:"data"`
}

func (x *CoreCommand) Reset() {
	*x = CoreCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_io4edge_core_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoreCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoreCommand) ProtoMessage() {}

func (x *CoreCommand) ProtoReflect() protoreflect.Message {
	mi := &file_io4edge_core_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoreCommand.ProtoReflect.Descriptor instead.
func (*CoreCommand) Descriptor() ([]byte, []int) {
	return file_io4edge_core_api_proto_rawDescGZIP(), []int{6}
}

func (x *CoreCommand) GetId() CommandId {
	if x != nil {
		return x.Id
	}
	return CommandId_IDENTIFY_HARDWARE
}

func (m *CoreCommand) GetData() isCoreCommand_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *CoreCommand) GetLoadFirmwareChunk() *LoadFirmwareChunkCommand {
	if x, ok := x.GetData().(*CoreCommand_LoadFirmwareChunk); ok {
		return x.LoadFirmwareChunk
	}
	return nil
}

func (x *CoreCommand) GetProgramHardwareIdentification() *ProgramHardwareIdentificationCommand {
	if x, ok := x.GetData().(*CoreCommand_ProgramHardwareIdentification); ok {
		return x.ProgramHardwareIdentification
	}
	return nil
}

func (x *CoreCommand) GetProgramDeviceIdentification() *ProgramDeviceIdentificationCommand {
	if x, ok := x.GetData().(*CoreCommand_ProgramDeviceIdentification); ok {
		return x.ProgramDeviceIdentification
	}
	return nil
}

type isCoreCommand_Data interface {
	isCoreCommand_Data()
}

type CoreCommand_LoadFirmwareChunk struct {
	LoadFirmwareChunk *LoadFirmwareChunkCommand `protobuf:"bytes,2,opt,name=load_firmware_chunk,json=loadFirmwareChunk,proto3,oneof"`
}

type CoreCommand_ProgramHardwareIdentification struct {
	ProgramHardwareIdentification *ProgramHardwareIdentificationCommand `protobuf:"bytes,3,opt,name=program_hardware_identification,json=programHardwareIdentification,proto3,oneof"`
}

type CoreCommand_ProgramDeviceIdentification struct {
	ProgramDeviceIdentification *ProgramDeviceIdentificationCommand `protobuf:"bytes,4,opt,name=program_device_identification,json=programDeviceIdentification,proto3,oneof"`
}

func (*CoreCommand_LoadFirmwareChunk) isCoreCommand_Data() {}

func (*CoreCommand_ProgramHardwareIdentification) isCoreCommand_Data() {}

func (*CoreCommand_ProgramDeviceIdentification) isCoreCommand_Data() {}

type CoreResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            CommandId `protobuf:"varint,1,opt,name=id,proto3,enum=io4edgeCoreApi.CommandId" json:"id,omitempty"`
	Status        Status    `protobuf:"varint,2,opt,name=status,proto3,enum=io4edgeCoreApi.Status" json:"status,omitempty"`
	RestartingNow bool      `protobuf:"varint,3,opt,name=restarting_now,json=restartingNow,proto3" json:"restarting_now,omitempty"`
	// Types that are assignable to Data:
	//	*CoreResponse_IdentifyHardware
	//	*CoreResponse_IdentifyFirmware
	Data isCoreResponse_Data `protobuf_oneof:"data"`
}

func (x *CoreResponse) Reset() {
	*x = CoreResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_io4edge_core_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoreResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoreResponse) ProtoMessage() {}

func (x *CoreResponse) ProtoReflect() protoreflect.Message {
	mi := &file_io4edge_core_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoreResponse.ProtoReflect.Descriptor instead.
func (*CoreResponse) Descriptor() ([]byte, []int) {
	return file_io4edge_core_api_proto_rawDescGZIP(), []int{7}
}

func (x *CoreResponse) GetId() CommandId {
	if x != nil {
		return x.Id
	}
	return CommandId_IDENTIFY_HARDWARE
}

func (x *CoreResponse) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_OK
}

func (x *CoreResponse) GetRestartingNow() bool {
	if x != nil {
		return x.RestartingNow
	}
	return false
}

func (m *CoreResponse) GetData() isCoreResponse_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *CoreResponse) GetIdentifyHardware() *IdentifyHardwareResponse {
	if x, ok := x.GetData().(*CoreResponse_IdentifyHardware); ok {
		return x.IdentifyHardware
	}
	return nil
}

func (x *CoreResponse) GetIdentifyFirmware() *IdentifyFirmwareResponse {
	if x, ok := x.GetData().(*CoreResponse_IdentifyFirmware); ok {
		return x.IdentifyFirmware
	}
	return nil
}

type isCoreResponse_Data interface {
	isCoreResponse_Data()
}

type CoreResponse_IdentifyHardware struct {
	IdentifyHardware *IdentifyHardwareResponse `protobuf:"bytes,4,opt,name=identify_hardware,json=identifyHardware,proto3,oneof"`
}

type CoreResponse_IdentifyFirmware struct {
	IdentifyFirmware *IdentifyFirmwareResponse `protobuf:"bytes,5,opt,name=identify_firmware,json=identifyFirmware,proto3,oneof"`
}

func (*CoreResponse_IdentifyHardware) isCoreResponse_Data() {}

func (*CoreResponse_IdentifyFirmware) isCoreResponse_Data() {}

var File_io4edge_core_api_proto protoreflect.FileDescriptor

var file_io4edge_core_api_proto_rawDesc = []byte{
	0x0a, 0x16, 0x69, 0x6f, 0x34, 0x65, 0x64, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x69, 0x6f, 0x34, 0x65, 0x64, 0x67,
	0x65, 0x43, 0x6f, 0x72, 0x65, 0x41, 0x70, 0x69, 0x22, 0x2e, 0x0a, 0x0c, 0x53, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x68, 0x69, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x06, 0x52, 0x02, 0x68, 0x69, 0x12, 0x0e, 0x0a, 0x02, 0x6c, 0x6f, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x06, 0x52, 0x02, 0x6c, 0x6f, 0x22, 0x75, 0x0a, 0x18, 0x4c, 0x6f, 0x61, 0x64,
	0x46, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x63, 0x68, 0x75, 0x6e,
	0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x0d, 0x69, 0x73, 0x5f, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b,
	0x69, 0x73, 0x4c, 0x61, 0x73, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0xcf, 0x01, 0x0a, 0x24, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x48, 0x61, 0x72, 0x64, 0x77,
	0x61, 0x72, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x6f, 0x6f, 0x74, 0x5f, 0x61,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x6f,
	0x6f, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x6d, 0x61, 0x6a,
	0x6f, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0c, 0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x41,
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x69, 0x6f, 0x34, 0x65, 0x64, 0x67, 0x65, 0x43,
	0x6f, 0x72, 0x65, 0x41, 0x70, 0x69, 0x2e, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x52, 0x0c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x22, 0x49, 0x0a, 0x22, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x69, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xa5, 0x01, 0x0a,
	0x18, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x79, 0x48, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x6f, 0x6f,
	0x74, 0x5f, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x72, 0x6f, 0x6f, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x23, 0x0a, 0x0d,
	0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0c, 0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x41, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x69, 0x6f, 0x34, 0x65, 0x64,
	0x67, 0x65, 0x43, 0x6f, 0x72, 0x65, 0x41, 0x70, 0x69, 0x2e, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x22, 0x48, 0x0a, 0x18, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x79,
	0x46, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x96,
	0x03, 0x0a, 0x0b, 0x43, 0x6f, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x29,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x69, 0x6f, 0x34,
	0x65, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x72, 0x65, 0x41, 0x70, 0x69, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x49, 0x64, 0x52, 0x02, 0x69, 0x64, 0x12, 0x5a, 0x0a, 0x13, 0x6c, 0x6f, 0x61,
	0x64, 0x5f, 0x66, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x5f, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x69, 0x6f, 0x34, 0x65, 0x64, 0x67, 0x65,
	0x43, 0x6f, 0x72, 0x65, 0x41, 0x70, 0x69, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x72, 0x6d,
	0x77, 0x61, 0x72, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x48, 0x00, 0x52, 0x11, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65,
	0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x7e, 0x0a, 0x1f, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d,
	0x5f, 0x68, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34,
	0x2e, 0x69, 0x6f, 0x34, 0x65, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x72, 0x65, 0x41, 0x70, 0x69, 0x2e,
	0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x48, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x49,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x1d, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x48,
	0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x78, 0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d,
	0x5f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x69,
	0x6f, 0x34, 0x65, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x72, 0x65, 0x41, 0x70, 0x69, 0x2e, 0x50, 0x72,
	0x6f, 0x67, 0x72, 0x61, 0x6d, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x48, 0x00, 0x52, 0x1b, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xca, 0x02, 0x0a, 0x0c, 0x43, 0x6f, 0x72, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x69, 0x6f, 0x34, 0x65, 0x64, 0x67, 0x65, 0x43, 0x6f,
	0x72, 0x65, 0x41, 0x70, 0x69, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x49, 0x64, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x2e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x69, 0x6f, 0x34, 0x65, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x72,
	0x65, 0x41, 0x70, 0x69, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e,
	0x67, 0x5f, 0x6e, 0x6f, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x72, 0x65, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x4e, 0x6f, 0x77, 0x12, 0x57, 0x0a, 0x11, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x66, 0x79, 0x5f, 0x68, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x69, 0x6f, 0x34, 0x65, 0x64, 0x67, 0x65, 0x43,
	0x6f, 0x72, 0x65, 0x41, 0x70, 0x69, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x79, 0x48,
	0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48,
	0x00, 0x52, 0x10, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x79, 0x48, 0x61, 0x72, 0x64, 0x77,
	0x61, 0x72, 0x65, 0x12, 0x57, 0x0a, 0x11, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x79, 0x5f,
	0x66, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28,
	0x2e, 0x69, 0x6f, 0x34, 0x65, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x72, 0x65, 0x41, 0x70, 0x69, 0x2e,
	0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x79, 0x46, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x10, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x66, 0x79, 0x46, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x42, 0x06, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x2a, 0xa7, 0x01, 0x0a, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x49, 0x64, 0x12, 0x15, 0x0a, 0x11, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x49, 0x46, 0x59, 0x5f, 0x48,
	0x41, 0x52, 0x44, 0x57, 0x41, 0x52, 0x45, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x49, 0x44, 0x45,
	0x4e, 0x54, 0x49, 0x46, 0x59, 0x5f, 0x46, 0x49, 0x52, 0x4d, 0x57, 0x41, 0x52, 0x45, 0x10, 0x01,
	0x12, 0x17, 0x0a, 0x13, 0x4c, 0x4f, 0x41, 0x44, 0x5f, 0x46, 0x49, 0x52, 0x4d, 0x57, 0x41, 0x52,
	0x45, 0x5f, 0x43, 0x48, 0x55, 0x4e, 0x4b, 0x10, 0x02, 0x12, 0x23, 0x0a, 0x1f, 0x50, 0x52, 0x4f,
	0x47, 0x52, 0x41, 0x4d, 0x5f, 0x48, 0x41, 0x52, 0x44, 0x57, 0x41, 0x52, 0x45, 0x5f, 0x49, 0x44,
	0x45, 0x4e, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x03, 0x12, 0x0b,
	0x0a, 0x07, 0x52, 0x45, 0x53, 0x54, 0x41, 0x52, 0x54, 0x10, 0x04, 0x12, 0x21, 0x0a, 0x1d, 0x50,
	0x52, 0x4f, 0x47, 0x52, 0x41, 0x4d, 0x5f, 0x44, 0x45, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x49, 0x44,
	0x45, 0x4e, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x05, 0x2a, 0xd8,
	0x01, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10,
	0x00, 0x12, 0x13, 0x0a, 0x0f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x43, 0x4f, 0x4d,
	0x4d, 0x41, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x49, 0x4c, 0x4c, 0x45, 0x47, 0x41,
	0x4c, 0x5f, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x45, 0x54, 0x45, 0x52, 0x10, 0x02, 0x12, 0x11, 0x0a,
	0x0d, 0x42, 0x41, 0x44, 0x5f, 0x43, 0x48, 0x55, 0x4e, 0x4b, 0x5f, 0x53, 0x45, 0x51, 0x10, 0x03,
	0x12, 0x12, 0x0a, 0x0e, 0x42, 0x41, 0x44, 0x5f, 0x43, 0x48, 0x55, 0x4e, 0x4b, 0x5f, 0x53, 0x49,
	0x5a, 0x45, 0x10, 0x04, 0x12, 0x12, 0x0a, 0x0e, 0x4e, 0x4f, 0x54, 0x5f, 0x43, 0x4f, 0x4d, 0x50,
	0x41, 0x54, 0x49, 0x42, 0x4c, 0x45, 0x10, 0x05, 0x12, 0x12, 0x0a, 0x0e, 0x49, 0x4e, 0x54, 0x45,
	0x52, 0x4e, 0x41, 0x4c, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x06, 0x12, 0x15, 0x0a, 0x11,
	0x50, 0x52, 0x4f, 0x47, 0x52, 0x41, 0x4d, 0x4d, 0x49, 0x4e, 0x47, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x10, 0x07, 0x12, 0x13, 0x0a, 0x0f, 0x4e, 0x4f, 0x5f, 0x48, 0x57, 0x5f, 0x49, 0x4e, 0x56,
	0x45, 0x4e, 0x54, 0x4f, 0x52, 0x59, 0x10, 0x08, 0x12, 0x1f, 0x0a, 0x1b, 0x54, 0x48, 0x49, 0x53,
	0x5f, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x5f,
	0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59, 0x10, 0x09, 0x42, 0x0f, 0x5a, 0x0d, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_io4edge_core_api_proto_rawDescOnce sync.Once
	file_io4edge_core_api_proto_rawDescData = file_io4edge_core_api_proto_rawDesc
)

func file_io4edge_core_api_proto_rawDescGZIP() []byte {
	file_io4edge_core_api_proto_rawDescOnce.Do(func() {
		file_io4edge_core_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_io4edge_core_api_proto_rawDescData)
	})
	return file_io4edge_core_api_proto_rawDescData
}

var file_io4edge_core_api_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_io4edge_core_api_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_io4edge_core_api_proto_goTypes = []interface{}{
	(CommandId)(0),                               // 0: io4edgeCoreApi.CommandId
	(Status)(0),                                  // 1: io4edgeCoreApi.Status
	(*SerialNumber)(nil),                         // 2: io4edgeCoreApi.SerialNumber
	(*LoadFirmwareChunkCommand)(nil),             // 3: io4edgeCoreApi.LoadFirmwareChunkCommand
	(*ProgramHardwareIdentificationCommand)(nil), // 4: io4edgeCoreApi.ProgramHardwareIdentificationCommand
	(*ProgramDeviceIdentificationCommand)(nil),   // 5: io4edgeCoreApi.ProgramDeviceIdentificationCommand
	(*IdentifyHardwareResponse)(nil),             // 6: io4edgeCoreApi.IdentifyHardwareResponse
	(*IdentifyFirmwareResponse)(nil),             // 7: io4edgeCoreApi.IdentifyFirmwareResponse
	(*CoreCommand)(nil),                          // 8: io4edgeCoreApi.CoreCommand
	(*CoreResponse)(nil),                         // 9: io4edgeCoreApi.CoreResponse
}
var file_io4edge_core_api_proto_depIdxs = []int32{
	2,  // 0: io4edgeCoreApi.ProgramHardwareIdentificationCommand.serial_number:type_name -> io4edgeCoreApi.SerialNumber
	2,  // 1: io4edgeCoreApi.IdentifyHardwareResponse.serial_number:type_name -> io4edgeCoreApi.SerialNumber
	0,  // 2: io4edgeCoreApi.CoreCommand.id:type_name -> io4edgeCoreApi.CommandId
	3,  // 3: io4edgeCoreApi.CoreCommand.load_firmware_chunk:type_name -> io4edgeCoreApi.LoadFirmwareChunkCommand
	4,  // 4: io4edgeCoreApi.CoreCommand.program_hardware_identification:type_name -> io4edgeCoreApi.ProgramHardwareIdentificationCommand
	5,  // 5: io4edgeCoreApi.CoreCommand.program_device_identification:type_name -> io4edgeCoreApi.ProgramDeviceIdentificationCommand
	0,  // 6: io4edgeCoreApi.CoreResponse.id:type_name -> io4edgeCoreApi.CommandId
	1,  // 7: io4edgeCoreApi.CoreResponse.status:type_name -> io4edgeCoreApi.Status
	6,  // 8: io4edgeCoreApi.CoreResponse.identify_hardware:type_name -> io4edgeCoreApi.IdentifyHardwareResponse
	7,  // 9: io4edgeCoreApi.CoreResponse.identify_firmware:type_name -> io4edgeCoreApi.IdentifyFirmwareResponse
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_io4edge_core_api_proto_init() }
func file_io4edge_core_api_proto_init() {
	if File_io4edge_core_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_io4edge_core_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SerialNumber); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_io4edge_core_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadFirmwareChunkCommand); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_io4edge_core_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProgramHardwareIdentificationCommand); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_io4edge_core_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProgramDeviceIdentificationCommand); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_io4edge_core_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdentifyHardwareResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_io4edge_core_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdentifyFirmwareResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_io4edge_core_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoreCommand); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_io4edge_core_api_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoreResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_io4edge_core_api_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*CoreCommand_LoadFirmwareChunk)(nil),
		(*CoreCommand_ProgramHardwareIdentification)(nil),
		(*CoreCommand_ProgramDeviceIdentification)(nil),
	}
	file_io4edge_core_api_proto_msgTypes[7].OneofWrappers = []interface{}{
		(*CoreResponse_IdentifyHardware)(nil),
		(*CoreResponse_IdentifyFirmware)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_io4edge_core_api_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_io4edge_core_api_proto_goTypes,
		DependencyIndexes: file_io4edge_core_api_proto_depIdxs,
		EnumInfos:         file_io4edge_core_api_proto_enumTypes,
		MessageInfos:      file_io4edge_core_api_proto_msgTypes,
	}.Build()
	File_io4edge_core_api_proto = out.File
	file_io4edge_core_api_proto_rawDesc = nil
	file_io4edge_core_api_proto_goTypes = nil
	file_io4edge_core_api_proto_depIdxs = nil
}
