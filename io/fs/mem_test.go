package fs

import (
	"context"
	"fmt"
	"io"
	gorand "math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/datarhei/core/v16/math/rand"
	"github.com/stretchr/testify/require"
)

func TestMemFromDir(t *testing.T) {
	mem, err := NewMemFilesystemFromDir(".", MemConfig{})
	require.NoError(t, err)

	names := []string{}
	for _, f := range mem.List("/", "/*.go") {
		names = append(names, f.Name())
	}

	require.ElementsMatch(t, []string{
		"/disk.go",
		"/fs_test.go",
		"/fs.go",
		"/mem_storage.go",
		"/mem_test.go",
		"/mem.go",
		"/readonly_test.go",
		"/readonly.go",
		"/s3.go",
		"/sized_test.go",
		"/sized.go",
	}, names)
}

func BenchmarkMemList(b *testing.B) {
	mem, err := NewMemFilesystem(MemConfig{})
	require.NoError(b, err)

	for i := 0; i < 1000; i++ {
		id := rand.StringAlphanumeric(8)
		path := fmt.Sprintf("/%d/%s.dat", i, id)
		mem.WriteFile(path, []byte("foobar"))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mem.List("/", "/5/**")
	}
}

func BenchmarkMemReadFile(b *testing.B) {
	mem, err := NewMemFilesystem(MemConfig{})
	require.NoError(b, err)

	nFiles := 1000

	for i := 0; i < nFiles; i++ {
		path := fmt.Sprintf("/%d.dat", i)
		mem.WriteFile(path, []byte(rand.StringAlphanumeric(2*1024)))
	}

	r := gorand.New(gorand.NewSource(42))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		num := r.Intn(nFiles)
		f := mem.Open("/" + strconv.Itoa(num) + ".dat")
		f.Close()
	}
}

func TestWriteWhileRead(t *testing.T) {
	fs, err := NewMemFilesystem(MemConfig{})
	require.NoError(t, err)

	_, _, err = fs.WriteFile("/foobar", []byte("xxxxx"))
	require.NoError(t, err)

	file := fs.Open("/foobar")
	require.NotNil(t, file)

	_, _, err = fs.WriteFile("/foobar", []byte("yyyyy"))
	require.NoError(t, err)

	data, err := io.ReadAll(file)
	require.NoError(t, err)
	require.Equal(t, []byte("xxxxx"), data)
}

func BenchmarkMemReadFileWhileWriting(b *testing.B) {
	mem, err := NewMemFilesystem(MemConfig{})
	require.NoError(b, err)

	nReaders := 500
	nWriters := 1000
	nFiles := 30

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writerWg := sync.WaitGroup{}

	data := []byte(rand.StringAlphanumeric(2 * 1024))

	for i := 0; i < nWriters; i++ {
		writerWg.Add(1)

		go func(ctx context.Context, from int) {
			for i := 0; i < nFiles; i++ {
				path := fmt.Sprintf("/%d.dat", i+from)
				mem.WriteFile(path, data)
			}

			ticker := time.NewTicker(40 * time.Millisecond)
			defer ticker.Stop()

			writerWg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					num := gorand.Intn(nFiles) + from
					path := fmt.Sprintf("/%d.dat", num)
					mem.WriteFile(path, data)
				}
			}
		}(ctx, i*nFiles)
	}

	// Wait for all writers to be started
	writerWg.Wait()

	b.ResetTimer()

	readerWg := sync.WaitGroup{}

	for i := 0; i < nReaders; i++ {
		readerWg.Add(1)
		go func() {
			defer readerWg.Done()

			for i := 0; i < b.N; i++ {
				num := gorand.Intn(nWriters * nFiles)
				f := mem.Open("/" + strconv.Itoa(num) + ".dat")
				f.Close()
			}
		}()
	}

	readerWg.Wait()
}
