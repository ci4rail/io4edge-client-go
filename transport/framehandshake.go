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
	"fmt"
	"log"
	"time"
)

// FramedStream represents a stream with message semantics
type FrameHandshake struct {
	Trans   TransportMsg
	recvSeq recvSeqNum
	sendSeq uint32
}

type recvSeqNum struct {
	lastSeq      uint32
	lastSeqValid bool
	// lastMsgReceivedTs time.Time
}

// NewFrameHandshakeFromTransport creates a message stream from transport t
func NewFrameHandshakeFromTransport(t TransportMsg) *FrameHandshake {
	t.SetReadTimeout(time.Millisecond * 100)

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
			fmt.Printf("Ignoring old ack #%d, expected %d", binary.LittleEndian.Uint32(ack), fh.sendSeq)
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
// func (fh *FrameHandshake) ReadMsg() ([]byte, error) {
func (fh *FrameHandshake) ReadMsg() MsgData {
	for {
		msgData := MsgData{
			Payload: nil,
			Err:     nil,
		}
		msg := make([]byte, 65507)
		n, err := fh.Trans.Read(msg)
		log.Printf("FrameHandshake ReadMsg: Read returned n=%d\n", n)
		if err != nil {
			if errors.Is(err, ErrTimeout) {
				continue
			}
			msgData.Err = err
			return msgData
		} else if n < 4 {
			msgData.Err = errors.New("FrameHandshake ReadMsg: Message too short")
			return msgData

		}

		// fh.recvSeq.lastMsgReceivedTs = time.Now()
		seq := binary.LittleEndian.Uint32(msg[:4])

		// Acknowledge the new message
		ackBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(ackBuf, seq)
		fh.Trans.Write(ackBuf)

		// evaluate sequence number
		if !fh.recvSeq.lastSeqValid || (seq-fh.recvSeq.lastSeq) >= 1 && (seq-fh.recvSeq.lastSeq) <= 100 {
			payload := msg[4:n]
			fh.recvSeq.lastSeq = seq
			fh.recvSeq.lastSeqValid = true
			msgData.Payload = payload
			return msgData
		} else {
			fmt.Printf("Ignoring DUP message #%d, lastSeq: %d", seq, fh.recvSeq.lastSeq)
		}
	}
}

// Close closes the transport stream
func (fh *FrameHandshake) Close() error {
	return fh.Trans.Close()
}
