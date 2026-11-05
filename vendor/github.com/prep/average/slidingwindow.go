// Package average implements sliding time window.
package average

import (
	"errors"
	"sync"
	"time"
)

// SlidingWindow provides a sliding time window with a custom size and
// granularity to store int64 counters. This can be used to determine the total
// or unweighted mean average of a subset of the window size.
type SlidingWindow struct {
	window      time.Duration
	granularity time.Duration
	samples     []int64
	pos         int
	size        int
	stopOnce    sync.Once
	stopC       chan struct{}
	sync.RWMutex
}

// MustNew returns a new SlidingWindow, but panics if an error occurs.
func MustNew(window, granularity time.Duration) *SlidingWindow {
	sw, err := New(window, granularity)
	if err != nil {
		panic(err.Error())
	}

	return sw
}

// New returns a new SlidingWindow.
func New(window, granularity time.Duration) (*SlidingWindow, error) {
	if window == 0 {
		return nil, errors.New("window cannot be 0")
	}
	if granularity == 0 {
		return nil, errors.New("granularity cannot be 0")
	}
	if window <= granularity || window%granularity != 0 {
		return nil, errors.New("window size has to be a multiplier of the granularity size")
	}

	sw := &SlidingWindow{
		window:      window,
		granularity: granularity,
		samples:     make([]int64, int(window/granularity)),
		stopC:       make(chan struct{}),
	}

	go sw.shifter()
	return sw, nil
}

func (sw *SlidingWindow) shifter() {
	ticker := time.NewTicker(sw.granularity)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sw.Lock()
			if sw.pos = sw.pos + 1; sw.pos >= len(sw.samples) {
				sw.pos = 0
			}
			sw.samples[sw.pos] = 0
			if sw.size < len(sw.samples) {
				sw.size++
			}
			sw.Unlock()

		case <-sw.stopC:
			return
		}
	}
}

// Add increments the value of the current sample.
func (sw *SlidingWindow) Add(v int64) {
	sw.Lock()
	sw.samples[sw.pos] += v
	sw.Unlock()
}

// Average returns the unweighted mean of the specified window.
func (sw *SlidingWindow) Average(window time.Duration) float64 {
	total, sampleCount := sw.Total(window)
	if sampleCount == 0 {
		return 0
	}

	return float64(total) / float64(sampleCount)
}

// Reset the samples in this sliding time window.
func (sw *SlidingWindow) Reset() {
	sw.Lock()
	defer sw.Unlock()

	sw.pos, sw.size = 0, 0
	for i := range sw.samples {
		sw.samples[i] = 0
	}
}

// Stop the shifter of this sliding time window. A stopped SlidingWindow cannot
// be started again.
func (sw *SlidingWindow) Stop() {
	sw.stopOnce.Do(func() {
		sw.stopC <- struct{}{}
	})
}

// Total returns the sum of all values over the specified window, as well as
// the number of samples.
func (sw *SlidingWindow) Total(window time.Duration) (int64, int) {
	if window > sw.window {
		window = sw.window
	}

	sampleCount := int(window / sw.granularity)
	if sampleCount > sw.size {
		sampleCount = sw.size
	}

	sw.RLock()
	defer sw.RUnlock()

	var total int64
	for i := 1; i <= sampleCount; i++ {
		pos := sw.pos - i
		if pos < 0 {
			pos += len(sw.samples)
		}

		total += sw.samples[pos]
	}

	return total, sampleCount
}
