package net

import (
	"errors"
	"fmt"
	"sync"
)

// The Portranger interface allows to get an available port from a pool and to put
// it back after use for later re-use.
type Portranger interface {
	// Get a new port from the pool. The port is from the range as providied at initialization.
	// If no more ports are available a negative port is returned and error is not nil.
	Get() (int, error)

	// Put a port back in the pool. It will be silently ignored if a port has already been returned back
	// to the pool or if the returned port is not in the range.
	Put(int)
}

type portrange struct {
	// Minimal port number
	min int

	// Array to store which ports are used. An
	// unused port is false.
	ports []bool

	// Smallest index in the ports array that
	// is an unused port.
	minUnused int

	lock sync.Mutex
}

// NewPortrange returns a new instance of a Portranger implementation. A minimal and
// maximal port number have to be provided for a valid port range. If the provided
// port range is invalid, nil and an error is retuned.
func NewPortrange(min, max int) (Portranger, error) {
	if max <= min {
		return nil, fmt.Errorf("invalid port range")
	}

	if min <= 0 {
		min = 1
	}

	if max > 65535 {
		max = 65535
	}

	r := &portrange{
		min:       min,
		minUnused: 0,
	}

	r.ports = make([]bool, max-min+1)

	return r, nil
}

func (r *portrange) Get() (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.minUnused == -1 {
		return -1, fmt.Errorf("no more ports available from range [%d,%d]", r.min, r.min+len(r.ports)-1)
	}

	// Calculate new port and mark as used
	var port int = r.min + r.minUnused

	r.ports[r.minUnused] = true

	// Find next unused index
	var minUnused int = -1

	for i := range r.ports {
		if !r.ports[i] {
			minUnused = i
			break
		}
	}

	r.minUnused = minUnused

	return port, nil
}

func (r *portrange) Put(port int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Check if the returned port is in our range
	if port < r.min || port > r.min+len(r.ports)-1 {
		return
	}

	// Translate to index
	port -= r.min

	r.ports[port] = false

	// Adjust the smallest index of the ports array that is unused
	if port < r.minUnused || r.minUnused == -1 {
		r.minUnused = port
	}
}

var ErrNoPortrangerProvided = errors.New("no portranger provided")

type dummy struct{}

func NewDummyPortrange() Portranger {
	return &dummy{}
}

func (d *dummy) Get() (int, error) {
	return 0, ErrNoPortrangerProvided
}

func (d *dummy) Put(port int) {}
