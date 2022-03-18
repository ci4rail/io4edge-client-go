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
	"fmt"
	"os/exec"
	"strconv"

	"github.com/ci4rail/io4edge-client-go/client"
)

var ttynvtInstanceMap = make(map[string]*exec.Cmd)
var major = 199

func serviceAdded(s client.ServiceInfo) error {
	name := s.GetInstanceName()
	fmt.Println("Added service ", name)
	ipPort := s.GetIPAddressPort()
	cmd := exec.Command("./ttynvt", "-D", "7", "-d", "-M", strconv.Itoa(major), "-m", "6", "-n", name, "-S", ipPort)
	err := cmd.Start()
	if err != nil {
		return err
	}
	ttynvtInstanceMap[name] = cmd
	major++
	return nil
}

func serviceRemoved(s client.ServiceInfo) error {
	name := s.GetInstanceName()
	fmt.Println("Removed service ", name)
	ttynvtInstanceMap[name].Process.Kill()
	return nil
}

func main() {
	client.ServiceObserver("_ttynvt._tcp", serviceAdded, serviceRemoved)
}
