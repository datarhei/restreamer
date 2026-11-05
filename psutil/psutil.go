package psutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var cgroup1Files = []string{
	"cpu/cpu.cfs_quota_us",
	"cpu/cpu.cfs_period_us",
	"cpuacct/cpuacct.usage",
	"memory/memory.limit_in_bytes",
	"memory/memory.usage_in_bytes",
}

var cgroup2Files = []string{
	"cpu.max",
	"cpu.stat",
	"memory.max",
	"memory.current",
}

// https://github.com/netdata/netdata/blob/master/collectors/cgroups.plugin/sys_fs_cgroup.c
// https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/resource_management_guide/sec-cpuacct
// https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/resource_management_guide/sect-cpu-example_usage
// https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html

var DefaultUtil Util

func init() {
	DefaultUtil, _ = New("/sys/fs/cgroup")
}

type MemoryInfoStat struct {
	Total     uint64
	Available uint64
	Used      uint64
}

type CPUInfoStat struct {
	System float64
	User   float64
	Idle   float64
	Other  float64
}

type cpuTimesStat struct {
	total  float64
	system float64
	user   float64
	idle   float64
	other  float64
}

type Util interface {
	Start()
	Stop()
	CPUCounts(logical bool) (float64, error)
	CPUPercent() (*CPUInfoStat, error)
	DiskUsage(path string) (*disk.UsageStat, error)
	VirtualMemory() (*MemoryInfoStat, error)
	NetIOCounters(pernic bool) ([]net.IOCountersStat, error)
	Process(pid int32) (Process, error)
}

type util struct {
	root fs.FS

	cpuLimit   uint64  // Max. allowed CPU time in nanoseconds per second
	ncpu       float64 // Actual available CPUs
	hasCgroup  bool
	cgroupType int

	stopTicker context.CancelFunc
	startOnce  sync.Once
	stopOnce   sync.Once

	lock             sync.RWMutex
	statCurrent      cpuTimesStat
	statCurrentTime  time.Time
	statPrevious     cpuTimesStat
	statPreviousTime time.Time
}

// New returns a new util, it will be started automatically
func New(root string) (Util, error) {
	u := &util{
		root: os.DirFS(root),
	}

	u.cgroupType = u.detectCgroupVersion()
	if u.cgroupType != 0 {
		u.hasCgroup = true
	}

	if u.hasCgroup {
		u.cpuLimit, u.ncpu = u.cgroupCPULimit(u.cgroupType)
	}

	if u.ncpu == 0 {
		var err error
		u.ncpu, err = u.CPUCounts(true)
		if err != nil {
			return nil, err
		}
	}

	u.stopOnce.Do(func() {})

	u.Start()

	return u, nil
}

func (u *util) Start() {
	u.startOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		u.stopTicker = cancel

		go u.tick(ctx, time.Second)
	})
}

func (u *util) Stop() {
	u.stopOnce.Do(func() {
		u.stopTicker()

		u.startOnce = sync.Once{}
	})
}

func (u *util) detectCgroupVersion() int {
	f, err := u.root.Open(".")
	if err != nil {
		// no cgroup available
		return 0
	}

	f.Close()

	for _, file := range cgroup1Files {
		if f, err := u.root.Open(file); err == nil {
			f.Close()
			return 1
		}
	}

	for _, file := range cgroup2Files {
		if f, err := u.root.Open(file); err == nil {
			f.Close()
			return 2
		}
	}

	return 0
}

func (u *util) cgroupCPULimit(version int) (uint64, float64) {
	if version == 1 {
		lines, err := u.readFile("cpu/cpu.cfs_quota_us")
		if err != nil {
			return 0, 0
		}

		quota, err := strconv.ParseFloat(lines[0], 64) // microseconds
		if err != nil {
			return 0, 0
		}

		if quota > 0 {
			lines, err := u.readFile("cpu/cpu.cfs_period_us")
			if err != nil {
				return 0, 0
			}

			period, err := strconv.ParseFloat(lines[0], 64) // microseconds
			if err != nil {
				return 0, 0
			}

			return uint64(1e6/period*quota) * 1e3, quota / period // nanoseconds
		}
	} else if version == 2 {
		lines, err := u.readFile("cpu.max")
		if err != nil {
			return 0, 0
		}

		if strings.HasPrefix(lines[0], "max") {
			return 0, 0
		}

		fields := strings.Split(lines[0], " ")
		if len(fields) != 2 {
			return 0, 0
		}

		quota, err := strconv.ParseFloat(fields[0], 64) // microseconds
		if err != nil {
			return 0, 0
		}

		period, err := strconv.ParseFloat(fields[1], 64) // microseconds
		if err != nil {
			return 0, 0
		}

		return uint64(1e6/period*quota) * 1e3, quota / period // nanoseconds
	}

	return 0, 0
}

func (u *util) tick(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			stat := u.collect()

			u.lock.Lock()
			u.statPrevious, u.statCurrent = u.statCurrent, stat
			u.statPreviousTime, u.statCurrentTime = u.statCurrentTime, t
			u.lock.Unlock()
		}
	}
}

func (u *util) collect() cpuTimesStat {
	stat, err := u.cpuTimes()
	if err != nil {
		return cpuTimesStat{
			total: float64(time.Now().Unix()),
			idle:  float64(time.Now().Unix()),
		}
	}

	return *stat
}

func (u *util) CPUCounts(logical bool) (float64, error) {
	if u.hasCgroup && u.ncpu > 0 {
		return u.ncpu, nil
	}

	ncpu, err := cpu.Counts(logical)
	if err != nil {
		return 0, err
	}

	return float64(ncpu), nil
}

