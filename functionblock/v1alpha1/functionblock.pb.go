// Code generated by protoc-gen-go. DO NOT EDIT.
// source: functionblock.proto

package v1alpha1

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// --------- Responses ------------
type Status int32

const (
	Status_OK    Status = 0
	Status_ERROR Status = 1
)

var Status_name = map[int32]string{
	0: "OK",
	1: "ERROR",
}

var Status_value = map[string]int32{
	"OK":    0,
	"ERROR": 1,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}

func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{0}
}

// -------- Meta ------------
type Context struct {
	// A message identifying key for a command-response pairs, e.g. an UUID the
	// clients sends on the request.
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Context) Reset()         { *m = Context{} }
func (m *Context) String() string { return proto.CompactTextString(m) }
func (*Context) ProtoMessage()    {}
func (*Context) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{0}
}

func (m *Context) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Context.Unmarshal(m, b)
}
func (m *Context) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Context.Marshal(b, m, deterministic)
}
func (m *Context) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Context.Merge(m, src)
}
func (m *Context) XXX_Size() int {
	return xxx_messageInfo_Context.Size(m)
}
func (m *Context) XXX_DiscardUnknown() {
	xxx_messageInfo_Context.DiscardUnknown(m)
}

var xxx_messageInfo_Context proto.InternalMessageInfo

func (m *Context) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// ------- Commands ---------
type Command struct {
	Context *Context `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	// Types that are valid to be assigned to Type:
	//	*Command_Configuration
	//	*Command_FunctionControl
	//	*Command_StreamControl
	Type                 isCommand_Type `protobuf_oneof:"type"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{1}
}

func (m *Command) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Command.Unmarshal(m, b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Command.Marshal(b, m, deterministic)
}
func (m *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(m, src)
}
func (m *Command) XXX_Size() int {
	return xxx_messageInfo_Command.Size(m)
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

func (m *Command) GetContext() *Context {
	if m != nil {
		return m.Context
	}
	return nil
}

type isCommand_Type interface {
	isCommand_Type()
}

type Command_Configuration struct {
	Configuration *ConfigurationControl `protobuf:"bytes,2,opt,name=configuration,proto3,oneof"`
}

type Command_FunctionControl struct {
	FunctionControl *FunctionControl `protobuf:"bytes,3,opt,name=functionControl,proto3,oneof"`
}

type Command_StreamControl struct {
	StreamControl *StreamControl `protobuf:"bytes,4,opt,name=streamControl,proto3,oneof"`
}

func (*Command_Configuration) isCommand_Type() {}

func (*Command_FunctionControl) isCommand_Type() {}

func (*Command_StreamControl) isCommand_Type() {}

func (m *Command) GetType() isCommand_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Command) GetConfiguration() *ConfigurationControl {
	if x, ok := m.GetType().(*Command_Configuration); ok {
		return x.Configuration
	}
	return nil
}

func (m *Command) GetFunctionControl() *FunctionControl {
	if x, ok := m.GetType().(*Command_FunctionControl); ok {
		return x.FunctionControl
	}
	return nil
}

func (m *Command) GetStreamControl() *StreamControl {
	if x, ok := m.GetType().(*Command_StreamControl); ok {
		return x.StreamControl
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Command) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Command_Configuration)(nil),
		(*Command_FunctionControl)(nil),
		(*Command_StreamControl)(nil),
	}
}

// Configuration contains the function blocks high level configuration
type ConfigurationControl struct {
	FunctionSpecificHardwareConfiguration *any.Any `protobuf:"bytes,10,opt,name=functionSpecificHardwareConfiguration,proto3" json:"functionSpecificHardwareConfiguration,omitempty"`
	XXX_NoUnkeyedLiteral                  struct{} `json:"-"`
	XXX_unrecognized                      []byte   `json:"-"`
	XXX_sizecache                         int32    `json:"-"`
}

func (m *ConfigurationControl) Reset()         { *m = ConfigurationControl{} }
func (m *ConfigurationControl) String() string { return proto.CompactTextString(m) }
func (*ConfigurationControl) ProtoMessage()    {}
func (*ConfigurationControl) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{2}
}

