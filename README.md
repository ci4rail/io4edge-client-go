# io4edge-client-go
go client sdk for io4edge

io4edge devices are intelligent I/O devices invented by [Ci4Rail](www.ci4rail.com), connected to the host via network.

This package currently provides a Go API to manage those devices, such as:
* Identify the currently running firmware
* Load new firmware
* Identify HW (name, revision, serial number)

Current version uses TCP sockets for communication. May be later extended to further transport protocols such as websockets.

## Installation

```bash
$ go get github.com/ci4rail/io4edge-client-go/pkg/io4edge
```

## Examples

### Indentify currently running firmware

Pass the `device-address` as an IP address with port, e.g. `192.168.7.1:9999`

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ci4rail/io4edge-client-go/pkg/io4edge/basefunc"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: identify <device-address>\n")
	}
	address := os.Args[1]

	c, err := basefunc.NewClientFromSocketAddress(address)
	if err != nil {
		log.Fatalf("Failed to create basefunc client: %v\n", err)
	}

	fwID, err := c.IdentifyFirmware(5 * time.Second)
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

	"github.com/ci4rail/io4edge-client-go/pkg/io4edge/basefunc"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: load  <device-address> <fwpkg>\n")
	}
	address := os.Args[1]
	file := os.Args[2]

	c, err := basefunc.NewClientFromSocketAddress(address)
	if err != nil {
		log.Fatalf("Failed to create basefunc client: %v\n", err)
	}

	err = c.LoadFirmware(file, 1024, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to load firmware package: %v\n", err)
	}
}
```

## Copyright

Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

io4edge package released unter Apache 2.0 License, see [LICENSE](LICENSE) for details.
