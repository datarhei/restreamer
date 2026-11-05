package fs

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/datarhei/core/v16/glob"
	"github.com/datarhei/core/v16/log"
)

// DiskConfig is the config required to create a new disk filesystem.
type DiskConfig struct {
	// For logging, optional
	Logger log.Logger
}

// RootedDiskConfig is the config required to create a new rooted disk filesystem.
type RootedDiskConfig struct {
	// Root is the path this filesystem is rooted to
	Root string

	// For logging, optional
	Logger log.Logger
}

// diskFileInfo implements the FileInfo interface
type diskFileInfo struct {
	root  string
	name  string
	mode  os.FileMode
	finfo os.FileInfo
}

func (fi *diskFileInfo) Name() string {
	return fi.name
}

func (fi *diskFileInfo) Size() int64 {
	if fi.finfo.IsDir() {
		return 0
	}

	return fi.finfo.Size()
}

func (fi *diskFileInfo) Mode() fs.FileMode {
	return fi.mode
}

func (fi *diskFileInfo) ModTime() time.Time {
	return fi.finfo.ModTime()
}

func (fi *diskFileInfo) IsLink() (string, bool) {
	if fi.mode&os.ModeSymlink == 0 {
		return fi.name, false
	}

	path, err := os.Readlink(filepath.Join(fi.root, fi.name))
	if err != nil {
		return fi.name, false
	}

	if !strings.HasPrefix(path, fi.root) {
		return fi.name, false
	}

	name := strings.TrimPrefix(path, fi.root)

	if name[0] != os.PathSeparator {
		name = string(os.PathSeparator) + name
	}

	return name, true
}

func (fi *diskFileInfo) IsDir() bool {
	return fi.finfo.IsDir()
}

// diskFile implements the File interface
type diskFile struct {
	root string
	name string
	mode os.FileMode
	file *os.File
}

func (f *diskFile) Name() string {
	return f.name
}

func (f *diskFile) Stat() (FileInfo, error) {
	finfo, err := f.file.Stat()
	if err != nil {
		return nil, err
	}

	dif := &diskFileInfo{
		root:  f.root,
		name:  f.name,
		mode:  f.mode,
		finfo: finfo,
	}

	return dif, nil
}

func (f *diskFile) Close() error {
	return f.file.Close()
}

func (f *diskFile) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

// diskFilesystem implements the Filesystem interface
type diskFilesystem struct {
	metadata map[string]string
	lock     sync.RWMutex

	root string
	cwd  string

	// Current size of the filesystem in bytes
	currentSize   int64
	lastSizeCheck time.Time

	// Logger from the config
	logger log.Logger
}

// NewDiskFilesystem returns a new filesystem that is backed by the disk filesystem.
// The root is / and the working directory is whatever is returned by os.Getwd(). The value
// of Root in the config will be ignored.
func NewDiskFilesystem(config DiskConfig) (Filesystem, error) {
	fs := &diskFilesystem{
		metadata: make(map[string]string),
		root:     "/",
		cwd:      "/",
		logger:   config.Logger,
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fs.cwd = cwd

	if len(fs.cwd) == 0 {
		fs.cwd = "/"
	}

	fs.cwd = filepath.Clean(fs.cwd)
	if !filepath.IsAbs(fs.cwd) {
		return nil, fmt.Errorf("the current working directory must be an absolute path")
	}

	if fs.logger == nil {
		fs.logger = log.New("")
	}

	return fs, nil
}

// NewRootedDiskFilesystem returns a filesystem that is backed by the disk filesystem. The
// root of the filesystem is defined by DiskConfig.Root. The working directory is "/". Root
// must be directory. If it doesn't exist, it will be created
func NewRootedDiskFilesystem(config RootedDiskConfig) (Filesystem, error) {
	fs := &diskFilesystem{
		metadata: make(map[string]string),
		root:     config.Root,
		cwd:      "/",
		logger:   config.Logger,
	}

	if len(fs.root) == 0 {
		fs.root = "/"
	}

	if root, err := filepath.Abs(fs.root); err != nil {
		return nil, err
	} else {
		fs.root = root
	}

	err := os.MkdirAll(fs.root, 0700)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(fs.root)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("root is not a directory")
	}

	if fs.logger == nil {
		fs.logger = log.New("")
	}

	return fs, nil
}

