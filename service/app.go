package main

import (
	"flag"
	"fmt"
	"github.com/kemingy/batching"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	address := flag.String("address", "batch.socket", "socket file or host:port")
	protocol := flag.String("protocol", "unix", "unix or tcp")
	batchSize := flag.Int("batch", 32, "max batch size")
	capacity := flag.Int("capacity", 1024, "max jobs in the queue")
	latency := flag.Int("latency", 10, "max latency (millisecond)")
	timeout := flag.Int("timeout", 5000, "timeout for a job (millisecond)")
	host := flag.String("host", "0.0.0.0", "host address")
	port := flag.Int("port", 8080, "service port")
	flag.Parse()
	batch := batching.NewBatching(
		*address,
		*protocol,
		*batchSize,
		*capacity,
		time.Millisecond*time.Dura