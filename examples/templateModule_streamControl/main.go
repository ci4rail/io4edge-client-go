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
	"os"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: identify svc <mdns-service-address>  OR  identify ip <ip:port>")
	}
	//addressType := os.Args[1]
	//address := os.Args[2]

	res := &fbv1.Response{
		Status: fbv1.Status_INVALID_PARAMETER,
		Error: &fbv1.Error{
			Error: "something",
		},
	}

	err := functionblock.ResponseErrorNew("wtf", res)
	status := err.(*functionblock.ResponseError)

	fmt.Println(err, status.Status)
}
