package core

import (
	"bufio"
	"io"
	"time"
)

// If represents a client for the io4edge core functions
// Implemented by pbcore and httpscore packages
type If interface {
	IdentifyFirmware(timeout time.Duration) (name string, version string, err error)
	LoadFirmware(file string, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error)
	LoadFirmwareFromBuffer(buffer []byte, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error)
	LoadFirmwareBinaryFromFile(file string, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error)
	LoadFirmwareBinary(r *bufio.Reader, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error)
	IdentifyHardware(timeout time.Duration) (*HardwareInventory, error)
	ProgramHardwareIdentification(i *HardwareInventory, timeout time.Duration) error
	ReadPartition(timeout time.Duration, partitionName string, offset uint32, w *bufio.Writer, prog func(bytes uint, msg string)) (err error)
	SetPersistentParameter(name string, value string, timeout time.Duration) (bool, error)
	GetPersistentParameter(name string, timeout time.Duration) (value string, err error)
	ResetReason(timeout time.Duration) (reason string, err error)
	Restart(timeout time.Duration) (restartingNow bool, err error)
	StreamLogs(streamTimeout time.Duration, infoCb func(msg string)) (io.ReadCloser, error)
	GetParameterSet(timeout time.Duration, namespace string) ([]byte, error)
	LoadParameterSet(timeout time.Duration, namespace string, data []byte) ([]byte, error)
	Close()
}

// FirmwareAlreadyPresentError is returned by LoadFirmware as a dummy error
type FirmwareAlreadyPresentError struct {
}

// Error returns the error string for FirmwareAlreadyPresentError
func (e *FirmwareAlreadyPresentError) Error() string {
	return "Requested Firmware already present"
}

// ParameterIsReadProtectedError is returned by SetPersistentParameter if the parameter is read protected
type ParameterIsReadProtectedError struct {
}

func (e *ParameterIsReadProtectedError) Error() string {
	return "Parameter is read protected"
}

// HardwareInventory represents the hardware inventory information
type HardwareInventory struct {
	PartNumber      string
	SerialNumber    string
	MajorVersion    uint32
	CustomExtension map[string]string // may contain e.g. customer part number
}
