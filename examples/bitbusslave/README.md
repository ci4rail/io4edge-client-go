# Bitbus slave loopback example

Requires one device providing both bitbusSlave and sender-capable bitbusSniffer
function blocks, with firmware supporting the current slave API. No Bitbus wiring
or external master is needed.

Run from the repository root, passing the slave and sniffer endpoints on the same
device, as mDNS instance names or `IP:port`:

```sh
go run ./examples/bitbusslave/loopback 192.168.210.1:10002 192.168.210.1:10001
```

The example enables internal loopback and full duplex at 375000 baud and configures
slave address 1, a 5000 ms application watchdog, and an idle response of `00`.
It performs the DISC/SNRM handshake, queues a slave message every 200 ms when no
message is pending, and sends a master message every 500 ms. Between master
messages it polls every 50 ms. Each slave information frame is acknowledged with
RNR before another message is queued. It checks sequence numbers and prints both
sides' received information fields plus the slave status once per second.

Press Ctrl+C to stop streams and close both clients. An in-flight exchange can
take up to three seconds. Device configuration remains set to loopback; configure
the sniffer again before using the external bus.

Based on the [Python slave loopback example](https://github.com/ci4rail/io4edge-client-python/tree/main/examples/bitbusslave).