func (m *ConfigurationControl) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigurationControl.Unmarshal(m, b)
}
func (m *ConfigurationControl) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigurationControl.Marshal(b, m, deterministic)
}
func (m *ConfigurationControl) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigurationControl.Merge(m, src)
}
func (m *ConfigurationControl) XXX_Size() int {
	return xxx_messageInfo_ConfigurationControl.Size(m)
}
func (m *ConfigurationControl) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigurationControl.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigurationControl proto.InternalMessageInfo

func (m *ConfigurationControl) GetFunctionSpecificHardwareConfiguration() *any.Any {
	if m != nil {
		return m.FunctionSpecificHardwareConfiguration
	}
	return nil
}

// FunctionControl specifies the direct function control for getting and setting
// values
type FunctionControl struct {
	// Types that are valid to be assigned to Action:
	//	*FunctionControl_FunctionSpecificFunctionControlSet
	//	*FunctionControl_FunctionSpecificFunctionControlGet
	Action               isFunctionControl_Action `protobuf_oneof:"action"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *FunctionControl) Reset()         { *m = FunctionControl{} }
func (m *FunctionControl) String() string { return proto.CompactTextString(m) }
func (*FunctionControl) ProtoMessage()    {}
func (*FunctionControl) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{3}
}

func (m *FunctionControl) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FunctionControl.Unmarshal(m, b)
}
func (m *FunctionControl) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FunctionControl.Marshal(b, m, deterministic)
}
func (m *FunctionControl) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FunctionControl.Merge(m, src)
}
func (m *FunctionControl) XXX_Size() int {
	return xxx_messageInfo_FunctionControl.Size(m)
}
func (m *FunctionControl) XXX_DiscardUnknown() {
	xxx_messageInfo_FunctionControl.DiscardUnknown(m)
}

var xxx_messageInfo_FunctionControl proto.InternalMessageInfo

type isFunctionControl_Action interface {
	isFunctionControl_Action()
}

type FunctionControl_FunctionSpecificFunctionControlSet struct {
	FunctionSpecificFunctionControlSet *any.Any `protobuf:"bytes,1,opt,name=functionSpecificFunctionControlSet,proto3,oneof"`
}

type FunctionControl_FunctionSpecificFunctionControlGet struct {
	FunctionSpecificFunctionControlGet *any.Any `protobuf:"bytes,2,opt,name=functionSpecificFunctionControlGet,proto3,oneof"`
}

func (*FunctionControl_FunctionSpecificFunctionControlSet) isFunctionControl_Action() {}

func (*FunctionControl_FunctionSpecificFunctionControlGet) isFunctionControl_Action() {}

func (m *FunctionControl) GetAction() isFunctionControl_Action {
	if m != nil {
		return m.Action
	}
	return nil
}

func (m *FunctionControl) GetFunctionSpecificFunctionControlSet() *any.Any {
	if x, ok := m.GetAction().(*FunctionControl_FunctionSpecificFunctionControlSet); ok {
		return x.FunctionSpecificFunctionControlSet
	}
	return nil
}

func (m *FunctionControl) GetFunctionSpecificFunctionControlGet() *any.Any {
	if x, ok := m.GetAction().(*FunctionControl_FunctionSpecificFunctionControlGet); ok {
		return x.FunctionSpecificFunctionControlGet
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*FunctionControl) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*FunctionControl_FunctionSpecificFunctionControlSet)(nil),
		(*FunctionControl_FunctionSpecificFunctionControlGet)(nil),
	}
}

// StreamControl specifies if the stream shall be started or stopped
type StreamControl struct {
	// Types that are valid to be assigned to Action:
	//	*StreamControl_FunctionSpecificStreamControlStart
	//	*StreamControl_Stop_
	Action               isStreamControl_Action `protobuf_oneof:"action"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *StreamControl) Reset()         { *m = StreamControl{} }
func (m *StreamControl) String() string { return proto.CompactTextString(m) }
func (*StreamControl) ProtoMessage()    {}
func (*StreamControl) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{4}
}

func (m *StreamControl) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamControl.Unmarshal(m, b)
}
func (m *StreamControl) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamControl.Marshal(b, m, deterministic)
}
func (m *StreamControl) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamControl.Merge(m, src)
}
func (m *StreamControl) XXX_Size() int {
	return xxx_messageInfo_StreamControl.Size(m)
}
func (m *StreamControl) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamControl.DiscardUnknown(m)
}

