package fs

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var ErrNoMinio = errors.New("minio binary not found")

func startMinio(t *testing.T, path string) (*exec.Cmd, error) {
	err := os.MkdirAll(path, 0700)
	require.NoError(t, err)

	minio, err := exec.LookPath("minio")
	if err != nil {
		return nil, ErrNoMinio
	}

	proc := exec.Command(minio, "server", path, "--address", "127.0.0.1:9000")
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout
	err = proc.Start()
	require.NoError(t, err)

	time.Sleep(5 * time.Second)

	return proc, nil
}

func stopMinio(t *testing.T, proc *exec.Cmd) {
	err := proc.Process.Signal(os.Interrupt)
	require.NoError(t, err)

	proc.Wait()
}

func TestFilesystem(t *testing.T) {
	miniopath, err := filepath.Abs("./minio")
	require.NoError(t, err)

	err = os.RemoveAll(miniopath)
	require.NoError(t, err)

	minio, err := startMinio(t, miniopath)
	if err != nil {
		if err != ErrNoMinio {
			require.NoError(t, err)
		}
	}

	os.RemoveAll("./testing/")

	filesystems := map[string]func(string) (Filesystem, error){
		"memfs": func(name string) (Filesystem, error) {
			return NewMemFilesystem(MemConfig{})
		},
		"diskfs": func(name string) (Filesystem, error) {
			return NewRootedDiskFilesystem(RootedDiskConfig{
				Root: "./testing/" + name,
			})
		},
		"s3fs": func(name string) (Filesystem, error) {
			return NewS3Filesystem(S3Config{
				Name:            name,
				Endpoint:        "127.0.0.1:9000",
				AccessKeyID:     "minioadmin",
				SecretAccessKey: "minioadmin",
				Region:          "",
				Bucket:          strings.ToLower(name),
				UseSSL:          false,
				Logger:          nil,
			})
		},
	}

	tests := map[string]func(*testing.T, Filesystem){
		"new":             testNew,
		"metadata":        testMetadata,
		"writeFile":       testWriteFile,
		"writeFileSafe":   testWriteFileSafe,
		"writeFileReader": testWriteFileReader,
		"delete":          testDelete,
		"files":           testFiles,
		"replace":         testReplace,
		"list":            testList,
		"listGlob":        testListGlob,
		"deleteAll":       testDeleteAll,
		"data":            testData,
		"statDir":         testStatDir,
		"mkdirAll":        testMkdirAll,
		"rename":          testRename,
		"renameOverwrite": testRenameOverwrite,
		"copy":            testCopy,
		"symlink":         testSymlink,
		"stat":            testStat,
		"copyOverwrite":   testCopyOverwrite,
		"symlinkErrors":   testSymlinkErrors,
		"symlinkOpenStat": testSymlinkOpenStat,
		"open":            testOpen,
	}

	for fsname, fs := range filesystems {
		for name, test := range tests {
			t.Run(fsname+"-"+name, func(t *testing.T) {
				if fsname == "s3fs" && minio == nil {
					t.Skip("minio server not available")
				}
				filesystem, err := fs(name)
				require.NoError(t, err)
				test(t, filesystem)
			})
		}
	}

	os.RemoveAll("./testing/")

	if minio != nil {
		stopMinio(t, minio)
	}

	os.RemoveAll(miniopath)
}

func testNew(t *testing.T, fs Filesystem) {
	cur, max := fs.Size()

	require.Equal(t, int64(0), cur, "current size")
	require.Equal(t, int64(-1), max, "max size")

	cur = fs.Files()

	require.Equal(t, int64(0), cur, "number of files")
}

func testMetadata(t *testing.T, fs Filesystem) {
	fs.SetMetadata("foo", "bar")
	require.Equal(t, "bar", fs.Metadata("foo"))
}

func testWriteFile(t *testing.T, fs Filesystem) {
	size, created, err := fs.WriteFile("/foobar", []byte("xxxxx"))

	require.Nil(t, err)
	require.Equal(t, int64(5), size)
	require.Equal(t, true, created)

	cur, max := fs.Size()

	require.Equal(t, int64(5), cur)
	require.Equal(t, int64(-1), max)

	cur = fs.Files()

	require.Equal(t, int64(1), cur)
}

