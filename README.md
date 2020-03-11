# Welcome to Syphon metric proxy 👋

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg?cacheSeconds=2592000)

> Versatile metrics processor, proxy and forwarder

## Install

```bash
go get github.com/butlerx/syphon/cmd/syphon
```

## Usage

```bash
$ syphon -h
NAME:
   syphon - Versatile metrics processor, proxy and forwarder

USAGE:
   syphon [options] COMMAND

VERSION:
   1.0.0

DESCRIPTION:

      syphon is designed to accept and route metrics traffic.
      Metrics can be received from socket, snooped from live traffic or read from file or kafka.
      Metrics can be exportered via file, kafka or udp/tcp


AUTHOR:
   Cian Butler <butlerx@notthe.cloud>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value  Config file to use
   --print         Print default config (default: false)
   --help, -h      show help (default: false)
   --version, -v   print the version (default: false)
```

## Configuration

```toml
[metric]
# Endpoint for store internal carbon metrics.
# Valid values: "" or "local", "tcp://host:port", "udp://host:port"
endpoint = "local"
# Interval of storing internal metrics. Like CARBON_METRIC_INTERVAL
interval = "1m0s"

[logging]
# "stderr", "stdout" can be used as file name
file = "stderr"
# Logging error level. Valid values: "debug", "info", "warn", "error"
level = "info"
# Logging encoding format. Valid values: "mixes", "json", "console"
encoding = "mixed"
encoding-time = "iso8601"
encoding-duration = "seconds"

[file]
enabled = false
path = ""

[prometheus]
listen = ":2006"
enabled = false

[tcp]
listen = ":2003"
enabled = true

[udp]
listen = ":2003"
# Setting mode to promiscuous sets the interface to promiscuously listen
# mode = "promiscuous"
enabled = true

[[uploader.file]]
enabled = true
path = "metrics_recieved.txt"
# RegEx pattern to use to Determine if metric should be sent
pattern = ".*"


[[uploader.udp]]
enabled = true
host = "localhost"
port = 2004
pattern = ".*"

[[uploader.tcp]]
enabled = true
host = "localhost"
port = 2004
pattern = ".*"

# Designed for use with carbon-clickhouse
# https://github.com/lomik/carbon-clickhouse/blob/master/grpc/carbon.proto
[[uploader.grpc]]
enabled = false
host = "localhost"
port = 2005
pattern = ".*"
```

## Run tests

```bash
make test
```

## Author

👤 **Cian Butler**
