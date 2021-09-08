package basefunc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ci4rail/io4edge-client-go/pkg/io4edge/fwpkg"
)

// IdentifyFirmware gets the firmware name and version from the device
func (c *Client) IdentifyFirmware(timeout time.Duration) (*ResIdentifyFirmware, error) {
	cmd := &BaseFuncCommand{
		Id: BaseFuncCommandId_IDENTIFY_FIRMWARE,
	}
	res := &BaseFuncResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return nil, err
	}
	return res.GetIdentifyFirmware(), nil
}

// LoadFirmware loads a binary from a firmware package to the device
//
// timeout is for each chunk
func (c *Client) LoadFirmware(file string, chunkSize uint, timeout time.Duration) error {

	pkg, err := fwpkg.NewFirmwarePackageFromFile(file)
	if err != nil {
		return err
	}
	manifest := pkg.Manifest()

	// get currently running firmware
	currentFWID, err := c.IdentifyFirmware(timeout)
	if err != nil {
		return err
	}

	// get devices hardware id
	hwID, err := c.IdentifyHardware(timeout)
	if err != nil {
		return err
	}

	// check compatibility
	err = AssertFirmwareIsCompatibleWithHardware(
		manifest.Compatibility.HW,
		manifest.Compatibility.MajorRevs,
		hwID.RootArticle,
		int(hwID.MajorVersion),
	)
	if err != nil {
		return err
	}

	// check if fw already running
	curVer := fmt.Sprintf("%d.%d.%d", currentFWID.MajorVersion, currentFWID.MinorVersion, currentFWID.PatchVersion)

	if currentFWID.Name == manifest.Name && curVer == manifest.Version {
		return errors.New("firmware variant/version is already active on device")
	}

	fwFile := new(bytes.Buffer)
	err = pkg.File(fwFile)
	if err != nil {
		return err
	}

	err = c.LoadFirmwareBinary(bufio.NewReader(fwFile), chunkSize, timeout)
	if err != nil {
		return err
	}
	return err
}

// LoadFirmwareBinaryFromFile loads new firmware from file into the device device
// timeout is for each chunk
func (c *Client) LoadFirmwareBinaryFromFile(file string, chunkSize uint, timeout time.Duration) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.New("cannot open file " + file + " : " + err.Error())
	}

	defer f.Close()

	r := bufio.NewReader(f)
	return c.LoadFirmwareBinary(r, chunkSize, timeout)
}

// LoadFirmwareBinary loads new firmware via r into the device device
// timeout is for each chunk
func (c *Client) LoadFirmwareBinary(r *bufio.Reader, chunkSize uint, timeout time.Duration) error {

	cmd := &BaseFuncCommand{
		Id: BaseFuncCommandId_LOAD_FIRMWARE_CHUNK,
		Data: &BaseFuncCommand_LoadFirmwareChunk{
			LoadFirmwareChunk: &CmdLoadFirmwareChunk{
				Data: make([]byte, chunkSize),
			},
		},
	}
	data := cmd.GetLoadFirmwareChunk().Data

	chunkNumber := uint32(0)

	for {
		atEOF := false

		n, err := r.Read(data)
		if err != nil {
			return errors.New("read firmware failed: " + err.Error())
		}

		// check if we are at EOF
		_, err = r.Peek(1)
		if err == io.EOF {
			atEOF = true
		}
		log.Printf("Read %d bytes at_eof=%v chunk %d\n", n, atEOF, chunkNumber)

		cmd.GetLoadFirmwareChunk().IsLastChunk = atEOF
		cmd.GetLoadFirmwareChunk().ChunkNumber = chunkNumber

		res := &BaseFuncResponse{}
		err = c.Command(cmd, res, timeout)
		if err != nil {
			return errors.New("load firmware chunk command failed: " + err.Error())
		}

		if atEOF {
			break
		}
		chunkNumber++
	}
	log.Printf("Device accepted all chunks.\n")

	return nil
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
