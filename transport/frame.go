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
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// This is the ascii representation of the string 'Start indicator of io4edge message'
	magicBytes = []byte{0xFE, 0xED}
)

// FramedStream represents a stream with message semantics
type FramedStream struct {
	Trans Transport
	mutex sync.Mutex
}

// NewFramedStreamFromTransport creates a message stream from transport t
func NewFramedStreamFromTransport(t Transport) *FramedStream {
	return &FramedStream{
		Trans: t,
		mutex: sync.Mutex{},
	}
}

// WriteMsg writes io4edge standard message to the transport stream
func (fs *FramedStream) WriteMsg(payload []byte) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	msg := make([]byte, 0)

	msg = append(msg, magicBytes...)
	msg = append(msg, fs.genLength(uint(len(payload)))...)
	msg = append(msg, payload...)

	err := fs.writeBytesSafe(msg)
	return err
}

// genLength writes 4 bytes to transport stream with the length
func (fs *FramedStream) genLength(length uint) []byte {
	lengthBytes := make([]byte, 4)

	lengthBytes[0] = byte(length & 0xFF)
	lengthBytes[1] = byte((length >> 8) & 0xFF)
	lengthBytes[2] = byte((length >> 16) & 0xFF)
	lengthBytes[3] = byte((length >> 24) & 0xFF)
	return lengthBytes
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
func (fs *FramedStream) ReadMsg(timeout time.Duration) ([]byte, error) {
	if timeout != 0 {
		// set deadline for read
		fs.Trans.SetReadDeadline(time.Now().Add(timeout))
		defer fs.Trans.SetReadDeadline(time.Time{})
	}

	// make sure we have the magic bytes
	err := fs.readMagicBytes()
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			err = ErrTimeout
		}
		return nil, err
	}

	length, err := fs.readLength()
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			err = ErrTimeout
		}
		return nil, err
	}

	msg, err := fs.readPayload(length)
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			err = ErrTimeout
		}
		return nil, err
	}
	return msg, nil
}

// readMagicBytes blocks until it receives the magic bytes 0xFE, 0xED from transport stream.
func (fs *FramedStream) readMagicBytes() error {
	// block until we get the magic bytes
	for {
		for i := 0; i < len(magicBytes); i++ {
			b := make([]byte, 1)
			_, err := fs.Trans.Read(b)
			if err != nil {
				return err
			}
			if b[0] != magicBytes[i] {
				i = 0
				continue
			}
		}
		return nil
	}
}

// readLength reads 4 bytes from transport stream and returns the length as uint of the message.
func (fs *FramedStream) readLength() (uint, error) {
	lengthBytes, err := fs.readAll(4)
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
	payload, err := fs.readAll(length)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// Close closes the transport stream
func (fs *FramedStream) Close() error {
	return fs.Trans.Close()
}

func (fs *FramedStream) readAll(length uint) ([]byte, error) {
	payload := []byte{}
	received := 0
	for {
		chunk := make([]byte, length-uint(received))
		n, err := fs.Trans.Read(chunk)
		log.Debugf("fstream readAll: Read returned n=%d err=%v\n", n, err)
		if err != nil {
			return nil, err
		}
		received += n
		payload = append(payload, chunk[:n]...)
		if received >= int(length) {
			return payload, nil
		}
	}
}
