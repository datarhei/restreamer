package psutil

import (
	"context"
	"sync"
	"time"

	psprocess "github.com/shirou/gopsutil/v3/process"
)

type Process interface {
	CPUPercent() (*CPUInfoStat, error)
	VirtualMemory() (uint64, error)
	Stop()
}

type process struct {
	pid       int32
	hasCgroup bool
	cpuLimit  uint64
	ncpu      float64
	proc      *psprocess.Process

	stopTicker context.CancelFunc

	lock             sync.RWMutex
	statCurrent      cpuTimesStat
	statCurrentTime  time.Time
	statPrevious     cpuTimesStat
	statPreviousTime time.Time
}

func (u *util) Process(pid int32) (Process, error) {
	p := &process{
		pid:       pid,
		hasCgroup: u.hasCgroup,
		cpuLimit:  u.cpuLimit,
		ncpu:      u.ncpu,
	}

	proc, err := psprocess.NewProcess(pid)
	if err != nil {
		return nil, err
	}

	p.proc = proc

	ctx, cancel := context.WithCancel(context.Background())
	p.stopTicker = cancel
	go p.tick(ctx)

	return p, nil
}

func NewProcess(pid int32) (Process, error) {
	return DefaultUtil.Process(pid)
}

func (p *process) tick(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			stat := p.collect()

			p.lock.Lock()
			p.statPrevious, p.statCurrent = p.statCurrent, stat
			p.statPreviousTime, p.statCurrentTime = p.statCurrentTime, t
			p.lock.Unlock()
		}
	}
}

func (p *process) collect() cpuTimesStat {
	stat, err := p.cpuTimes()
	if err != nil {
		return cpuTimesStat{
			total: float64(time.Now().Unix()),
			idle:  float64(time.Now().Unix()),
		}
	}

	return *stat
}

func (p *process) Stop() {
	p.stopTicker()
}

func (p *process) cpuTimes() (*cpuTimesStat, error) {
	times, err := p.proc.Times()
	if err != nil {
		return nil, err
	}

	s := &cpuTimesStat{
		total:  cpuTotal(times),
		system: times.System,
		user:   times.User,
	}

	s.other = s.total - s.system - s.user

	return s, nil
}

func (p *process) CPUPercent() (*CPUInfoStat, error) {
	var diff float64

	p.lock.RLock()
	defer p.lock.RUnlock()

	if p.hasCgroup && p.cpuLimit > 0 {
		diff = float64(p.cpuLimit) * (p.statCurrentTime.Sub(p.statPreviousTime)).Seconds() / 1e9
	} else {
		diff = p.statCurrentTime.Sub(p.statPreviousTime).Seconds() * p.ncpu
	}

	s := &CPUInfoStat{
		System: 0,
		User:   0,
		Idle:   0,
		Other:  0,
	}

	if diff <= 0 {
		return s, nil
	}

	s.System = 100 * (p.statCurrent.system - p.statPrevious.system) / diff
	s.User = 100 * (p.statCurrent.user - p.statPrevious.user) / diff
	s.Idle = 100 * (p.statCurrent.idle - p.statPrevious.idle) / diff
	s.Other = 100 * (p.statCurrent.other - p.statPrevious.other) / diff

	return s, nil
}

func (p *process) VirtualMemory() (uint64, error) {
	info, err := p.proc.MemoryInfo()
	if err != nil {
		return 0, err
	}

	return info.RSS, nil
}
