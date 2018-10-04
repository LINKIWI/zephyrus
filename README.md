# zephyrus

**Zephyrus** is an end-to-end system for reading data from a hardware temperature sensor and emitting them as statsd metrics.

## Architecture

Zephyrus is divided into two main components:

* **server:** Provides abstractions for interacting with the hardware device in the form of formalized gRPC APIs
* **collector:** Daemon, acting as a gRPC client, that streams device data to a statsd server for ingestion and aggregation

The server process is expected to run on the same physical machine as the temperature sensor, while the collector can run on any machine that has access to both the Zephyrus and statsd server. The decision to split the responsibility in a client-server system is motivated by the desire to not couple the device interaction logic with its use case, thus permitting additional clients with unrelated use cases to be developed without needing to change the device abstraction layer.

However, the current implementation makes some assumptions specific to my particular use case:

* The server's device driver assumes the Temper USB sensor (hardware ID `413d:2107`).
* The collector sends statsd metrics in InfluxDB-style, e.g. with tags like `some.metric.name,tag=value`.

## Building

Building requires the Go toolchain, version 1.11 or greater.

```bash
$ make
# This will compile protobuf schemas, followed by the server and collector.
```

## Running

Run `./bin/server --help` and `./bin/collector --help` for usage instructions. The two services can run on the same machine or different machines (as long as they are properly networked and the relevant ports are allowed through firewall).

Daemonize by editing `init/zephyrus-server.service` and `init/zephyrus-collector.service` as necessary and installing as a `systemd` service:

```bash
$ cp init/zephyrus-{server,collector}.service /lib/systemd/system/
$ sudo systemctl daemon-reload
$ sudo systemctl enable zephyrus-{server,collector}
```
