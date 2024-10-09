package socketcore

/*
Copyright © 2023 Ci4Rail GmbH <engineering@ci4rail.com>

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

import (
	"fmt"
	"io"
	"net"
	"time"
)

type logReadCloser struct {
	io.Reader
	io.Closer
	dataChan chan []byte
	close    chan bool
	err      error
}

// StreamLogs streams the device log
// streamTimeout is the timeout after which the stream is assumed to be dead and the connection is reestablished
// infoCb is a callback function that is called with status information (such as connection, stream errors)
// It returns an io.ReadCloser that can be used to read the log stream
func (c *Client) StreamLogs(streamTimeout time.Duration, infoCb func(msg string)) (io.ReadCloser, error) {
	// lookup aux port of core client, this is the port where the log stream is available
	protocol, port, err := c.funcClient.FuncInfo.AuxPort()

	if err != nil {
		// no aux port available, assume default port
		protocol = "tcp"
		port = 9998
	}
	if protocol != "tcp" {
		return nil, fmt.Errorf("aux port has unsupported protocol %s", protocol)
	}
	if port == 0 {
		return nil, fmt.Errorf("aux port is not set")
	}
	host, _, err := c.funcClient.FuncInfo.NetAddress()
	if err != nil {
		return nil, err
	}

	l := &logReadCloser{
		dataChan: make(chan []byte),
	}

	// log reader goroutine
	go func(l *logReadCloser) {
		for {
			// open tcp connection to log port
			infoCb(fmt.Sprintf("*** Connecting to %s:%d\n", host, port))
			conn, err := net.Dial(protocol, fmt.Sprintf("%s:%d", host, port))
			if err != nil {
				infoCb(fmt.Sprintf("*** failed to connect to log stream: %s. retrying...\n", err))
				time.Sleep(time.Second)
				continue
			}

			infoCb(fmt.Sprintf("*** Connected to %s:%d\n", host, port))

			if handleConnection(conn, l.dataChan, l.close, streamTimeout, infoCb) {
				conn.Close()
				return // reader close requested
			}
			conn.Close()
		}
	}(l)
	return l, nil
}

func handleConnection(conn net.Conn, dataChan chan []byte, close chan bool, streamTimeout time.Duration, infoCb func(msg string)) bool {
	for {
		buf := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(streamTimeout))

		n, err := conn.Read(buf)
		if err != nil {
			infoCb(fmt.Sprintf("*** Log Read err=%s\n", err))
			return false
		}
		dataChan <- buf[:n]

		select {
		case <-close:
			return true
		default:
		}
	}
}

func (l *logReadCloser) Read(p []byte) (n int, err error) {
	if l.err != nil {
		return 0, l.err
	}
	buf := <-l.dataChan
	copy(p, buf)
	return len(buf), nil
}

func (l *logReadCloser) Close() error {
	l.close <- true
	return nil
}
