# io4edge-cli

[![.github/workflows/cli-release.yaml](https://github.com/ci4rail/io4edge-client-go/v2/actions/workflows/cli-release.yaml/badge.svg)](https://github.com/ci4rail/io4edge-client-go/v2/actions/workflows/cli-release.yaml)

Command line tool to manage io4edge devices.

## Description

Io4edge devices are intelligent I/O devices invented by Ci4Rail, connected to the host via network.

There are two flavours of io4edge devices:
- Original (O): Using TCP and protobuf for management functions
- New (N): Using HTTPS REST API for management functions

Some functions are only available for the original devices, some only for the new devices, and some for both.

The `io4edge-cli` tool is intended to run on the host machine. Via the command line tool, you can:
* Scan for io4edge devices in the network (O)
* Identify the currently running firmware (O/N)
* Load new firmware (O/N)
* Identify HW (name, revision, serial number) (O/N)
* Program HW identification (O/N)
* Restart device (O/N)
* Set device id (O/N)
* Set/Get persistent parameter (O/N)
* Read a partition from the device (O)
* Get reset reason (O)
* Load/Get parameterset (N)

## Examples

The device can be addressed either with its device-ID or with its IP address and port (see examples).
The device-ID is either a preset persistent parameter in the non-volatile storage of the device (which can be set with the program-devid command) or consisting of the article name and the serial number of the HW inventory of the device (for example `S101-IOU04-70a3b920-7eb7-434e-b20d-6d0a12618ffe`). If both are not programmed, the device is still addressable with it's IP address and the port of the io4edge core (for example `192.168.200.1:9999`).
If the device-ID is not known, it can be found out with `io4edge-cli scan`.

### Scan for available devices
```bash
$ io4edge-cli scan
DEVICE ID               IP              HARDWARE        SERIAL
S101-CPU01UC            192.168.200.1   S101-CPU01UC    (Unknown)
S101-IOU01-USB-EXT-1    192.168.201.1   S101-IOU01      9999
S101-IOU03-USB-EXT-3    192.168.203.1   S101-IOU03      7
```

### Scan for available services/functions
```bash
$ io4edge-cli scan -f
DEVICE ID               SERVICE TYPE                    SERVICE NAME                            IP:PORT
S101-CPU01UC            _io4edge-core._tcp              S101-CPU01UC                            192.168.200.1:9999
                        _ttynvt._tcp                    S101-CPU01UC-com                        192.168.200.1:10000
                        _io4edge_motionSensor._tcp      S101-CPU01UC-accel                      192.168.200.1:10001
S101-IOU01-USB-EXT-1    _io4edge-core._tcp              S101-IOU01-USB-EXT-1                    192.168.201.1:9999
                        _io4edge_analogInTypeA._tcp     S101-IOU01-USB-EXT-1-analogInTypeA1     192.168.201.1:10000
                        _io4edge_analogInTypeA._tcp     S101-IOU01-USB-EXT-1-analogInTypeA2     192.168.201.1:10001
                        _io4edge_binaryIoTypeA._tcp     S101-IOU01-USB-EXT-1-binaryIoTypeA      192.168.201.1:10002
S101-IOU03-USB-EXT-3    _io4edge-core._tcp              S101-IOU03-USB-EXT-3                    192.168.203.1:9999
                        _io4edge_mvbSniffer._tcp        S101-IOU03-USB-EXT-3-mvbSniffer         192.168.203.1:10000
                        _io4edge_canL2._tcp             S101-IOU03-USB-EXT-3-can                192.168.203.1:10001
```

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

### Load firmware from a firmware package
A firmware package contains the firmware binary and a manifest file. The io4edge-cli checks if the firmware is suitable for the device before loading it.

```bash
$ io4edge-cli -d S101-CPU01UC-USB-IO-CTRL load-firmware fw-cpu01uc-tty_accdl-1.1.0.beta1.fwpkg
Reconnect to restarted device
...
```

### Load raw firmware
"Raw Firmware" means: Load a firmware binary that is not embedded in a firmware package file. In this case, the HW/SW compatibility check is NOT performed! Use with caution! You may brick your device!

```bash
$ io4edge-cli -d S101-IOU04-USB-EXT-1 load-raw-firmware build/fw_esp_io4edge_default.bin
Reconnect to restarted device
Reading back firmware id
Firmware name: fw_esp_io4edge_default, Version 1f3f2a2-dirty
```

### Program HW inventory
```bash
$ io4edge-cli -d S101-IOU04-USB-EXT-1 program-hwid S101-IOU04 2 70a3b920-7eb7-434e-b20d-6d0a12618ffe
Success. Read back programmed ID
Hardware name: S101-IOU04, rev: 2, serial: 70a3b920-7eb7-434e-b20d-6d0a12618ffe
```

### Set device-ID
```bash
$ io4edge-cli -i 192.168.201.1:9999 program-devid S101-IOU04-USB-EXT-1
Device id was set to S101-IOU04-USB-EXT-1
Restart of the device required to apply the new device id.
```
Or
```bash
$ io4edge-cli -d S101-IOU04-70a3b920-7eb7-434e-b20d-6d0a12618ffe program-devid S101-IOU04-USB-EXT-1
Device id was set to S101-IOU04-USB-EXT-1
Restart of the device required to apply the new device id.
```

All other commands also accept both addresses (device-ID or IP address with port).

### Set persistent parameter
An io4edge device can provide persistent parameters in its non-volatile storage (NVS). These parameters can be set by the io4edge-cli tool with the set-parameter command.

Attention: The device's firmware must have already reserved a place in the NVS for the parameter to set, which means that the set-parameter command cannot create new parameters, only set or change existing ones. Non-existing parameters return an error `ILLEGAL_PARAMETER`.

```bash
$ io4edge-cli -d S103-MIO04-1 set-parameter wifi-ssid Ci4Rail-Guest
Parameter wifi-ssid was set to Ci4Rail-Guest
```

#### Set persistent parameter from file
Multiple parameters can be set at once by providing a yaml file with the parameters to set.
Example file:
```bash
$ cat wifi.yaml
wifi-ssid: Ci4Rail-Guest
wifi-pw: abc123
```

```bash
$ io4edge-cli -d S103-MIO04-1 set-parameter -f wifi.yaml
```

### Get persistent parameter
An io4edge device can provide persistent parameters in its non-volatile storage (NVS). These parameters can be read by the io4edge-cli tool with the get-parameter command.

Some parameters, such as `wifi-pw` contain secrets and cannot be read by the get-parameter command.

```bash
$ io4edge-cli -d S103-MIO04-1 get-parameter wifi-ssid
Read parameter name: wifi-ssid, value: Ci4Rail-Guest
```
Non-existing parameters return an error `ILLEGAL_PARAMETER`.

### Read a partition from the device
The io4edge-cli tool can read a partition from the device and write it to a file. The partition to read is specified by the partition name. It is especially useful to read the `coredump` partition to perform a post-mortem analysis of a device crash.

Not all versions of the io4edge firmware support this feature. If the firmware does not support the command `UNKNOWN_COMMAND` is returned. If the partition does not exist, the error `ILLEGAL_PARAMETER` is returned.

ESP32 devices support the `coredump` partition only when the device is flashed from scratch with a firmware version that supports this feature. It is not possible to create the `coredump` partition by flashing a firmware package only (as this only updates the application partition, not the partition table).

```bash
$ io4edge-cli -d S103-MIO04-1 read-partition coredump coredump.bin
```

### Get reset reason
The io4edge-cli tool can read the reset reason from the device. The reset reason is returned as a string.

```bash
$ io4edge-cli -d S103-MIO04-1 get-reset-reason
Reset reason: software
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

Releases can be found [here](https://github.com/ci4rail/io4edge-client-go/v2/releases).