func (fs *diskFilesystem) Name() string {
	return "disk"
}

func (fs *diskFilesystem) Type() string {
	return "disk"
}

func (fs *diskFilesystem) Metadata(key string) string {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	return fs.metadata[key]
}

func (fs *diskFilesystem) SetMetadata(key, data string) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.metadata[key] = data
}

func (fs *diskFilesystem) Size() (int64, int64) {
	// This is to cache the size for some time in order not to
	// stress the underlying filesystem too much.
	if time.Since(fs.lastSizeCheck) >= 10*time.Second {
		var size int64 = 0

		fs.walk(fs.root, func(path string, info os.FileInfo) {
			if info.IsDir() {
				return
			}

			size += info.Size()
		})

		fs.currentSize = size

		fs.lastSizeCheck = time.Now()
	}

	return fs.currentSize, -1
}

func (fs *diskFilesystem) Purge(size int64) int64 {
	return 0
}

func (fs *diskFilesystem) Files() int64 {
	var nfiles int64 = 0

	fs.walk(fs.root, func(path string, info os.FileInfo) {
		if info.IsDir() {
			return
		}

		nfiles++
	})

	return nfiles
}

func (fs *diskFilesystem) Symlink(oldname, newname string) error {
	oldname = fs.cleanPath(oldname)
	newname = fs.cleanPath(newname)

	info, err := os.Lstat(oldname)
	if err != nil {
		return err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("%s can't link to another link (%s)", newname, oldname)
	}

	if info.IsDir() {
		return fmt.Errorf("can't symlink directories")
	}

	return os.Symlink(oldname, newname)
}

func (fs *diskFilesystem) Open(path string) File {
	path = fs.cleanPath(path)

	df := &diskFile{
		root: fs.root,
		name: strings.TrimPrefix(path, fs.root),
	}

	info, err := os.Lstat(path)
	if err != nil {
		return nil
	}

	df.mode = info.Mode()

	f, err := os.Open(path)
	if err != nil {
		return nil
	}

	df.file = f

	return df
}

func (fs *diskFilesystem) ReadFile(path string) ([]byte, error) {
	path = fs.cleanPath(path)

	return os.ReadFile(path)
}

func (fs *diskFilesystem) WriteFileReader(path string, r io.Reader) (int64, bool, error) {
	path = fs.cleanPath(path)

	replace := true

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return -1, false, fmt.Errorf("creating file failed: %w", err)
	}

	var f *os.File
	var err error

	f, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return -1, false, fmt.Errorf("creating file failed: %w", err)
		}

		replace = false
	}

	defer f.Close()

	size, err := f.ReadFrom(r)
	if err != nil {
		return -1, false, fmt.Errorf("reading data failed: %w", err)
	}

	fs.lastSizeCheck = time.Time{}

	return size, !replace, nil
}

func (fs *diskFilesystem) WriteFile(path string, data []byte) (int64, bool, error) {
	return fs.WriteFileReader(path, bytes.NewBuffer(data))
}

func (fs *diskFilesystem) WriteFileSafe(path string, data []byte) (int64, bool, error) {
	path = fs.cleanPath(path)
	dir, filename := filepath.Split(path)

	tmpfile, err := os.CreateTemp(dir, filename)
	if err != nil {
		return -1, false, err
	}

	defer os.Remove(tmpfile.Name())

	size, err := tmpfile.Write(data)
	if err != nil {
		return -1, false, err
	}

	if err := tmpfile.Close(); err != nil {
		return -1, false, err
	}

	replace := false
	if _, err := fs.Stat(path); err == nil {
		replace = true
	}

	if err := fs.rename(tmpfile.Name(), path); err != nil {
		return -1, false, err
	}

	fs.lastSizeCheck = time.Time{}

	return int64(size), !replace, nil
}

func (fs *diskFilesystem) Rename(src, dst string) error {
	src = fs.cleanPath(src)
	dst = fs.cleanPath(dst)

	return fs.rename(src, dst)
}