var xxx_messageInfo_StreamControl proto.InternalMessageInfo

type isStreamControl_Action interface {
	isStreamControl_Action()
}

type StreamControl_FunctionSpecificStreamControlStart struct {
	FunctionSpecificStreamControlStart *any.Any `protobuf:"bytes,1,opt,name=functionSpecificStreamControlStart,proto3,oneof"`
}

type StreamControl_Stop_ struct {
	Stop *StreamControl_Stop `protobuf:"bytes,2,opt,name=stop,proto3,oneof"`
}

func (*StreamControl_FunctionSpecificStreamControlStart) isStreamControl_Action() {}

func (*StreamControl_Stop_) isStreamControl_Action() {}

func (m *StreamControl) GetAction() isStreamControl_Action {
	if m != nil {
		return m.Action
	}
	return nil
}

func (m *StreamControl) GetFunctionSpecificStreamControlStart() *any.Any {
	if x, ok := m.GetAction().(*StreamControl_FunctionSpecificStreamControlStart); ok {
		return x.FunctionSpecificStreamControlStart
	}
	return nil
}

func (m *StreamControl) GetStop() *StreamControl_Stop {
	if x, ok := m.GetAction().(*StreamControl_Stop_); ok {
		return x.Stop
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*StreamControl) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*StreamControl_FunctionSpecificStreamControlStart)(nil),
		(*StreamControl_Stop_)(nil),
	}
}

type StreamControl_Stop struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamControl_Stop) Reset()         { *m = StreamControl_Stop{} }
func (m *StreamControl_Stop) String() string { return proto.CompactTextString(m) }
func (*StreamControl_Stop) ProtoMessage()    {}
func (*StreamControl_Stop) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{4, 0}
}

func (m *StreamControl_Stop) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamControl_Stop.Unmarshal(m, b)
}
func (m *StreamControl_Stop) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamControl_Stop.Marshal(b, m, deterministic)
}
func (m *StreamControl_Stop) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamControl_Stop.Merge(m, src)
}
func (m *StreamControl_Stop) XXX_Size() int {
	return xxx_messageInfo_StreamControl_Stop.Size(m)
}
func (m *StreamControl_Stop) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamControl_Stop.DiscardUnknown(m)
}

var xxx_messageInfo_StreamControl_Stop proto.InternalMessageInfo

type Error struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{5}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type HardwareConfigurationResponse struct {
	Status                                        Status   `protobuf:"varint,1,opt,name=status,proto3,enum=functionblock.Status" json:"status,omitempty"`
	Error                                         *Error   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	FunctionSpecificHardwareConfigurationResponse *any.Any `protobuf:"bytes,10,opt,name=functionSpecificHardwareConfigurationResponse,proto3" json:"functionSpecificHardwareConfigurationResponse,omitempty"`
	XXX_NoUnkeyedLiteral                          struct{} `json:"-"`
	XXX_unrecognized                              []byte   `json:"-"`
	XXX_sizecache                                 int32    `json:"-"`
}

func (m *HardwareConfigurationResponse) Reset()         { *m = HardwareConfigurationResponse{} }
func (m *HardwareConfigurationResponse) String() string { return proto.CompactTextString(m) }
func (*HardwareConfigurationResponse) ProtoMessage()    {}
func (*HardwareConfigurationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{6}
}

func (m *HardwareConfigurationResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HardwareConfigurationResponse.Unmarshal(m, b)
}
func (m *HardwareConfigurationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HardwareConfigurationResponse.Marshal(b, m, deterministic)
}
func (m *HardwareConfigurationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HardwareConfigurationResponse.Merge(m, src)
}
func (m *HardwareConfigurationResponse) XXX_Size() int {
	return xxx_messageInfo_HardwareConfigurationResponse.Size(m)
}
func (m *HardwareConfigurationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HardwareConfigurationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HardwareConfigurationResponse proto.InternalMessageInfo

func (m *HardwareConfigurationResponse) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_OK
}

func (m *HardwareConfigurationResponse) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *HardwareConfigurationResponse) GetFunctionSpecificHardwareConfigurationResponse() *any.Any {
	if m != nil {
		return m.FunctionSpecificHardwareConfigurationResponse
	}
	return nil
}

