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

f