package fs

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newMemFS() Filesystem {
	mem, _ := NewMemFilesystem(MemConfig{})

	return mem
}

func TestNewSized(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, false)

	cur, max := fs.Size()

	require.Equal(t, int64(0), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(0), cur)
}

func TestSizedResize(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, false)

	cur, max := fs.Size()

	require.Equal(t, int64(0), cur)
	require.Equal(t, int64(10), max)

	err := fs.Resize(20, false)
	require.NoError(t, err)

	cur, max = fs.Size()

	require.Equal(t, int64(0), cur)
	require.Equal(t, int64(20), max)
}

func TestSizedResizeSetPurge(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, false)

	_, _, err := fs.WriteFileReader("/foobar1", strings.NewReader("xxxxx"))
	require.NoError(t, err)

	_, _, err = fs.WriteFileReader("/foobar2", strings.NewReader("xxxxx"))
	require.NoError(t, err)

	_, _, err = fs.WriteFileReader("/foobar3", strings.NewReader("xxxxx"))
	require.Error(t, err)

	fs.Resize(10, true)

	_, _, err = fs.WriteFileReader("/foobar3", strings.NewReader("xxxxx"))
	require.NoError(t, err)

	require.Equal(t, int64(2), fs.Files())
}

func TestSizedResizePurge(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, false)

	cur, max := fs.Size()

	require.Equal(t, int64(0), cur)
	require.Equal(t, int64(10), max)

	fs.WriteFileReader("/foobar", strings.NewReader("xxxxxxxxxx"))

	cur, max = fs.Size()

	require.Equal(t, int64(10), cur)
	require.Equal(t, int64(10), max)

	err := fs.Resize(5, false)
	require.NoError(t, err)

	cur, max = fs.Size()

	require.Equal(t, int64(0), cur)
	require.Equal(t, int64(5), max)
}

func TestSizedWrite(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, false)

	cur, max := fs.Size()

	require.Equal(t, int64(0), cur)
	require.Equal(t, int64(10), max)

	size, created, err := fs.WriteFileReader("/foobar", strings.NewReader("xxxxx"))
	require.NoError(t, err)
	require.Equal(t, int64(5), size)
	require.Equal(t, true, created)

	cur, max = fs.Size()

	require.Equal(t, int64(5), cur)
	require.Equal(t, int64(10), max)

	_, _, err = fs.WriteFile("/foobaz", []byte("xxxxxx"))
	require.Error(t, err)

	_, _, err = fs.WriteFileReader("/foobaz", strings.NewReader("xxxxxx"))
	require.Error(t, err)

	_, _, err = fs.WriteFileSafe("/foobaz", []byte("xxxxxx"))
	require.Error(t, err)
}

func TestSizedReplaceNoPurge(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, false)

	data := strings.NewReader("xxxxx")

	size, created, err := fs.WriteFileReader("/foobar", data)

	require.Nil(t, err)
	require.Equal(t, int64(5), size)
	require.Equal(t, true, created)

	cur, max := fs.Size()

	require.Equal(t, int64(5), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(1), cur)

	data = strings.NewReader("yyy")

	size, created, err = fs.WriteFileReader("/foobar", data)

	require.Nil(t, err)
	require.Equal(t, int64(3), size)
	require.Equal(t, false, created)

	cur, max = fs.Size()

	require.Equal(t, int64(3), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(1), cur)

	data = strings.NewReader("zzzzzzz")

	size, created, err = fs.WriteFileReader("/foobar", data)

	require.Nil(t, err)
	require.Equal(t, int64(7), size)
	require.Equal(t, false, created)

	cur, max = fs.Size()

	require.Equal(t, int64(7), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(1), cur)

	data = strings.NewReader("zzzzzzzz")

	size, created, err = fs.WriteFileReader("/foobar", data)

	require.Nil(t, err)
	require.Equal(t, int64(8), size)
	require.Equal(t, false, created)
}

func TestSizedReplacePurge(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, true)

	data1 := strings.NewReader("xxx")
	data2 := strings.NewReader("yyy")
	data3 := strings.NewReader("zzz")

	fs.WriteFileReader("/foobar1", data1)
	fs.WriteFileReader("/foobar2", data2)
	fs.WriteFileReader("/foobar3", data3)

	cur, max := fs.Size()

	require.Equal(t, int64(9), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(3), cur)

	data4 := strings.NewReader("zzzzz")

	size, _, _ := fs.WriteFileReader("/foobar1", data4)

	require.Equal(t, int64(5), size)

	cur, max = fs.Size()

	require.Equal(t, int64(8), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(2), cur)
}

