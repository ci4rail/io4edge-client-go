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

package transport

import (
	"encoding/binary"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

// FrameHandshake represents a stream with message semantics
type FrameHandshake struct {
	Trans   Transport
	recvSeq recvSeqNum
	sendSeq uint32
}

type recvSeqNum struct {
	lastSeq      uint32
	lastSeqValid bool
}

// NewFrameHandshakeFromTransport creates a message stream from transport t
func NewFrameHandshakeFromTransport(t Transport) *FrameHandshake {
	return &FrameHandshake{
		Trans: t,
		recvSeq: recvSeqNum{
			lastSeq:      0,
			lastSeqValid: false,
		},
		sendSeq: 0,
	}
}

func (fh *FrameHandshake) receiveAck() error {
	fh.Trans.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer fh.Trans.SetReadDeadline(time.Time{})
	for {
		// wait for ack but stop on timeout
		ack := make([]byte, 4)
		n, err := fh.Trans.Read(ack)
		if err != nil {
			return err
		}
		if n < 4 {
			return errors.New("FrameHandshake receiveAck: Ack message too short")
		}

		if binary.LittleEndian.Uint32(ack) == fh.sendSeq {
			break
		} else {
			log.Debugf("Ignoring old ack #%d, expected %d", binary.LittleEndian.Uint32(ack), fh.sendSeq)
		}
	}

	return nil
}

// WriteMsg writes io4edge standard message to the transport stream
func (fh *FrameHandshake) WriteMsg(payload []byte) error {
	// send message via transport
	fh.sendSeq++
	msg := make([]byte, 4+len(payload))
	binary.LittleEndian.PutUint32(msg, fh.sendSeq)
	copy(msg[4:], payload)
	_, err := fh.Trans.Write(msg)
	if err != nil {
		return err
	}

	// wait for ack but stop on timeout
	err = fh.receiveAck()

	return err
}

// ReadMsg reads a io4edge standard message from transport without timeout
func (fh *FrameHandshake) ReadMsg(timeout time.Duration) ([]byte, error) {
	if timeout != 0 {
		// set deadline for read
		fh.Trans.SetReadDeadline(time.Now().Add(timeout))
		defer fh.Trans.SetReadDeadline(time.Time{})
	}

	for {
		msg := make([]byte, 65507)
		n, err := fh.Trans.Read(msg)
		if err != nil {
			return nil, err

		} else if n < 4 {
			err := errors.New("FrameHandshake ReadMsg: Message too short")
			return nil, err
		}

		// log.Printf("FrameHandshake ReadMsg: Read returned n=%d\n", n)

		// Acknowledge the new message
		seq := binary.LittleEndian.Uint32(msg[:4])
		ackBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(ackBuf, seq)
		// log.Printf("FrameHandshake ReadMsg: Sending ACK for message #%d", seq)
		fh.Trans.Write(ackBuf)

		// evaluate sequence number
		if !fh.recvSeq.lastSeqValid || (seq-fh.recvSeq.lastSeq) >= 1 && (seq-fh.recvSeq.lastSeq) <= 100 {
			fh.recvSeq.lastSeq = seq
			fh.recvSeq.lastSeqValid = true
			payload := msg[4:n]
			log.Debugf("FrameHandshake ReadMsg: Received message #%d, len %d\n", seq, len(payload))
			return payload, nil
		}
		log.Debugf("Ignoring DUP message #%d, lastSeq: %d\n", seq, fh.recvSeq.lastSeq)
	}
}

// Close closes the transport stream
func (fh *FrameHandshake) Close() error {
	return fh.Trans.Close()
}
