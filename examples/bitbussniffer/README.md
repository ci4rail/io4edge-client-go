# Bitbus sniffer examples

Run from the repository root. Endpoints can be mDNS instance names or `IP:port`.
Options must precede positional arguments.

## Dump captured frames

```sh
go run ./examples/bitbussniffer/dumpstream -lowlatency -prepare_sender 192.168.210.1:10001
```

The example configures 375000 baud, captures addresses 0, 1 and 2, keeps frames
with bad CRCs, and discards empty frames. It prints stream sequence numbers,
delivery timestamps, frame timestamps in device microseconds, address, control,
information bytes, and CRC/loss/overrun flags. Change the 32-byte address mask in
the source to capture other addresses; all `0xff` bytes accept every address.

`-lowlatency` and `-prepare_sender` are optional and default to false. Press Ctrl+C
to stop the stream and close the connection (a pending read can take three seconds).

## Send a frame

With the sniffer configured with `-prepare_sender`, use a second terminal:

```sh
go run ./examples/bitbussniffer/send 192.168.210.1:10001 0x01 0x93
go run ./examples/bitbussniffer/send 192.168.210.1:10001 1 0x10 0102a0
```

Address and control accept decimal or `0x` hexadecimal byte values. The optional
information argument contains hexadecimal bytes (quoted whitespace is accepted).
The sender preserves the existing device configuration. Sending requires hardware
with a sender (or loopback enabled) and `prepare_sender` enabled.