type FunctionControlResponse struct {
	Status                                  Status   `protobuf:"varint,1,opt,name=status,proto3,enum=functionblock.Status" json:"status,omitempty"`
	Error                                   *Error   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	FunctionSpecificFunctionControlResponse *any.Any `protobuf:"bytes,10,opt,name=functionSpecificFunctionControlResponse,proto3" json:"functionSpecificFunctionControlResponse,omitempty"`
	XXX_NoUnkeyedLiteral                    struct{} `json:"-"`
	XXX_unrecognized                        []byte   `json:"-"`
	XXX_sizecache                           int32    `json:"-"`
}

func (m *FunctionControlResponse) Reset()         { *m = FunctionControlResponse{} }
func (m *FunctionControlResponse) String() string { return proto.CompactTextString(m) }
func (*FunctionControlResponse) ProtoMessage()    {}
func (*FunctionControlResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{7}
}

func (m *FunctionControlResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FunctionControlResponse.Unmarshal(m, b)
}
func (m *FunctionControlResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FunctionControlResponse.Marshal(b, m, deterministic)
}
func (m *FunctionControlResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FunctionControlResponse.Merge(m, src)
}
func (m *FunctionControlResponse) XXX_Size() int {
	return xxx_messageInfo_FunctionControlResponse.Size(m)
}
func (m *FunctionControlResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FunctionControlResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FunctionControlResponse proto.InternalMessageInfo

func (m *FunctionControlResponse) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_OK
}

func (m *FunctionControlResponse) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *FunctionControlResponse) GetFunctionSpecificFunctionControlResponse() *any.Any {
	if m != nil {
		return m.FunctionSpecificFunctionControlResponse
	}
	return nil
}

type StreamControlResponse struct {
	Status               Status   `protobuf:"varint,1,opt,name=status,proto3,enum=functionblock.Status" json:"status,omitempty"`
	Error                *Error   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamControlResponse) Reset()         { *m = StreamControlResponse{} }
func (m *StreamControlResponse) String() string { return proto.CompactTextString(m) }
func (*StreamControlResponse) ProtoMessage()    {}
func (*StreamControlResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{8}
}

func (m *StreamControlResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamControlResponse.Unmarshal(m, b)
}
func (m *StreamControlResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamControlResponse.Marshal(b, m, deterministic)
}
func (m *StreamControlResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamControlResponse.Merge(m, src)
}
func (m *StreamControlResponse) XXX_Size() int {
	return xxx_messageInfo_StreamControlResponse.Size(m)
}
func (m *StreamControlResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamControlResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StreamControlResponse proto.InternalMessageInfo

func (m *StreamControlResponse) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_OK
}

func (m *StreamControlResponse) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