func testWriteFileSafe(t *testing.T, fs Filesystem) {
	size, created, err := fs.WriteFileSafe("/foobar", []byte("xxxxx"))

	require.Nil(t, err)
	require.Equal(t, int64(5), size)
	require.Equal(t, true, created)

	cur, max := fs.Size()

	require.Equal(t, int64(5), cur)
	require.Equal(t, int64(-1), max)

	cur = fs.Files()

	require.Equal(t, int64(1), cur)
}

func testWriteFileReader(t *testing.T, fs Filesystem) {
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
}

func testOpen(t *testing.T, fs Filesystem) {
	file := fs.Open("/foobar")
	require.Nil(t, file)

	_, _, err := fs.WriteFileReader("/foobar", strings.NewReader("xxxxx"))
	require.NoError(t, err)

	file = fs.Open("/foobar")
	require.NotNil(t, file)
	require.Equal(t, "/foobar", file.Name())

	stat, err := file.Stat()
	require.NoError(t, err)
	require.Equal(t, "/foobar", stat.Name())
	require.Equal(t, int64(5), stat.Size())
	require.Equal(t, false, stat.IsDir())
}

func testDelete(t *testing.T, fs Filesystem) {
	size := fs.Remove("/foobar")

	require.Equal(t, int64(-1), size)

	data := strings.NewReader("xxxxx")

	fs.WriteFileReader("/foobar", data)

	size = fs.Remove("/foobar")

	require.Equal(t, int64(5), size)

	cur, max := fs.Size()

	require.Equal(t, int64(0), cur)
	require.Equal(t, int64(-1), max)

	cur = fs.Files()

	require.Equal(t, int64(0), cur)
}

func testFiles(t *testing.T, fs Filesystem) {
	require.Equal(t, int64(0), fs.Files())

	fs.WriteFileReader("/foobar.txt", strings.NewReader("bar"))

	require.Equal(t, int64(1), fs.Files())

	fs.MkdirAll("/path/to/foo", 0777)

	require.Equal(t, int64(1), fs.Files())

	fs.Remove("/foobar.txt")

	require.Equal(t, int64(0), fs.Files())
}

