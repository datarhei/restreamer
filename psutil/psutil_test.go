package psutil

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getUtil(path string) *util {
	u, _ := New(path)

	return u.(*util)
}

func TestCgroup2Limited(t *testing.T) {
	u := getUtil("./fixtures/cgroup2-limited")

	assert.Equal(t, true, u.hasCgroup)
	assert.Equal(t, 2, u.cgroupType)

	assert.Equal(t, uint64(1000000000), u.cpuLimit)

	stat, err := u.cgroupCPUTimes(2)

	assert.Nil(t, err)
	assert.Equal(t, float64(111275261000), stat.system)

	mem, err := u.cgroupVirtualMemory(2)

	assert.Nil(t, err)
	assert.Equal(t, uint64(524288000), mem.Total)
	assert.Equal(t, uint64(43745280), mem.Used)
}

func TestCgroup2(t *testing.T) {
	u := getUtil("./fixtures/cgroup2")

	assert.Equal(t, true, u.hasCgroup)
	assert.Equal(t, 2, u.cgroupType)

	assert.Equal(t, uint64(0), u.cpuLimit)

	stat, err := u.cgroupCPUTimes(2)

	assert.Nil(t, err)
	assert.Equal(t, float64(97868879000), stat.system)

	mem, err := u.cgroupVirtualMemory(2)

	assert.Nil(t, err)
	assert.Equal(t, uint64(math.MaxUint64), mem.Total)
	assert.Equal(t, uint64(41603072), mem.Used)
}

func TestCgroupLimited(t *testing.T) {
	u := getUtil("./fixtures/cgroup-limited")

	assert.Equal(t, true, u.hasCgroup)
	assert.Equal(t, 1, u.cgroupType)

	assert.Equal(t, uint64(300000000), u.cpuLimit)

	stat, err := u.cgroupCPUTimes(1)

	assert.Nil(t, err)
	assert.Equal(t, float64(5487224113), stat.system)

	mem, err := u.cgroupVirtualMemory(1)

	assert.Nil(t, err)
	assert.Equal(t, uint64(536870912), mem.Total)
	assert.Equal(t, uint64(34197504), mem.Used)
}

func TestCgroup(t *testing.T) {
	u := getUtil("./fixtures/cgroup")

	assert.Equal(t, true, u.hasCgroup)
	assert.Equal(t, 1, u.cgroupType)

	assert.Equal(t, uint64(0), u.cpuLimit)

	stat, err := u.cgroupCPUTimes(1)

	assert.Nil(t, err)
	assert.Equal(t, float64(809843976), stat.system)

	mem, err := u.cgroupVirtualMemory(1)

	assert.Nil(t, err)
	assert.Equal(t, uint64(9223372036854771712), mem.Total)
	assert.Equal(t, uint64(34070528), mem.Used)
}
