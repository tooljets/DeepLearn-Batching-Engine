
package batching

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"github.com/vmihailenco/msgpack/v4"
	"go.uber.org/zap"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	// IntByteLength defines the byte length of `length` for data
	IntByteLength = 4
	// UUIDLength defines the string bits length
	UUIDLength = 36
	// ErrorIDsKey defines the key for error IDs
	ErrorIDsKey = "error_ids"
)

// String2Bytes structure used in socket communication protocol
type String2Bytes map[string][]byte

// Job wrap the new request as a job waiting to be done by workers
type Job struct {
	id         string
	done       chan bool
	data       []byte // request data
	statusCode int    // HTTP Error Code
	expire     time.Time
}

func newJob(data []byte, timeout time.Duration) *Job {
	return &Job{
		id:         uuid.New().String(),
		done:       make(chan bool, 1),
		data:       data,
		statusCode: 200,
		expire:     time.Now().Add(timeout),
	}
}

// Batching provides HTTP handler and socket communication.
// It generate batch jobs when workers request and send the inference results (or error)
// to the right client.
type Batching struct {
	Address    string // socket file or "{host}:{port}"
	protocol   string // "unix" (Unix domain socket) or "tcp"
	socket     net.Listener
	maxLatency time.Duration // max latency for a batch inference to wait
	batchSize  int           // max batch size for a batch inference
	capacity   int           // the capacity of the batching queue
	timeout    time.Duration // timeout for jobs in the queue
	logger     *zap.Logger
	queue      chan *Job       // job queue
	jobs       map[string]*Job // use job id as the key to find the job
	jobsLock   sync.Mutex      // lock for jobs
}

// NewBatching creates a Batching instance
func NewBatching(address, protocol string, batchSize, capacity int, maxLatency, timeout time.Duration) *Batching {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Cannot create a zap logger")
	}
	if protocol == "unix" {
		// check the socket file (remove if already exists)
		if _, err := os.Stat(address); err == nil {
			logger.Info("Socket file already exists. Try to remove it", zap.String("address", address))
			if err := os.Remove(address); err != nil {
				logger.Fatal("Remove socket file error", zap.Error(err))
			}
		}

	}

	socket, err := net.Listen(protocol, address)
	if err != nil {
		logger.Error("Cannot listen to the socket", zap.Error(err))
		panic("Cannot listen to the socket")
	}

	logger.Info("Listen on socket", zap.String("address", address))
	return &Batching{
		Address:    address,
		protocol:   protocol,
		socket:     socket,
		maxLatency: maxLatency,
		batchSize:  batchSize,
		capacity:   capacity,
		timeout:    timeout,
		logger:     logger,
		queue:      make(chan *Job, capacity),
		jobs:       make(map[string]*Job),
	}
}

// HandleHTTP is the handler for fasthttp
func (b *Batching) HandleHTTP(ctx *fasthttp.RequestCtx) {
	data := ctx.PostBody()
	if len(data) == 0 {