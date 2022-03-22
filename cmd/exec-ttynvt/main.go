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

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/client"
)

const maxMinorNumbers = 256

type ttynvtInstanceInfo struct {
	cmd   *exec.Cmd
	minor int
}

var ttynvtInstanceMap = make(map[string]ttynvtInstanceInfo)
var programPath string
var major int
var minorMap [maxMinorNumbers]bool

func initMinorMap() {
	for i := range minorMap {
		minorMap[i] = true
	}
}

func getFreeMinor() (minor int, err error) {
	for minor, isFree := range minorMap {
		if isFree {
			return minor + 1, nil
		}
	}
	err = errors.New("no free minor numbers are left")
	return 0, err
}

func setMinorOccupied(minor int) {
	minorMap[minor-1] = false
}

func setMinorFree(minor int) {
	minorMap[minor-1] = true
}

func serviceAdded(s client.ServiceInfo) error {
	var err error
	var instanceInfo ttynvtInstanceInfo
	name := s.GetInstanceName()
	fmt.Println("Added service ", name)
	ipPort := s.GetIPAddressPort()
	instanceInfo.minor, err = getFreeMinor()
	if err != nil {
		log.Errorf("Start ttynvt instance (%s) failed: %v\n", name, err)
		// return nil, that all other ttynvt instances are not terminated
		return nil
	}
	instanceInfo.cmd = exec.Command(programPath, "-f", "-E", "-M", strconv.Itoa(major), "-m", strconv.Itoa(instanceInfo.minor), "-n", name, "-S", ipPort)
	err = instanceInfo.cmd.Start()
	if err != nil {
		log.Errorf("Start ttynvt instance %d (%s) failed: %v\n", instanceInfo.minor, name, err)
		// return nil, that all other ttynvt instances are not terminated
		return nil
	}
	setMinorOccupied(instanceInfo.minor)
	ttynvtInstanceMap[name] = instanceInfo
	return nil
}

func serviceRemoved(s client.ServiceInfo) error {
	name := s.GetInstanceName()
	fmt.Println("Removed service ", name)
	err := ttynvtInstanceMap[name].cmd.Process.Kill()
	if err != nil {
		log.Errorf("Kill ttynvt instance %d (%s) failed: %v\n", ttynvtInstanceMap[name].minor, name, err)
		// return nil, that all other ttynvt instances are not terminated (maybe current instance wasn't ever started)
		return nil
	}
	setMinorFree(ttynvtInstanceMap[name].minor)
	return nil
}

func main() {
	var err error
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <ttynvt-path> <driver-major-number>", os.Args[0])
	}
	programPath = os.Args[1]
	_, err = os.Stat(programPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("error: %s: path not exists!", os.Args[0])
		} else {
			log.Fatalf("error: %v", err)
		}
	}
	major, err = strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("error: driver-major-number must be a number!")
	}
	initMinorMap()
	client.ServiceObserver("_ttynvt._tcp", serviceAdded, serviceRemoved)
}
