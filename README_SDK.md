# io4edge-client-go
go client sdk for io4edge

io4edge devices are intelligent I/O devices invented by [Ci4Rail](https://www.ci4rail.com), connected to the host via network.

This package provides a Go API to
* manage those devices, such as:
	* Identify the currently running firmware
	* Load new firmware
	* Identify HW (name, revision, serial number)
	* Program HW identification
	* Set and get persistent parameter

* make use of the function blocks, such as
	* [Analog In TypeA](analogintypea) - IOU01, MIO01
	* [Binary IO TypeA](binaryiotypea) - IOU01, MIO01
	* [Binary IO TypeB](binaryiotypeb) - IOU06
	* [Binary IO TypeC](binaryiotypec) - IOU07
	* [CAN Layer2](canl2) - IOU04, MIO04, IOU03, MIO03, IOU06
	* [Motion Sensor](motionsensor) - CPU01UC
	* [MVB Sniffer](mvbsniffer) - IOU03, MIO03
	* [Template Module](templatemodule)

# Documentation

https://pkg.go.dev/github.com/ci4rail/io4edge-client-go


## System Dependencies

This sdk uses the Avahi Go package to browse for mdns services, which provides bindings for DBus interfaces exposed by the Avahi daemon.

## Installation

```bash
$ go get github.com/ci4rail/io4edge-client-go
```

## Examples for Management of io4edge Devices

### Indentify currently running firmware

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ci4rail/io4edge-client-go/core"
)

func main() {
	const timeout = 5 * time.Second

	address := "192.168.7.1:9999"

	// Create a client object to work with the io4edge device at <address>
	c, err := core.NewClientFromSocketAddress(address)
	if err != nil {
		log.Fatalf("Failed to create basefunc client: %v\n", err)
	}

	// Get the active firmware version from the device
	fwName, fwVersion, err := c.IdentifyFirmware(timeout)
	if err != nil {
		log.Fatalf("Failed to identify firmware: %v\n", err)
	}

	fmt.Printf("Firmware name: %s, Version %s\n", fwName, fwVersion)
}
```

### Load New Firmware

The following example loads a new firmware contained in a firmware package.

A firmware package is a tar file, ending with `.fwpkg` containing the firmware binary and a manifest.json. See [this example](pkg/io4edge/fwpkg/testdata/t1.fwpkg)

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ci4rail/io4edge-client-go/core"
)

func main() {
	const timeout = 5 * time.Second
	const chunkSize = 1024

	address := "192.168.7.1:9999"
	file := "myfirmware.fwpkg"

	// Create a client object to work with the io4edge device at <address>
	c, err := core.NewClientFromSocketAddress(address)
	if err != nil {
		log.Fatalf("Failed to create basefunc client: %v\n", err)
	}

	// Load the firmware package into the device
	// Loading happens in chunks of <chunkSize>. 1024 should work with each device
	// <timeout> is the time to wait for responses from device
	err = c.LoadFirmware(file, chunkSize, timeout)
	if err != nil {
		log.Fatalf("Failed to load firmware package: %v\n", err)
	}
}
```

### Connect to a Target with service name

The following example shows how to connect with a service. The service address consists of <instance_name>.<service_name>.<protocol>.

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ci4rail/io4edge-client-go/core"
)

func main() {
	const timeoutService = 1 * time.Second
	const timeoutCmd = 5 * time.Second

	address := "iou04-usb-ext1._io4edge-core._tcp"

	// Create a client object to work with the io4edge device with the service <address>
	c, err := core.NewClientFromSocketAddress(address, timeoutService)
	if err != nil {
		log.Fatalf("Failed to create basefunc client: %v\n", err)
	}

	// Get the active firmware version from the device
	fwName, fwVersion, err := c.IdentifyFirmware(timeoutCmd)
	if err != nil {
		log.Fatalf("Failed to identify firmware: %v\n", err)
	}

	fmt.Printf("Firmware name: %s, Version %s\n", fwName, fwVersion)
}
```

## Copyright

Copyright © 2021-2022 Ci4Rail GmbH <engineering@ci4rail.com>

io4edge package released unter Apache 2.0 License, see [LICENSE](LICENSE) for details.
