# go-template

A Go application with RESTful API template.

## Setup

This application requires Go 1.18 version.

To install for development purposes execute these commands:

```shell
$ git clone git@github.com:exepirit/go-template.git
$ cd go-template/
$ go mod download
```

## Build

### Terminal

```shell
$ go build ./cmd/gateway
```

## Configuration

Application can be configured in two ways: via `config.yaml` file in the current working directory, or environment variables.

| Parameter | Env | Description |
| --- | --- | --- |
| `Debug` | `DEBUG` | Debugging mode. In this mode lots of debugging information will be printed on stdout. |
| `ListenAddress` | `LISTENADDRESS` | Host and port, that the HTTP server will listen on. |

## Contribution

Simply create an issue or a pull request.