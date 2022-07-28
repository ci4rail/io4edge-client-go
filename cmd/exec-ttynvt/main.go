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

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/client"
)

const maxMinorNumbers = 256

type ttynvtInstanceInfo struct {
	cmd    *exec.Cmd
	minor  int
	ipPort string
}

var ttynvtInstanceMap = make(map[string]*ttynvtInstanceInfo)
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

func ttyName(instanceName string) string {
	return "tty" + instanceName
}

func killCmd(name string, cmd *exec.Cmd) error {
	if cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			log.Warnf("Kill ttynvt instance for %s failed: %v\n", name, err)
			return err
		}
		cmd.Wait()
		return nil
	}
	log.Warnf("ttynvt instance for %s not running", name)
	return nil
}

func delInfo(name string) {
	if info, ok := ttynvtInstanceMap[name]; ok {
		if info.minor != 0 {
			setMinorFree(info.minor)
		}
		delete(ttynvtInstanceMap, name)
	}
}

func serviceAdded(s client.ServiceInfo) error {
	var info *ttynvtInstanceInfo

	fmt.Println("Added service", s.GetInstanceName())

	name := ttyName(s.GetInstanceName())
	ipPort := s.GetIPAddressPort()

	info, ok := ttynvtInstanceMap[name]
	if ok {
		// instance already exists, check if ip or port changed
		if info.ipPort == ipPort {
			log.Infof("no change in ip/port for instance %s", name)
			return nil
		}
		// ip or port changed, kill old instance and start new one
		log.Infof("ip/port changed for instance %s, %s->%s killing old instance", name, info.ipPort, ipPort)
		killCmd(name, info.cmd)
		info.cmd = nil

	} else {
		// instance does not exist. start new instance
		info = &ttynvtInstanceInfo{}
		info.ipPort = ipPort
		minor, err := getFreeMinor()

		if err != nil {
			log.Errorf("No free minor numbers for %s: %v\n", name, err)
			return nil
		}
		info.minor = minor
		setMinorOccupied(info.minor)
		log.Infof("start process for instance %d (%s)", minor, name)
		ttynvtInstanceMap[name] = info
	}
	info.cmd = exec.Command(programPath, "-f", "-E", "-M", strconv.Itoa(major), "-m", strconv.Itoa(info.minor), "-n", name, "-S", ipPort)
	err := info.cmd.Start()
	if err != nil {
		log.Errorf("Start ttynvt instance %d (%s) failed: %v\n", info.minor, name, err)
		delInfo(name)
		return nil
	}

	return nil
}

func serviceRemoved(s client.ServiceInfo) error {
	name := ttyName(s.GetInstanceName())
	fmt.Println("Removed service", s.GetInstanceName())

	info, ok := ttynvtInstanceMap[name]
	if ok {
		log.Infof("Killing ttynvt instance for %s", name)
		killCmd(name, info.cmd)
		delInfo(name)
	} else {
		log.Warnf("ttynvt instance for %s not in map", name)
	}
	return nil
}

func main() {
	var err error

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <ttynvt-program-path> <driver-major-number>\n", os.Args[0])
		os.Exit(1)
	}

	logLevel := flag.String("loglevel", "info", "loglevel (debug, info, warn, error)")
	// parse command line arguments
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
	}

	level, err := log.ParseLevel(*logLevel)

	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	log.SetLevel(level)

	programPath = flag.Arg(0)
	_, err = os.Stat(programPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("error: %s: path not exists!", os.Args[0])
		} else {
			log.Fatalf("error: %v", err)
		}
	}
	major, err = strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.Fatalf("error: driver-major-number must be a number!")
	}
	initMinorMap()
	client.ServiceObserver("_ttynvt._tcp", serviceAdded, serviceRemoved)
}