func (fs *diskFilesystem) rename(src, dst string) error {
	if src == dst {
		return nil
	}

	// First try to rename the file
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// If renaming the file fails, copy the data
	if err := fs.copy(src, dst); err != nil {
		os.Remove(dst)
		return fmt.Errorf("failed to copy files: %w", err)
	}

	if err := os.Remove(src); err != nil {
		os.Remove(dst)
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	return nil
}

func (fs *diskFilesystem) Copy(src, dst string) error {
	src = fs.cleanPath(src)
	dst = fs.cleanPath(dst)

	return fs.copy(src, dst)
}

func (fs *diskFilesystem) copy(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}

	destination, err := os.Create(dst)
	if err != nil {
		source.Close()
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		source.Close()
		os.Remove(dst)
		return fmt.Errorf("failed to copy data from source to destination: %w", err)
	}

	source.Close()

	fs.lastSizeCheck = time.Time{}

	return nil
}

func (fs *diskFilesystem) MkdirAll(path string, perm os.FileMode) error {
	path = fs.cleanPath(path)

	return os.MkdirAll(path, perm)
}

func (fs *diskFilesystem) Stat(path string) (FileInfo, error) {
	path = fs.cleanPath(path)

	dif := &diskFileInfo{
		root: fs.root,
		name: strings.TrimPrefix(path, fs.root),
	}

	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	dif.mode = info.Mode()

	if info.Mode()&os.ModeSymlink != 0 {
		info, err = os.Stat(path)
		if err != nil {
			return nil, err
		}
	}

	dif.finfo = info

	return dif, nil
}

func (fs *diskFilesystem) Remove(path string) int64 {
	path = fs.cleanPath(path)

	finfo, err := os.Stat(path)
	if err != nil {
		return -1
	}

	size := finfo.Size()

	if err := os.Remove(path); err != nil {
		return -1
	}

	fs.lastSizeCheck = time.Time{}

	return size
}

func (fs *diskFilesystem) RemoveAll() int64 {
	return 0
}

func (fs *diskFilesystem) List(path, pattern string) []FileInfo {
	path = fs.cleanPath(path)
	files := []FileInfo{}

	fs.walk(path, func(path string, info os.FileInfo) {
		if path == fs.root {
			return
		}

		name := strings.TrimPrefix(path, fs.root)
		if name[0] != os.PathSeparator {
			name = string(os.PathSeparator) + name
		}

		if info.IsDir() {
			return
		}

		if len(pattern) != 0 {
			if ok, _ := glob.Match(pattern, name, '/'); !ok {
				return
			}
		}

		files = append(files, &diskFileInfo{
			root:  fs.root,
			name:  name,
			finfo: info,
		})
	})

	return files
}

func (fs *diskFilesystem) LookPath(file string) (string, error) {
	if strings.Contains(file, "/") {
		file = fs.cleanPath(file)
		err := fs.findExecutable(file)
		if err == nil {
			return file, nil
		}
		return "", os.ErrNotExist
	}
	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		path := filepath.Join(dir, file)
		path = fs.cleanPath(path)
		if err := fs.findExecutable(path); err == nil {
			if !filepath.IsAbs(path) {
				return path, os.ErrNotExist
			}
			return path, nil
		}
	}
	return "", os.ErrNotExist
}

func (fs *diskFilesystem) findExecutable(file string) error {
	d, err := fs.Stat(file)
	if err != nil {
		return err
	}

	m := d.Mode()
	if m.IsDir() {
		return fmt.Errorf("is a directory")
	}

	if m&0111 != 0 {
		return nil
	}

	return os.ErrPermission
}

func (fs *diskFilesystem) walk(path string, walkfn func(path string, info os.FileInfo)) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			walkfn(path, info)
			return nil
		}

		mode := info.Mode()
		if !mode.IsRegular() && mode&os.ModeSymlink == 0 {
			return nil
		}

		walkfn(path, info)

		return nil
	})
}

func (fs *diskFilesystem) cleanPath(path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(fs.cwd, path)
	}

	return filepath.Join(fs.root, filepath.Clean(path))
}