type StreamResponse struct {
	// timestamp when the message has been sent out
	DeliveryTimestampUs uint64 `protobuf:"fixed64,1,opt,name=deliveryTimestampUs,proto3" json:"deliveryTimestampUs,omitempty"`
	// sample series sequence number (counted from 0, rolls over)
	Sequence uint32 `protobuf:"fixed32,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	// Function specific data type
	FunctionSpecificStreamResponse *any.Any `protobuf:"bytes,10,opt,name=functionSpecificStreamResponse,proto3" json:"functionSpecificStreamResponse,omitempty"`
	XXX_NoUnkeyedLiteral           struct{} `json:"-"`
	XXX_unrecognized               []byte   `json:"-"`
	XXX_sizecache                  int32    `json:"-"`
}

func (m *StreamResponse) Reset()         { *m = StreamResponse{} }
func (m *StreamResponse) String() string { return proto.CompactTextString(m) }
func (*StreamResponse) ProtoMessage()    {}
func (*StreamResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{9}
}

func (m *StreamResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamResponse.Unmarshal(m, b)
}
func (m *StreamResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamResponse.Marshal(b, m, deterministic)
}
func (m *StreamResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamResponse.Merge(m, src)
}
func (m *StreamResponse) XXX_Size() int {
	return xxx_messageInfo_StreamResponse.Size(m)
}
func (m *StreamResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StreamResponse proto.InternalMessageInfo

func (m *StreamResponse) GetDeliveryTimestampUs() uint64 {
	if m != nil {
		return m.DeliveryTimestampUs
	}
	return 0
}

func (m *StreamResponse) GetSequence() uint32 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *StreamResponse) GetFunctionSpecificStreamResponse() *any.Any {
	if m != nil {
		return m.FunctionSpecificStreamResponse
	}
	return nil
}

type Response struct {
	Context *Context `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	// Types that are valid to be assigned to Type:
	//	*Response_Configuration
	//	*Response_FunctionControl
	//	*Response_StreamControl
	//	*Response_Stream
	Type                 isResponse_Type `protobuf_oneof:"type"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_296c9f3b85f1c061, []int{10}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetContext() *Context {
	if m != nil {
		return m.Context
	}
	return nil
}

type isResponse_Type interface {
	isResponse_Type()
}

type Response_Configuration struct {
	Configuration *HardwareConfigurationResponse `protobuf:"bytes,2,opt,name=configuration,proto3,oneof"`
}

type Response_FunctionControl struct {
	FunctionControl *FunctionControlResponse `protobuf:"bytes,3,opt,name=functionControl,proto3,oneof"`
}

type Response_StreamControl struct {
	StreamControl *StreamControlResponse `protobuf:"bytes,4,opt,name=streamControl,proto3,oneof"`
}

type Response_Stream struct {
	Stream *StreamResponse `protobuf:"bytes,5,opt,name=stream,proto3,oneof"`
}

func (*Response_Configuration) isResponse_Type() {}

func (*Response_FunctionControl) isResponse_Type() {}

func (*Response_StreamControl) isResponse_Type() {}

func (*Response_Stream) isResponse_Type() {}

func (m *Response) GetType() isResponse_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Response) GetConfiguration() *HardwareConfigurationResponse {
	if x, ok := m.GetType().(*Response_Configuration); ok {
		return x.Configuration
	}
	return nil
}

func (m *Response) GetFunctionControl() *FunctionControlResponse {
	if x, ok := m.GetType().(*Response_FunctionControl); ok {
		return x.FunctionControl
	}
	return nil
}

func (m *Response) GetStreamControl() *StreamControlResponse {
	if x, ok := m.GetType().(*Response_StreamControl); ok {
		return x.StreamControl
	}
	return nil
}

func (m *Response) GetStream() *StreamResponse {
	if x, ok := m.GetType().(*Response_Stream); ok {
		return x.Stream
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Response) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Response_Configuration)(nil),
		(*Response_FunctionControl)(nil),
		(*Response_StreamControl)(nil),
		(*Response_Stream)(nil),
	}
}

func init() {
	proto.RegisterEnum("functionblock.Status", Status_name, Status_value)
	proto.RegisterType((*Context)(nil), "functionblock.Context")
	proto.RegisterType((*Command)(nil), "functionblock.Command")
	proto.RegisterType((*ConfigurationControl)(nil), "functionblock.ConfigurationControl")
	proto.RegisterType((*FunctionControl)(nil), "functionblock.FunctionControl")
	proto.RegisterType((*StreamControl)(nil), "functionblock.StreamControl")
	proto.RegisterType((*StreamControl_Stop)(nil), "functionblock.StreamControl.Stop")
	proto.RegisterType((*Error)(nil), "functionblock.Error")
	proto.RegisterType((*HardwareConfigurationResponse)(nil), "functionblock.HardwareConfigurationResponse")
	proto.RegisterType((*FunctionControlResponse)(nil), "functionblock.FunctionControlResponse")
	proto.RegisterType((*StreamControlResponse)(nil), "functionblock.StreamControlResponse")
	proto.RegisterType((*StreamResponse)(nil), "functionblock.StreamResponse")
	proto.RegisterType((*Response)(nil), "functionblock.Response")
}

func init() { proto.RegisterFile("functionblock.proto", fileDescriptor_296c9f3b85f1c061) }

var fileDescriptor_296c9f3b85f1c061 = []byte{
	// 629 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0xcd, 0x72, 0xd3, 0x3e,
	0x10, 0x8f, 0xf3, 0x4f, 0xdc, 0x76, 0xff, 0x93, 0xb6, 0xa3, 0x7e, 0x60, 0x0a, 0x29, 0x60, 0x3e,
	0xa7, 0x43, 0x9d, 0xb6, 0x1c, 0x7a, 0x26, 0x50, 0x9a, 0xa1, 0xcc, 0x74, 0x46, 0x29, 0x17, 0x86,
	0x8b, 0xe2, 0x28, 0xc1, 0xe0, 0x58, 0x46, 0x96, 0x03, 0xe1, 0xc8, 0xa3, 0xf0, 0x1c, 0x9c, 0x78,
	0x16, 0x6e, 0x9c, 0xe1, 0xcc, 0x58, 0xb6, 0x8b, 0xe5, 0x98, 0xd4, 0xe1, 0xd0, 0x5b, 0x56, 0xfa,
	0x7d, 0xec, 0x66, 0x57, 0x5e, 0x58, 0x1b, 0x84, 0x9e, 0x2d, 0x1c, 0xe6, 0xf5, 0x5c, 0x66, 0xbf,
	0xb3, 0x7c, 0xce, 0x04, 0x43, 0x0d, 0xe5, 0x70, 0xeb, 0xea, 0x90, 0xb1, 0xa1, 0x4b, 0x5b, 0xf2,
	0xb2, 0x17, 0x0e, 0x5a, 0xc4, 0x9b, 0xc4, 0x48, 0xf3, 0x06, 0x2c, 0x3c, 0x61, 0x9e, 0xa0, 0x1f,
	0x05, 0x5a, 0x87, 0xfa, 0x98, 0xb8, 0x21, 0x35, 0xb4, 0x9b, 0xda, 0x83, 0x25, 0x1c, 0x07, 0xe6,
	0x97, 0x6a, 0x84, 0x18, 0x8d, 0x88, 0xd7, 0x47, 0x7b, 0xb0, 0x60, 0xc7, 0x60, 0x89, 0xf9, 0xff,
	0x60, 0xd3, 0x52, 0xdd, 0x13, 0x29, 0x9c, 0xc2, 0xd0, 0x09, 0x34, 0x6c, 0xe6, 0x0d, 0x9c, 0x61,
	0xc8, 0x49, 0x04, 0x33, 0xaa, 0x92, 0x77, 0x7b, 0x9a, 0xf7, 0x07, 0x13, 0x89, 0x70, 0xe6, 0x76,
	0x2a, 0x58, 0xe5, 0xa2, 0xe7, 0xb0, 0x92, 0xd2, 0x12, 0x8c, 0xf1, 0x9f, 0x94, 0xdb, 0xce, 0xc9,
	0x3d, 0x53, 0x51, 0x9d, 0x0a, 0xce, 0x13, 0xd1, 0x53, 0x68, 0x04, 0x82, 0x53, 0x32, 0x4a, 0x95,
	0x6a, 0x52, 0xe9, 0x7a, 0x4e, 0xa9, 0x9b, 0xc5, 0x44, 0x19, 0x29, 0xa4, 0xb6, 0x0e, 0x35, 0x31,
	0xf1, 0xa9, 0xf9, 0x59, 0x83, 0xf5, 0xa2, 0x1a, 0xd0, 0x5b, 0xb8, 0x9b, 0x0a, 0x76, 0x7d, 0x6a,
	0x3b, 0x03, 0xc7, 0xee, 0x10, 0xde, 0xff, 0x40, 0x38, 0x55, 0xf0, 0x06, 0x48, 0xfb, 0x75, 0x2b,
	0xee, 0x94, 0x95, 0x76, 0xca, 0x7a, 0xec, 0x4d, 0x70, 0x39, 0x09, 0xf3, 0x87, 0x06, 0x2b, 0xb9,
	0xca, 0xd1, 0x00, 0xcc, 0x3c, 0x39, 0x07, 0xe9, 0xd2, 0xb4, 0x99, 0x85, 0xe6, 0x9d, 0x0a, 0x2e,
	0xa1, 0x50, 0xc2, 0xe7, 0x98, 0x8a, 0xa4, 0xf9, 0xff, 0xea, 0x73, 0x4c, 0x45, 0x7b, 0x11, 0x74,
	0x22, 0xcf, 0xcc, 0x6f, 0x1a, 0x34, 0x94, 0xee, 0x14, 0xe5, 0xa0, 0x00, 0xba, 0x82, 0xf0, 0xb9,
	0x6b, 0x9d, 0x56, 0x40, 0x87, 0x50, 0x0b, 0x04, 0xf3, 0x93, 0x6a, 0x6e, 0xcd, 0x9a, 0x18, 0xab,
	0x2b, 0x98, 0xdf, 0xa9, 0x60, 0x49, 0xd8, 0xd2, 0xa1, 0x16, 0xc5, 0x99, 0x22, 0x9a, 0x50, 0x3f,
	0xe2, 0x9c, 0xf1, 0xe8, 0xed, 0xd1, 0xe8, 0x47, 0xfa, 0xf6, 0x64, 0x60, 0xfe, 0xd2, 0xa0, 0x59,
	0xd8, 0x6b, 0x4c, 0x03, 0x9f, 0x79, 0x01, 0x45, 0xbb, 0xa0, 0x07, 0x82, 0x88, 0x30, 0x90, 0xc4,
	0xe5, 0x83, 0x8d, 0xa9, 0x6c, 0xa2, 0x4b, 0x9c, 0x80, 0xd0, 0x4e, 0x6a, 0x93, 0x76, 0x42, 0x45,
	0xcb, 0x5c, 0x12, 0x73, 0xf4, 0x09, 0x76, 0x4b, 0xcd, 0x5d, 0x9a, 0xcb, 0xcc, 0x11, 0x9e, 0x4f,
	0xca, 0xfc, 0xae, 0xc1, 0x95, 0x5c, 0xf7, 0x2f, 0xa3, 0x64, 0x0f, 0xee, 0x5f, 0x30, 0x83, 0xa5,
	0x8a, 0x2d, 0x2b, 0x62, 0x72, 0xd8, 0x50, 0xc6, 0xe5, 0x12, 0x6a, 0x34, 0xbf, 0x6a, 0xb0, 0x1c,
	0x9b, 0x9e, 0xbb, 0xed, 0xc1, 0x5a, 0x9f, 0xba, 0xce, 0x98, 0xf2, 0xc9, 0x99, 0x33, 0xa2, 0x81,
	0x20, 0x23, 0xff, 0x65, 0x6c, 0xad, 0xe3, 0xa2, 0x2b, 0xb4, 0x05, 0x8b, 0x01, 0x7d, 0x1f, 0x52,
	0xcf, 0xa6, 0xd2, 0x73, 0x01, 0x9f, 0xc7, 0xe8, 0x35, 0x6c, 0x17, 0x3f, 0xa2, 0x52, 0xff, 0xdd,
	0x05, 0x5c, 0xf3, 0x67, 0x15, 0x16, 0x33, 0x89, 0xcf, 0xbb, 0x8f, 0xce, 0x8a, 0xf7, 0xd1, 0xc3,
	0x1c, 0x6f, 0xe6, 0x74, 0x4e, 0x2f, 0x26, 0xfc, 0xb7, 0xc5, 0x74, 0x6f, 0xf6, 0x62, 0xca, 0x28,
	0x4e, 0x2d, 0xa8, 0x17, 0xc5, 0x0b, 0xea, 0xce, 0xac, 0xcf, 0x4d, 0x36, 0x43, 0x85, 0x8c, 0x0e,
	0xa3, 0x81, 0x8a, 0x0e, 0x8c, 0xba, 0x94, 0x69, 0x16, 0xca, 0x64, 0xf8, 0x09, 0x3c, 0xdd, 0x70,
	0x3b, 0xd7, 0x40, 0x8f, 0x87, 0x0e, 0xe9, 0x50, 0x3d, 0x3d, 0x59, 0xad, 0xa0, 0x25, 0xa8, 0x1f,
	0x61, 0x7c, 0x8a, 0x57, 0xb5, 0xb6, 0xf1, 0x6a, 0x53, 0x91, 0x6b, 0x8d, 0xf7, 0x89, 0xeb, 0xbf,
	0x21, 0xfb, 0x3d, 0x5d, 0x36, 0xf7, 0xd1, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8f, 0xc5, 0xd7,
	0xfb, 0xa6, 0x08, 0x00, 0x00,
}
