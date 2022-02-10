package client

import (
	"errors"
)

// FuncInfoDefault provides the default FunctionInfo implementation
type FuncInfoDefault struct {
	address     string
	auxPort     int
	auxProtocol string
	auxSchemaID string
}

// NewFuncInfoDefault creates a new FuncInfoDefault object with no aux port
func NewFuncInfoDefault(address string) *FuncInfoDefault {
	return &FuncInfoDefault{address: address}
}

// NewFuncInfoDefaultWithAuxPort creates a new FuncInfoDefault object with an aux port
func NewFuncInfoDefaultWithAuxPort(address string, auxPort int, auxProtocol string, auxSchemaID string) *FuncInfoDefault {
	return &FuncInfoDefault{address: address, auxPort: auxPort, auxProtocol: auxProtocol, auxSchemaID: auxSchemaID}
}

// NetAddress gives the caller the ip address and port of the service
func (f *FuncInfoDefault) NetAddress() (string, int, error) {
	return netAddressSplit(f.address)
}

// FuncClass gives the caller the funcclass value of the service
func (f *FuncInfoDefault) FuncClass() (string, error) {
	return "", errors.New("not implemented")
}

// Security gives the caller the security value of the service
func (f *FuncInfoDefault) Security() (string, error) {
	return "", errors.New("not implemented")
}

// AuxPort gives the caller the auxport value of the service protocol and port
func (f *FuncInfoDefault) AuxPort() (string, int, error) {
	if f.auxProtocol == "" || f.auxPort == 0 {
		return "", 0, errors.New("no aux port")
	}
	return f.auxProtocol, f.auxPort, nil
}

// AuxSchemaID gives the caller the auxschema value of the service
func (f *FuncInfoDefault) AuxSchemaID() (string, error) {
	if f.auxSchemaID == "" {
		return "", errors.New("no aux schema")
	}
	return f.auxSchemaID, nil
}
