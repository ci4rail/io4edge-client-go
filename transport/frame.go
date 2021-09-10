/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Implements message framing on a streaming transport (e.g. sockets)
//
// Frame format:
// - Magic Word (2 Bytes) 0xFE 0xED
// - Payload Len (4 Bytes, Little Endian)
// - Payload

package transport

import (
	"errors"
)

// FramedStream represents a stream with message semantics
type FramedStream struct {
	Trans Transport
}

// NewFramedStreamFromTransport creates a message stream from transport t
func NewFramedStreamFromTransport(t Transport) (*FramedStream, error) {
	return &FramedStream{
		Trans: t,
	}, nil
}

// WriteMsg writes io4edge standard message to the transport stream
func (fs *FramedStream) WriteMsg(payload []byte) error {
	// make sure we have the magic bytes
	err := fs.writeMagicBytes()
	if err != nil {
		return err
	}

	length := uint(len(payload))
	err = fs.writeLength(length)
	if err != nil {
		return err
	}

	err = fs.writePayload(payload)
	if err != nil {
		return err
	}
	return nil
}

// writeMagicBytes write the magic bytes 0xFE, 0xED to transport stream
func (fs *FramedStream) writeMagicBytes() error {
	magicBytes := []byte{0xFE, 0xED}

	err := fs.writeBytesSafe(magicBytes)
	return err
}

// writeLength writes 4 bytes to transport stream with the length
func (fs *FramedStream) writeLength(length uint) error {
	lengthBytes := make([]byte, 4)

	lengthBytes[0] = byte(length & 0xFF)
	lengthBytes[1] = byte((length >> 8) & 0xFF)
	lengthBytes[2] = byte((length >> 16) & 0xFF)
	lengthBytes[3] = byte((length >> 24) & 0xFF)

	err := fs.writeBytesSafe(lengthBytes)
	return err
}

// writePayload write the payload to transport stream.
func (fs *FramedStream) writePayload(payload []byte) error {
	err := fs.writeBytesSafe(payload)
	return err
}

// writeBytesSafe retries writing to transport stream until all bytes are written
func (fs *FramedStream) writeBytesSafe(payload []byte) error {
	for {
		written, err := fs.Trans.Write(payload)
		if err != nil {
			return err
		}
		if written == len(payload) {
			return nil
		}
		payload = payload[written:]
	}
}

// ReadMsg reads a io4edge standard message from transport stream
func (fs *FramedStream) ReadMsg() ([]byte, error) {
	// make sure we have the magic bytes
	err := fs.readMagicBytes()
	if err != nil {
		return nil, err
	}

	length, err := fs.readLength()
	if err != nil {
		return nil, err
	}
	payload, err := fs.readPayload(length)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// readMagicBytes blocks until it receives the magic bytes 0xFE, 0xED from transport stream.
func (fs *FramedStream) readMagicBytes() error {
	// block until we get the magic bytes
	for {
		magicBytes := make([]byte, 2)
		for i := 0; i < 2; i++ {
			b := make([]byte, 1)

			_, err := fs.Trans.Read(b)
			if err != nil {
				return err
			}
			magicBytes[i] = b[0]
		}
		if magicBytes[0] == 0xFE && magicBytes[1] == 0xED {
			return nil
		}
	}
}

// readLength reads 4 bytes from transport stream and returns the length as uint of the message.
func (fs *FramedStream) readLength() (uint, error) {
	lengthBytes := make([]byte, 4)
	_, err := fs.Trans.Read(lengthBytes)
	if err != nil {
		return 0, err
	}
	length := uint(lengthBytes[0])
	length |= uint(lengthBytes[1]) << 8
	length |= uint(lengthBytes[2]) << 16
	length |= uint(lengthBytes[3]) << 24
	return length, nil
}

// readPayload reads the payload from transport stream and returns it as []byte.
func (fs *FramedStream) readPayload(length uint) ([]byte, error) {
	payload := make([]byte, length)
	n, err := fs.Trans.Read(payload)
	if err != nil {
		return nil, err
	}
	if n != int(length) {
		return nil, errors.New("read too few bytes")
	}
	return payload, nil
}

// Close closes the transport stream
func (fs *FramedStream) Close() error {
	return fs.Trans.Close()
}
