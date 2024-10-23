package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// If represents a client for the io4edge core functions
// Implemented by pbcore and httpscore packages
type If interface {
	IdentifyFirmware(timeout time.Duration) (name string, version string, err error)
	LoadFirmware(file string, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error)
	LoadFirmwareBinaryFromFile(file string, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error)
	LoadFirmwareBinary(r *bufio.Reader, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error)
	IdentifyHardware(timeout time.Duration) (name string, major uint32, serial string, err error)
	ProgramHardwareIdentification(name string, major uint32, serial string, timeout time.Duration) error
	ReadPartition(timeout time.Duration, partitionName string, offset uint32, w *bufio.Writer, prog func(bytes uint, msg string)) (err error)
	SetPersistentParameter(name string, value string, timeout time.Duration) error
	GetPersistentParameter(name string, timeout time.Duration) (value string, err error)
	ResetReason(timeout time.Duration) (reason string, err error)
	Restart(timeout time.Duration) (restartingNow bool, err error)
	StreamLogs(streamTimeout time.Duration, infoCb func(msg string)) (io.ReadCloser, error)
}

// FirmwareAlreadyPresentError is returned by LoadFirmware as a dummy error
type FirmwareAlreadyPresentError struct {
}

// Error returns the error string for FirmwareAlreadyPresentError
func (e *FirmwareAlreadyPresentError) Error() string {
	return "Requested Firmware already present"
}

// AssertFirmwareIsCompatibleWithHardware checks if the firmware specified by fwHw and fwMajorRevs is compatible
// with hardware hwName, hwMajor
func AssertFirmwareIsCompatibleWithHardware(fwHw string, fwMajorRevs []int, hwName string, hwMajor int) error {
	if !strings.EqualFold(fwHw, hwName) {
		return errors.New("firmware " + fwHw + " not suitable for hardware " + hwName)
	}
	var ok = false
	for _, b := range fwMajorRevs {
		if hwMajor == b {
			ok = true
		}
	}
	if !ok {
		return fmt.Errorf("firmware doesn't support hardware version %d", hwMajor)
	}
	return nil
}
