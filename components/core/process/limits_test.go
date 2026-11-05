package process

import (
	"sync"
	"testing"
	"time"

	"github.com/datarhei/core/v16/psutil"

	"github.com/stretchr/testify/assert"
)

type psproc struct{}

func (p *psproc) CPUPercent() (*psutil.CPUInfoStat, error) {
	return &psutil.CPUInfoStat{
		System: 50,
		User:   0,
		Idle:   0,
		Other:  0,
	}, nil
}

func (p *psproc) VirtualMemory() (uint64, error) {
	return 197, nil
}

func (p *psproc) Stop() {}

func TestCPULimit(t *testing.T) {
	lock := sync.Mutex{}

	lock.Lock()
	done := false
	lock.Unlock()

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(1)

		l := NewLimiter(LimiterConfig{
			CPU: 42,
			OnLimit: func(float64, uint64) {
				wg.Done()
			},
		})

		l.Start(&psproc{})
		defer l.Stop()

		wg.Wait()

		lock.Lock()
		done = true
		lock.Unlock()
	}()

	assert.Eventually(t, func() bool {
		lock.Lock()
		defer lock.Unlock()

		return done
	}, 2*time.Second, 100*time.Millisecond)
}

func TestCPULimitWaitFor(t *testing.T) {
	lock := sync.Mutex{}

	lock.Lock()
	done := false
	lock.Unlock()

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(1)

		l := NewLimiter(LimiterConfig{
			CPU:     42,
			WaitFor: 3 * time.Second,
			OnLimit: func(float64, uint64) {
				wg.Done()
			},
		})

		l.Start(&psproc{})
		defer l.Stop()

		wg.Wait()

		lock.Lock()
		done = true
		lock.Unlock()
	}()

	assert.Eventually(t, func() bool {
		lock.Lock()
		defer lock.Unlock()

		return done
	}, 10*time.Second, 1*time.Second)
}

func TestMemoryLimit(t *testing.T) {
	lock := sync.Mutex{}

	lock.Lock()
	done := false
	lock.Unlock()

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(1)

		l := NewLimiter(LimiterConfig{
			Memory: 42,
			OnLimit: func(float64, uint64) {
				wg.Done()
			},
		})

		l.Start(&psproc{})
		defer l.Stop()

		wg.Wait()

		lock.Lock()
		done = true
		lock.Unlock()
	}()

	assert.Eventually(t, func() bool {
		lock.Lock()
		defer lock.Unlock()

		return done
	}, 2*time.Second, 100*time.Millisecond)
}

func TestMemoryLimitWaitFor(t *testing.T) {
	lock := sync.Mutex{}

	lock.Lock()
	done := false
	lock.Unlock()

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(1)

		l := NewLimiter(LimiterConfig{
			Memory:  42,
			WaitFor: 3 * time.Second,
			OnLimit: func(float64, uint64) {
				wg.Done()
			},
		})

		l.Start(&psproc{})
		defer l.Stop()

		wg.Wait()

		lock.Lock()
		done = true
		lock.Unlock()
	}()

	assert.Eventually(t, func() bool {
		lock.Lock()
		defer lock.Unlock()

		return done
	}, 10*time.Second, 1*time.Second)
}
