/*
Copyright © 2022 Ci4Rail GmbH <engineering@ci4rail.com>

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

package mvbsniffer

import (
	"fmt"
)

// CommandList holds the list of generator commands generated by AddxxxFrame
type CommandList struct {
	cmds []*command
}

type command struct {
	ram []uint8
}

// ErrorInject defines the parameters for generator error injection
type ErrorInject struct {
	ErrorInA     bool
	ErrorInB     bool
	FullBitError bool  // half otherwise
	Position     uint8 // Positiion = 2.66us + (n * 5.33us)  [n=0..31]
}

// NewCommandList starts a new list of commands
func NewCommandList() *CommandList {
	return &CommandList{cmds: make([]*command, 0)}
}

// AddMasterFrame adds a MVB master frame to the command list
func (c *CommandList) AddMasterFrame(redundantFrameDelayNs int, lineB bool, pauseAfterUs int, fcode uint8, address uint16) error {
	return c.addMasterFrame(redundantFrameDelayNs, lineB, pauseAfterUs, fcode, address, ErrorInject{})
}

// AddMasterFrameWithError adds a MVB master frame with an injected error to the command list
func (c *CommandList) AddMasterFrameWithError(redundantFrameDelayNs int, lineB bool, pauseAfterUs int, fcode uint8, address uint16, injectedError ErrorInject) error {
	return c.addMasterFrame(redundantFrameDelayNs, lineB, pauseAfterUs, fcode, address, injectedError)
}

// AddMasterFrame adds a MVB master frame to the command list
func (c *CommandList) addMasterFrame(redundantFrameDelayNs int, lineB bool, pauseAfterUs int, fcode uint8, address uint16, injectedError ErrorInject) error {
	if fcode > 15 {
		return fmt.Errorf("fcode must be <16")
	}
	if address > 4095 {
		return fmt.Errorf("address must be <4096")
	}

	return c.addFrame(true, redundantFrameDelayNs, lineB, pauseAfterUs, []uint8{(fcode << 4) + uint8((address >> 8)), uint8((address & 0xff))}, injectedError)
}

// AddSlaveFrame adds a MVB slave frame to the command list
func (c *CommandList) AddSlaveFrame(redundantFrameDelayNs int, lineB bool, pauseAfterUs int, data []uint8) error {
	return c.addFrame(false, redundantFrameDelayNs, lineB, pauseAfterUs, data, ErrorInject{})
}

// AddSlaveFrameWithError adds a MVB slave frame with an injected error to the command list
func (c *CommandList) AddSlaveFrameWithError(redundantFrameDelayNs int, lineB bool, pauseAfterUs int, data []uint8, injectedError ErrorInject) error {
	return c.addFrame(false, redundantFrameDelayNs, lineB, pauseAfterUs, data, injectedError)
}

// StartGeneratorString generates from the command list a string that can be sent to the MVB pattern generator to start the pattern
func (c *CommandList) StartGeneratorString(internalLoop bool) string {
	s := "0"
	if internalLoop {
		s += "3"
	} else {
		s += "2"
	}
	s += "4" // clear ram pointer

	c.cmds[len(c.cmds)-1].ram[0] |= 0x1 // set restart bit on last command

	for _, cmd := range c.cmds {
		for _, b := range cmd.ram {
			s += string(rune(64 + ((b & 0xf0) >> 4)))
			s += string(rune(96 + (b & 0xf)))
		}
	}
	s += "1" // start generator
	return s
}

// DumpCommandBytes prints the list of commands (before encoding)
func (c *CommandList) DumpCommandBytes() {
	for _, cmd := range c.cmds {
		fmt.Printf("Command: ")
		for _, b := range cmd.ram {
			fmt.Printf("%02x ", b)
		}
		fmt.Printf("\n")
	}
}

// StopGeneratorString generates a string to stop the MVB pattern generator
func (c *Client) StopGeneratorString() string {
	return "0"
}

// SelectExternalInputString generates a string to switch the MVB receiver to the external signal
func (c *Client) SelectExternalInputString() string {
	return "2"
}

func (c *CommandList) addFrame(isMaster bool, redundantFrameDelayNs int, lineB bool, pauseAfterUs int, data []uint8, injectedError ErrorInject) error {
	cmd := &command{ram: make([]uint8, 3)}

	b, err := cmdByte0(isMaster, len(data), redundantFrameDelayNs, lineB)
	if err != nil {
		return err
	}
	cmd.ram[0] = b
	b, err = cmdByte1(injectedError)
	if err != nil {
		return err
	}
	cmd.ram[1] = b

	b, err = cmdByte2(pauseAfterUs)
	if err != nil {
		return err
	}
	cmd.ram[2] = b
	cmd.ram = append(cmd.ram, data...)

	c.cmds = append(c.cmds, cmd)
	return nil
}

func cmdByte0(isMaster bool, frameLen int, redundantFrameDelayNs int, lineB bool) (uint8, error) {
	b := uint8(0)
	if isMaster {
		b |= 0x80
	}
	switch frameLen {
	case 2:
	case 4:
		b |= 0x10
	case 8:
		b |= 0x20
	case 16:
		b |= 0x30
	case 32:
		b |= 0x40
	default:
		return 0, fmt.Errorf("illegal framelen %d", frameLen)
	}
	if lineB {
		b |= 0x8
	}
	switch redundantFrameDelayNs {
	case 0:
	case 7500:
		b |= 0x2
	case 8500:
		b |= 0x6
	default:
		return 0, fmt.Errorf("illegal redundantFrameDelayNs %d", redundantFrameDelayNs)
	}
	return b, nil
}

func cmdByte1(e ErrorInject) (uint8, error) {
	b := uint8(0)
	if e.Position > 31 {
		return 0, fmt.Errorf("illegal error position %d", e.Position)
	}
	if e.ErrorInB {
		b |= 0x80
	}
	if e.ErrorInA {
		b |= 0x40
	}
	if e.FullBitError {
		b |= 0x20
	}
	b |= e.Position
	return b, nil
}

func cmdByte2(pauseAfterUs int) (uint8, error) {
	if pauseAfterUs <= 127 {
		return uint8(pauseAfterUs), nil
	}
	if (pauseAfterUs%100) != 0 || (pauseAfterUs > (100 * 127)) {
		return 0, fmt.Errorf("illegal pauseAfterUs %d", pauseAfterUs)
	}
	return uint8(pauseAfterUs/100) | 0x80, nil
}