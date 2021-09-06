package basefunc

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"time"
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

// LoadFirmwareFromFile loads new firmware from file into the device device
// timeout is for each chunk
func (c *Client) LoadFirmwareFromFile(file string, chunkSize uint, timeout time.Duration) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.New("cannot open file " + file + " : " + err.Error())
	}

	defer f.Close()

	r := bufio.NewReader(f)
	return c.LoadFirmware(r, chunkSize, timeout)
}

// LoadFirmware loads new firmware via r into the device device
// timeout is for each chunk
func (c *Client) LoadFirmware(r *bufio.Reader, chunkSize uint, timeout time.Duration) error {

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