func testReplace(t *testing.T, fs Filesystem) {
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

func testList(t *testing.T, fs Filesystem) {
	fs.WriteFileReader("/foobar1", strings.NewReader("a"))
	fs.WriteFileReader("/foobar2", strings.NewReader("bb"))
	fs.WriteFileReader("/foobar3", strings.NewReader("ccc"))
	fs.WriteFileReader("/foobar4", strings.NewReader("dddd"))
	fs.WriteFileReader("/path/foobar3", strings.NewReader("ccc"))
	fs.WriteFileReader("/path/to/foobar4", strings.NewReader("dddd"))

	cur, max := fs.Size()

	require.Equal(t, int64(17), cur)
	require.Equal(t, int64(-1), max)

	cur = fs.Files()

	require.Equal(t, int64(6), cur)

	getNames := func(files []FileInfo) []string {
		names := []string{}
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}

	files := fs.List("/", "")

	require.Equal(t, 6, len(files))
	require.ElementsMatch(t, []string{"/foobar1", "/foobar2", "/foobar3", "/foobar4", "/path/foobar3", "/path/to/foobar4"}, getNames(files))

	files = fs.List("/path", "")

	require.Equal(t, 2, len(files))
	require.ElementsMatch(t, []string{"/path/foobar3", "/path/to/foobar4"}, getNames(files))
}

func testListGlob(t *testing.T, fs Filesystem) {
	fs.WriteFileReader("/foobar1", strings.NewReader("a"))
	fs.WriteFileReader("/path/foobar2", strings.NewReader("a"))
	fs.WriteFileReader("/path/to/foobar3", strings.NewReader("a"))
	fs.WriteFileReader("/foobar4", strings.NewReader("a"))

	cur := fs.Files()

	require.Equal(t, int64(4), cur)

	getNames := func(files []FileInfo) []string {
		names := []string{}
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}

	files := getNames(fs.List("/", "/foo*"))
	require.Equal(t, 2, len(files))
	require.ElementsMatch(t, []string{"/foobar1", "/foobar4"}, files)

	files = getNames(fs.List("/", "/*bar?"))
	require.Equal(t, 2, len(files))
	require.ElementsMatch(t, []string{"/foobar1", "/foobar4"}, files)

	files = getNames(fs.List("/", "/path/*"))
	require.Equal(t, 1, len(files))
	require.ElementsMatch(t, []string{"/path/foobar2"}, files)

	files = getNames(fs.List("/", "/path/**"))
	require.Equal(t, 2, len(files))
	require.ElementsMatch(t, []string{"/path/foobar2", "/path/to/foobar3"}, files)

	files = getNames(fs.List("/path", "/**"))
	require.Equal(t, 2, len(files))
	require.ElementsMatch(t, []string{"/path/foobar2", "/path/to/foobar3"}, files)
}

func testDeleteAll(t *testing.T, fs Filesystem) {
	if _, ok := fs.(*diskFilesystem); ok {
		return
	}

	fs.WriteFileReader("/foobar1", strings.NewReader("abc"))
	fs.WriteFileReader("/path/foobar2", strings.NewReader("abc"))
	fs.WriteFileReader("/path/to/foobar3", strings.NewReader("abc"))
	fs.WriteFileReader("/foobar4", strings.NewReader("abc"))

	cur := fs.Files()

	require.Equal(t, int64(4), cur)

	size := fs.RemoveAll()
	require.Equal(t, int64(12), size)

	cur = fs.Files()

	require.Equal(t, int64(0), cur)
}

func testData(t *testing.T, fs Filesystem) {
	file := fs.Open("/foobar")
	require.Nil(t, file)

	_, err := fs.ReadFile("/foobar")
	require.Error(t, err)

	data := "gduwotoxqb"

	data1 := strings.NewReader(data)

	_, _, err = fs.WriteFileReader("/foobar", data1)
	require.NoError(t, err)

	file = fs.Open("/foobar")
	require.NotNil(t, file)

	data2 := make([]byte, len(data)+1)
	n, err := file.Read(data2)
	if err != nil {
		if err != io.EOF {
			require.NoError(t, err)
		}
	}

	require.Equal(t, len(data), n)
	require.Equal(t, []byte(data), data2[:n])

	data3, err := fs.ReadFile("/foobar")

	require.NoError(t, err)
	require.Equal(t, []byte(data), data3)
}

func testStatDir(t *testing.T, fs Filesystem) {
	info, err := fs.Stat("/")
	require.NoError(t, err)
	require.NotNil(t, info)
	require.Equal(t, true, info.IsDir())

	data := strings.NewReader("gduwotoxqb")
	fs.WriteFileReader("/these/are/some/directories/foobar", data)

	info, err = fs.Stat("/foobar")
	require.Error(t, err)
	require.Nil(t, info)

	info, err = fs.Stat("/these/are/some/directories/foobar")
	require.NoError(t, err)
	require.Equal(t, "/these/are/some/directories/foobar", info.Name())
	require.Equal(t, int64(10), info.Size())
	require.Equal(t, false, info.IsDir())

	info, err = fs.Stat("/these")
	require.NoError(t, err)
	require.Equal(t, "/these", info.Name())
	require.Equal(t, int64(0), info.Size())
	require.Equal(t, true, info.IsDir())

	info, err = fs.Stat("/these/are/")
	require.NoError(t, err)
	require.Equal(t, "/these/are", info.Name())
	require.Equal(t, int64(0), info.Size())
	require.Equal(t, true, info.IsDir())

	info, err = fs.Stat("/these/are/some")
	require.NoError(t, err)
	require.Equal(t, "/these/are/some", info.Name())
	require.Equal(t, int64(0), info.Size())
	require.Equal(t, true, info.IsDir())

	info, err = fs.Stat("/these/are/some/directories")
	require.NoError(t, err)
	require.Equal(t, "/these/are/some/directories", info.Name())
	require.Equal(t, int64(0), info.Size())
	require.Equal(t, true, info.IsDir())
}

func testMkdirAll(t *testing.T, fs Filesystem) {
	info, err := fs.Stat("/foo/bar/dir")
	require.Error(t, err)
	require.Nil(t, info)

	err = fs.MkdirAll("/foo/bar/dir", 0755)
	require.NoError(t, err)

	err = fs.MkdirAll("/foo/bar", 0755)
	require.NoError(t, err)

	info, err = fs.Stat("/foo/bar/dir")
	require.NoError(t, err)
	require.NotNil(t, info)
	require.Equal(t, int64(0), info.Size())
	require.Equal(t, true, info.IsDir())

	info, err = fs.Stat("/")
	require.NoError(t, err)
	require.NotNil(t, info)
	require.Equal(t, int64(0), info.Size())
	require.Equal(t, true, info.IsDir())

	info, err = fs.Stat("/foo")
	require.NoError(t, err)
	require.NotNil(t, info)
	require.Equal(t, int64(0), info.Size())
	require.Equal(t, true, info.IsDir())

	info, err = fs.Stat("/foo/bar")
	require.NoError(t, err)
	require.NotNil(t, info)
	require.Equal(t, int64(0), info.Size())
	require.Equal(t, true, info.IsDir())

	_, _, err = fs.WriteFileReader("/foobar", strings.NewReader("gduwotoxqb"))
	require.NoError(t, err)

	err = fs.MkdirAll("/foobar", 0755)
	require.Error(t, err)
}

func testRename(t *testing.T, fs Filesystem) {
	err := fs.Rename("/foobar", "/foobaz")
	require.Error(t, err)

	_, err = fs.Stat("/foobar")
	require.Error(t, err)

	_, err = fs.Stat("/foobaz")
	require.Error(t, err)

	_, _, err = fs.WriteFileReader("/foobar", strings.NewReader("gduwotoxqb"))
	require.NoError(t, err)

	_, err = fs.Stat("/foobar")
	require.NoError(t, err)

	err = fs.Rename("/foobar", "/foobaz")
	require.NoError(t, err)

	_, err = fs.Stat("/foobar")
	require.Error(t, err)

	_, err = fs.Stat("/foobaz")
	require.NoError(t, err)
}

func testRenameOverwrite(t *testing.T, fs Filesystem) {
	_, err := fs.Stat("/foobar")
	require.Error(t, err)

	_, err = fs.Stat("/foobaz")
	require.Error(t, err)

	_, _, err = fs.WriteFileReader("/foobar", strings.NewReader("foobar"))
	require.NoError(t, err)

	_, _, err = fs.WriteFileReader("/foobaz", strings.NewReader("foobaz"))
	require.NoError(t, err)

	_, err = fs.Stat("/foobar")
	require.NoError(t, err)

	_, err = fs.Stat("/foobaz")
	require.NoError(t, err)

	err = fs.Rename("/foobar", "/foobaz")
	require.NoError(t, err)

	_, err = fs.Stat("/foobar")
	require.Error(t, err)

	_, err = fs.Stat("/foobaz")
	require.NoError(t, err)

	data, err := fs.ReadFile("/foobaz")
	require.NoError(t, err)
	require.Equal(t, "foobar", string(data))
}

func testSymlink(t *testing.T, fs Filesystem) {
	if _, ok := fs.(*s3Filesystem); ok {
		return
	}

	err := fs.Symlink("/foobar", "/foobaz")
	require.Error(t, err)

	_, _, err = fs.WriteFileReader("/foobar", strings.NewReader("foobar"))
	require.NoError(t, err)

	err = fs.Symlink("/foobar", "/foobaz")
	require.NoError(t, err)

	file := fs.Open("/foobaz")
	require.NotNil(t, file)
	require.Equal(t, "/foobaz", file.Name())

	data := make([]byte, 10)
	n, err := file.Read(data)
	if err != nil {
		if err != io.EOF {
			require.NoError(t, err)
		}
	}
	require.NoError(t, err)
	require.Equal(t, 6, n)
	require.Equal(t, "foobar", string(data[:n]))

	stat, err := fs.Stat("/foobaz")
	require.NoError(t, err)
	require.Equal(t, "/foobaz", stat.Name())
	require.Equal(t, int64(6), stat.Size())
	require.NotEqual(t, 0, int(stat.Mode()&os.ModeSymlink))

	link, ok := stat.IsLink()
	require.Equal(t, "/foobar", link)
	require.Equal(t, true, ok)

	data, err = fs.ReadFile("/foobaz")
	require.NoError(t, err)
	require.Equal(t, "foobar", string(data))
}

func testSymlinkOpenStat(t *testing.T, fs Filesystem) {
	if _, ok := fs.(*s3Filesystem); ok {
		return
	}

	_, _, err := fs.WriteFileReader("/foobar", strings.NewReader("foobar"))
	require.NoError(t, err)

	err = fs.Symlink("/foobar", "/foobaz")
	require.NoError(t, err)

	file := fs.Open("/foobaz")
	require.NotNil(t, file)
	require.Equal(t, "/foobaz", file.Name())

	fstat, err := file.Stat()
	require.NoError(t, err)

	stat, err := fs.Stat("/foobaz")
	require.NoError(t, err)

	require.Equal(t, "/foobaz", fstat.Name())
	require.Equal(t, fstat.Name(), stat.Name())

	require.Equal(t, int64(6), fstat.Size())
	require.Equal(t, fstat.Size(), stat.Size())

	require.NotEqual(t, 0, int(fstat.Mode()&os.ModeSymlink))
	require.Equal(t, fstat.Mode(), stat.Mode())
}

func testStat(t *testing.T, fs Filesystem) {
	_, _, err := fs.WriteFileReader("/foobar", strings.NewReader("foobar"))
	require.NoError(t, err)

	file := fs.Open("/foobar")
	require.NotNil(t, file)

	stat1, err := fs.Stat("/foobar")
	require.NoError(t, err)

	stat2, err := file.Stat()
	require.NoError(t, err)

	require.Equal(t, stat1, stat2)
}

func testCopy(t *testing.T, fs Filesystem) {
	err := fs.Rename("/foobar", "/foobaz")
	require.Error(t, err)

	_, err = fs.Stat("/foobar")
	require.Error(t, err)

	_, err = fs.Stat("/foobaz")
	require.Error(t, err)

	_, _, err = fs.WriteFileReader("/foobar", strings.NewReader("gduwotoxqb"))
	require.NoError(t, err)

	_, err = fs.Stat("/foobar")
	require.NoError(t, err)

	err = fs.Copy("/foobar", "/foobaz")
	require.NoError(t, err)

	_, err = fs.Stat("/foobar")
	require.NoError(t, err)

	_, err = fs.Stat("/foobaz")
	require.NoError(t, err)
}

func testCopyOverwrite(t *testing.T, fs Filesystem) {
	_, err := fs.Stat("/foobar")
	require.Error(t, err)

	_, err = fs.Stat("/foobaz")
	require.Error(t, err)

	_, _, err = fs.WriteFileReader("/foobar", strings.NewReader("foobar"))
	require.NoError(t, err)

	_, _, err = fs.WriteFileReader("/foobaz", strings.NewReader("foobaz"))
	require.NoError(t, err)

	_, err = fs.Stat("/foobar")
	require.NoError(t, err)

	_, err = fs.Stat("/foobaz")
	require.NoError(t, err)

	err = fs.Copy("/foobar", "/foobaz")
	require.NoError(t, err)

	_, err = fs.Stat("/foobar")
	require.NoError(t, err)

	_, err = fs.Stat("/foobaz")
	require.NoError(t, err)

	data, err := fs.ReadFile("/foobaz")
	require.NoError(t, err)
	require.Equal(t, "foobar", string(data))
}

func testSymlinkErrors(t *testing.T, fs Filesystem) {
	if _, ok := fs.(*s3Filesystem); ok {
		return
	}

	err := fs.Symlink("/foobar", "/foobaz")
	require.Error(t, err)

	_, _, err = fs.WriteFileReader("/foobar", strings.NewReader("foobar"))
	require.NoError(t, err)

	_, _, err = fs.WriteFileReader("/foobaz", strings.NewReader("foobaz"))
	require.NoError(t, err)

	err = fs.Symlink("/foobar", "/foobaz")
	require.Error(t, err)

	err = fs.Symlink("/foobar", "/bazfoo")
	require.NoError(t, err)

	err = fs.Symlink("/bazfoo", "/barfoo")
	require.Error(t, err)
}
