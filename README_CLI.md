# io4edge-cli

[![.github/workflows/cli-release.yaml](https://github.com/ci4rail/io4edge-client-go/actions/workflows/cli-release.yaml/badge.svg)](https://github.com/ci4rail/io4edge-client-go/actions/workflows/cli-release.yaml)

Command line tool to manage io4edge devices.

## Description

Io4edge devices are intelligent I/O devices invented by Ci4Rail, connected to the host via network.

The `io4edge-cli` tool is intended to run on the host machine. Via the command line tool, you can:
* Identify the currently running firmware
* Load new firmware
* Identify HW (name, revision, serial number)
* Program HW identification
* Restart device
* Set device id
* Set/Get persistent parameter

## Examples

The device can be addressed either with its device id or with its ip address and port (see examples).
The device id is either a preset persistent parameter in the non volatile storage of the device (which can be set with the program-devid command) or consisting of the article name and the serial number of the hw inventory of the device (for example `S101-IOU04-70a3b920-7eb7-434e-b20d-6d0a12618ffe`). If both are not programmed, the device is still addressable with its ip address and the port of the io4edge core (for example `192.168.200.1:9999`).
If the device id is not known, it can be find out with a mDNS browser, such as `avahi-browse`.
List the services named `_io4edge-core._tcp` and look up for the instance name:

```shell
$ avahi-browse -t _io4edge-core._tcp
+ usb_ext_2 IPv4 S101-IOU04-USB-EXT-2     _io4edge-core._tcp   local
```

In this example the device id is `S101-IOU04-USB-EXT-2`.

### Identify currently running firmware:
```bash
$ io4edge-cli -d S101-IOU04-USB-EXT-1 fw
Firmware name: fw_esp_io4edge_default, Version 1f3f2a2-dirty
```

### Get hardware inventory information:
```bash
$ io4edge-cli -d S101-IOU04-USB-EXT-1 hw
Hardware name: S101-IOU04, rev: 2, serial: 70a3b920-7eb7-434e-b20d-6d0a12618ffe
```

### Load firmware from a firmware package:
A firmware package contains the firmware binary and a manifest file. The io4edge-cli checks if the firmware is suitable for the device before loading it.

```bash
$ io4edge-cli -d S101-CPU01UC-USB-IO-CTRL load-firmware fw-cpu01uc-tty_accdl-1.1.0.beta1.fwpkg
Reconnect to restarted device
...
```

### Load raw firmware:
"Raw Firmware" means: Load a firmware binary that is not embedded in a firmware package file. In this case, the HW/SW compatibility check is NOT performed!

```bash
$ io4edge-cli -d S101-IOU04-USB-EXT-1 load-raw-firmware build/fw_esp_io4edge_default.bin
Reconnect to restarted device
Reading back firmware id
Firmware name: fw_esp_io4edge_default, Version 1f3f2a2-dirty
```

### Program HW inventory:
```bash
$ io4edge-cli -d S101-IOU04-USB-EXT-1 program-hwid S101-IOU04 2 70a3b920-7eb7-434e-b20d-6d0a12618ffe
Success. Read back programmed ID
Hardware name: S101-IOU04, rev: 2, serial: 70a3b920-7eb7-434e-b20d-6d0a12618ffe
```

### Set device id:
```bash
$ ./io4edge-cli -i 192.168.201.1:9999 program-devid S101-IOU04-USB-EXT-1
Device id was set to S101-IOU04-USB-EXT-1
Restart of the device required to apply the new device id.
```
Or
```bash
$ ./io4edge-cli -d S101-IOU04-70a3b920-7eb7-434e-b20d-6d0a12618ffe program-devid S101-IOU04-USB-EXT-1
Device id was set to S101-IOU04-USB-EXT-1
Restart of the device required to apply the new device id.
```

All other commands also accept both addresses (device id or ip address with port).

### Set persistent parameter:
An io4edge device can provide persistent parameters in its non volatile storage (nvs). These parameter can be set by the io4edge-cli tool with the set-parameter command.

Attention: The devices firmware must already reserved a place in the nvs for the parameter to set, which means that the set-parameter command cannot create new parameters, only set or change existing ones.

```bash
$ ./io4edge-cli -d S101-IOU04-USB-EXT-1 set-parameter wifi-ssid Ci4Rail-Guest
Parameter wifi-ssid was set to Ci4Rail-Guest
```

### Get persistent parameter:
An io4edge device can provide persistent parameters in its non volatile storage (nvs). These parameter can be read by the io4edge-cli tool with the get-parameter command.

```bash
$ ./io4edge-cli -d S101-IOU04-USB-EXT-1 get-parameter wifi-ssid Ci4Rail-Guest
Read parameter name: wifi-ssid, value: Ci4Rail-Guest
```

## Building

### Local Builds

go compiler has to be installed.

```
make
```

## CI

github action `.github/workflows/cli-release.yaml` builds the binaries of io4edge-cli for linux and several CPU architectures. 

Create a Release via the GitHub UI (or gh cli), provide a semantic version name (`vx.y.z`), this will trigger the release action.

Releases can be found [here](https://github.com/ci4rail/io4edge-client-go/releases).

