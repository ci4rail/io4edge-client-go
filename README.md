# io4edge-client-go
go client sdk for io4edge

io4edge devices are intelligent I/O devices invented by [Ci4Rail](www.ci4rail.com), connected to the host via network.

This package currently provides a Go API to manage those devices, such as:
* Identify the currently running firmware
* Load new firmware
* Identify HW (name, revision, serial number)
* Program HW identification

Current version uses TCP sockets for communication. May be later extended to further transport protocols such as websockets.

## Installation

```bash
$ go get github.com/ci4rail/io4edge-client-go
```

## Examples

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
	fwID, err := c.IdentifyFirmware(timeout)
	if err != nil {
		log.Fatalf("Failed to identify firmware: %v\n", err)
	}

	fmt.Printf("Firmware name: %s, Version %d.%d.%d\n", fwID.Name, fwID.MajorVersion, fwID.MinorVersion, fwID.PatchVersion)
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

## Copyright

Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

io4edge package released unter Apache 2.0 License, see [LICENSE](LICENSE) for details.