func CPUCounts(logical bool) (float64, error) {
	return DefaultUtil.CPUCounts(logical)
}

func (u *util) cpuTimes() (*cpuTimesStat, error) {
	if u.hasCgroup && u.cpuLimit > 0 {
		if stat, err := u.cgroupCPUTimes(u.cgroupType); err == nil {
			return stat, nil
		}
	}

	times, err := cpu.Times(false)
	if err != nil {
		return nil, err
	}

	if len(times) == 0 {
		return nil, errors.New("cpu.Times() returned an empty slice")
	}

	s := &cpuTimesStat{
		total:  cpuTotal(&times[0]),
		system: times[0].System,
		user:   times[0].User,
		idle:   times[0].Idle,
	}

	s.other = s.total - s.system - s.user - s.idle

	return s, nil
}

func (u *util) CPUPercent() (*CPUInfoStat, error) {
	var total float64

	u.lock.RLock()
	defer u.lock.RUnlock()

	if u.hasCgroup && u.cpuLimit > 0 {
		total = float64(u.cpuLimit) * (u.statCurrentTime.Sub(u.statPreviousTime)).Seconds()
	} else {
		total = (u.statCurrent.total - u.statPrevious.total)
	}

	s := &CPUInfoStat{
		System: 0,
		User:   0,
		Idle:   100,
		Other:  0,
	}

	if total == 0 {
		return s, nil
	}

	s.System = 100 * (u.statCurrent.system - u.statPrevious.system) / total
	s.User = 100 * (u.statCurrent.user - u.statPrevious.user) / total
	s.Idle = 100 * (u.statCurrent.idle - u.statPrevious.idle) / total
	s.Other = 100 * (u.statCurrent.other - u.statPrevious.other) / total

	if u.hasCgroup && u.cpuLimit > 0 {
		s.Idle = 100 - s.User - s.System
	}

	return s, nil
}

func CPUPercent() (*CPUInfoStat, error) {
	return DefaultUtil.CPUPercent()
}

func (u *util) cgroupCPUTimes(version int) (*cpuTimesStat, error) {
	info := &cpuTimesStat{}

	if version == 1 {
		lines, err := u.readFile("cpuacct/cpuacct.usage")
		if err != nil {
			return nil, err
		}

		usage, err := strconv.ParseFloat(lines[0], 64) // nanoseconds
		if err != nil {
			return nil, err
		}

		info.system = usage
	} else if version == 2 {
		lines, err := u.readFile("cpu.stat")
		if err != nil {
			return nil, err
		}

		var usage float64

		if _, err := fmt.Sscanf(lines[0], "usage_usec %f", &usage); err != nil {
			return nil, err
		}

		info.system = usage * 1e3 // convert to nanoseconds
	}

	return info, nil
}

func (u *util) DiskUsage(path string) (*disk.UsageStat, error) {
	return disk.Usage(path)
}

func DiskUsage(path string) (*disk.UsageStat, error) {
	return DefaultUtil.DiskUsage(path)
}

func (u *util) VirtualMemory() (*MemoryInfoStat, error) {
	info, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	if u.hasCgroup {
		if cginfo, err := u.cgroupVirtualMemory(u.cgroupType); err == nil {
			// if total is a huge garbage number, then there are no limits set
			if cginfo.Total <= info.Total {
				return cginfo, nil
			}
		}
	}

	return &MemoryInfoStat{
		Total:     info.Total,
		Available: info.Available,
		Used:      info.Used,
	}, nil
}

func VirtualMemory() (*MemoryInfoStat, error) {
	return DefaultUtil.VirtualMemory()
}

func (u *util) cgroupVirtualMemory(version int) (*MemoryInfoStat, error) {
	info := &MemoryInfoStat{}

	if version == 1 {
		lines, err := u.readFile("memory/memory.limit_in_bytes")
		if err != nil {
			return nil, err
		}

		total, err := strconv.ParseUint(lines[0], 10, 64)
		if err != nil {
			return nil, err
		}

		lines, err = u.readFile("memory/memory.usage_in_bytes")
		if err != nil {
			return nil, err
		}

		used, err := strconv.ParseUint(lines[0], 10, 64)
		if err != nil {
			return nil, err
		}

		info.Total = total
		info.Available = total - used
		info.Used = used
	} else if version == 2 {
		lines, err := u.readFile("memory.max")
		if err != nil {
			return nil, err
		}

		total, err := strconv.ParseUint(lines[0], 10, 64)
		if err != nil {
			total = uint64(math.MaxUint64)
		}

		lines, err = u.readFile("memory.current")
		if err != nil {
			return nil, err
		}

		used, err := strconv.ParseUint(lines[0], 10, 64)
		if err != nil {
			return nil, err
		}

		info.Total = total
		info.Available = total - used
		info.Used = used
	}

	return info, nil
}

func (u *util) NetIOCounters(pernic bool) ([]net.IOCountersStat, error) {
	return net.IOCounters(pernic)
}

func NetIOCounters(pernic bool) ([]net.IOCountersStat, error) {
	return DefaultUtil.NetIOCounters(pernic)
}

func (u *util) readFile(path string) ([]string, error) {
	file, err := u.root.Open(path)
	if err != nil {
		return nil, err
	}

	data := []byte{}
	buf := make([]byte, 20148)
	for {
		n, err := file.Read(buf)

		if n > 0 {
			data = append(data, buf[:n]...)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	return lines, nil
}

func cpuTotal(c *cpu.TimesStat) float64 {
	return c.User + c.System + c.Idle + c.Nice + c.Iowait + c.Irq +
		c.Softirq + c.Steal + c.Guest + c.GuestNice
}
