# DeepLearn Batching Engine

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/tooljets/DeepLearn-Batching-Engine)](https://goreportcard.com/report/github.com/tooljets/DeepLearn-Batching-Engine)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/tooljets/DeepLearn-Batching-Engine)
![GitHub Actions](https://github.com/tooljets/DeepLearn-Batching-Engine/workflows/Go/badge.svg)
[![LICENSE](https://img.shields.io/github/license/tooljets/DeepLearn-Batching-Engine.svg)](https://github.com/tooljets/DeepLearn-Batching-Engine/blob/master/LICENSE)

Dynamic Batching Engine for Deep Learning Serving. A tool that implements dynamic batching with batch size and latency factors.

## Warning

**This tool is currently a proof of concept. Refer to [MOSEC](https://github.com/mosecorg/mosec) for production usage.**

## Main Features

* Dynamic batching with control over batch size and latency
* Prevents invalid requests from affecting others in the same batch
* Communicates with workers through Unix domain socket or TCP
* Supports load balancing

Click [here](https://tooljets.github.io/blogs/deep-learning-serving/) to read more about the design concept.

## Configuration Principles

```shell script
go run service/app.go --help
```

```
Usage app:
  -address string
        socket file or host:port (default "batch.socket")
  -batch int
        max batch size (default 32)
  -capacity int
        max jobs in the queue (default 1024)
  -host string
        host address (default "0.0.0.0")
  -latency int
        max latency (millisecond) (default 10)
  -port int
        service port (default 8080)
  -protocol string
        unix or tcp (default "unix")
  -timeo