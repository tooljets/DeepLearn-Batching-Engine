
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