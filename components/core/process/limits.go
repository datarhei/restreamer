package process

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/datarhei/core/v16/psutil"
)

type LimitFunc func(cpu float64, memory uint64)

type LimiterConfig struct {
	CPU     float64       // Max. CPU usage in percent
	Memory  uint64        // Max. memory usage in bytes
	WaitFor time.Duration // Duration one of the limits has to be above the limit until OnLimit gets triggered
	OnLimit LimitFunc     // Function to be triggered if limits are exceeded
}

type Limiter interface {
	// Start starts the limiter with a psutil.Process.
	Start(process psutil.Process) error

	// Stop stops the limiter. The limiter can be reused by calling Start() again
	Stop()

	// Current returns the current CPU and memory values
	Current() (cpu float64, memory uint64)

	// Limits returns the defined CPU and memory limits. Values < 0 means no limit
	Limits() (cpu float64, memory uint64)
}

type limiter struct {
	proc    psutil.Process
	lock    sync.Mutex
	cancel  context.CancelFunc
	onLimit LimitFunc

	cpu              float64
	cpuCurrent       float64
	cpuLast          float64
	cpuLimitSince    time.Time
	memory           uint64
	memoryCurrent    uint64
	memoryLast       uint64
	memoryLimitSince time.Time
	waitFor          time.Duration
}

// NewLimiter returns a new Limiter
func NewLimiter(config LimiterConfig) Limiter {
	l := &limiter{
		cpu:     config.CPU,
		memory:  config.Memory,
		waitFor: config.WaitFor,
		onLimit: config.OnLimit,
	}

	if l.onLimit == nil {
		l.onLimit = func(float64, uint64) {}
	}

	return l
}

func (l *limiter) reset() {
	l.cpuCurrent = 0
	l.cpuLast = 0
	l.memoryCurrent = 0
	l.memoryLast = 0
}

func (l *limiter) Start(process psutil.Process) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.proc != nil {
		return fmt.Errorf("limiter is already running")
	}

	l.reset()

	l.proc = process

	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel

	go l.ticker(ctx)

	return nil
}

func (l *limiter) Stop() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.proc == nil {
		return
	}

	l.cancel()

	l.proc.Stop()
	l.proc = nil

	l.reset()
}

func (l *limiter) ticker(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			l.collect(t)
		}
	}
}

func (l *limiter) collect(t time.Time) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.proc == nil {
		return
	}

	if mstat, err := l.proc.VirtualMemory(); err == nil {
		l.memoryLast, l.memoryCurrent = l.memoryCurrent, mstat
	}

	if cpustat, err := l.proc.CPUPercent(); err == nil {
		l.cpuLast, l.cpuCurrent = l.cpuCurrent, cpustat.System+cpustat.User+cpustat.Other
	}

	isLimitExceeded := false

	if l.cpu > 0 {
		if l.cpuCurrent > l.cpu {
			// Current value is higher than the limit
			if l.cpuLast <= l.cpu {
				// If the previous value is below the limit, then we reached the
				// limit as of now
				l.cpuLimitSince = time.Now()
			}

			if time.Since(l.cpuLimitSince) >= l.waitFor {
				isLimitExceeded = true
			}
		}
	}

	if l.memory > 0 {
		if l.memoryCurrent > l.memory {
			// Current value is higher than the limit
			if l.memoryLast <= l.memory {
				// If the previous value is below the limit, then we reached the
				// limit as of now
				l.memoryLimitSince = time.Now()
			}

			if time.Since(l.memoryLimitSince) >= l.waitFor {
				isLimitExceeded = true
			}
		}
	}

	if isLimitExceeded {
		go l.onLimit(l.cpuCurrent, l.memoryCurrent)
	}
}

func (l *limiter) Current() (cpu float64, memory uint64) {
	l.lock.Lock()
	defer l.lock.Unlock()

	cpu = l.cpuCurrent
	memory = l.memoryCurrent

	return
}

func (l *limiter) Limits() (cpu float64, memory uint64) {
	return l.cpu, l.memory
}