func TestSizedReplaceUnlimited(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), -1, false)

	data := strings.NewReader("xxxxx")

	size, created, err := fs.WriteFileReader("/foobar", data)

	require.Nil(t, err)
	require.Equal(t, int64(5), size)
	require.Equal(t, true, created)

	cur, max := fs.Size()

	require.Equal(t, int64(5), cur)
	require.Equal(t, int64(-1), max)

	cur = fs.Files()

	require.Equal(t, int64(1), cur)

	data = strings.NewReader("yyy")

	size, created, err = fs.WriteFileReader("/foobar", data)

	require.Nil(t, err)
	require.Equal(t, int64(3), size)
	require.Equal(t, false, created)

	cur, max = fs.Size()

	require.Equal(t, int64(3), cur)
	require.Equal(t, int64(-1), max)

	cur = fs.Files()

	require.Equal(t, int64(1), cur)
}

func TestSizedTooBigNoPurge(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, false)

	data := strings.NewReader("xxxxxyyyyyz")

	size, _, err := fs.WriteFileReader("/foobar", data)
	require.Error(t, err)
	require.Equal(t, int64(-1), size)
}

func TestSizedTooBigPurge(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, true)

	data1 := strings.NewReader("xxxxx")
	data2 := strings.NewReader("yyyyy")

	fs.WriteFileReader("/foobar1", data1)
	fs.WriteFileReader("/foobar2", data2)

	data := strings.NewReader("xxxxxyyyyyz")

	size, _, err := fs.WriteFileReader("/foobar", data)
	require.Error(t, err)
	require.Equal(t, int64(-1), size)

	require.Equal(t, int64(2), fs.Files())
}

func TestSizedFullSpaceNoPurge(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, false)

	data1 := strings.NewReader("xxxxx")
	data2 := strings.NewReader("yyyyy")

	fs.WriteFileReader("/foobar1", data1)
	fs.WriteFileReader("/foobar2", data2)

	cur, max := fs.Size()

	require.Equal(t, int64(10), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(2), cur)

	data3 := strings.NewReader("zzzzz")

	size, _, err := fs.WriteFileReader("/foobar3", data3)
	require.Error(t, err)
	require.Equal(t, int64(-1), size)
}

func TestSizedFullSpacePurge(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, true)

	data1 := strings.NewReader("xxxxx")
	data2 := strings.NewReader("yyyyy")

	fs.WriteFileReader("/foobar1", data1)
	fs.WriteFileReader("/foobar2", data2)

	cur, max := fs.Size()

	require.Equal(t, int64(10), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(2), cur)

	data3 := strings.NewReader("zzzzz")

	size, _, _ := fs.WriteFileReader("/foobar3", data3)

	require.Equal(t, int64(5), size)

	cur, max = fs.Size()

	require.Equal(t, int64(10), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(2), cur)
}

func TestSizedFullSpacePurgeMulti(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, true)

	data1 := strings.NewReader("xxx")
	data2 := strings.NewReader("yyy")
	data3 := strings.NewReader("zzz")

	fs.WriteFileReader("/foobar1", data1)
	fs.WriteFileReader("/foobar2", data2)
	fs.WriteFileReader("/foobar3", data3)

	cur, max := fs.Size()

	require.Equal(t, int64(9), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(3), cur)

	data4 := strings.NewReader("zzzzz")

	size, _, _ := fs.WriteFileReader("/foobar4", data4)

	require.Equal(t, int64(5), size)

	cur, max = fs.Size()

	require.Equal(t, int64(8), cur)
	require.Equal(t, int64(10), max)

	cur = fs.Files()

	require.Equal(t, int64(2), cur)
}

func TestSizedPurgeOrder(t *testing.T) {
	fs, _ := NewSizedFilesystem(newMemFS(), 10, true)

	data1 := strings.NewReader("xxxxx")
	data2 := strings.NewReader("yyyyy")
	data3 := strings.NewReader("zzzzz")

	fs.WriteFileReader("/foobar1", data1)
	time.Sleep(1 * time.Second)
	fs.WriteFileReader("/foobar2", data2)
	time.Sleep(1 * time.Second)
	fs.WriteFileReader("/foobar3", data3)

	file := fs.Open("/foobar1")

	require.Nil(t, file)
}
