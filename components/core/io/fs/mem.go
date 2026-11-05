package fs

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/datarhei/core/v16/glob"
	"github.com/datarhei/core/v16/log"
)

// MemConfig is the config that is required for creating
// a new memory filesystem.
type MemConfig struct {
	Logger log.Logger // For logging, optional
}

type memFileInfo struct {
	name    string    // Full name of the file (including path)
	size    int64     // The size of the file in bytes
	dir     bool      // Whether this file represents a directory
	lastMod time.Time // The time of the last modification of the file
	linkTo  string    // Where the file links to, empty if it's not a link
}

func (f *memFileInfo) Name() string {
	return f.name
}

func (f *memFileInfo) Size() int64 {
	return f.size
}

func (f *memFileInfo) Mode() fs.FileMode {
	mode := fs.FileMode(fs.ModePerm)

	if f.dir {
		mode |= fs.ModeDir
	}

	if len(f.linkTo) != 0 {
		mode |= fs.ModeSymlink
	}

	return mode
}

func (f *memFileInfo) ModTime() time.Time {
	return f.lastMod
}

func (f *memFileInfo) IsLink() (string, bool) {
	return f.linkTo, len(f.linkTo) != 0
}

func (f *memFileInfo) IsDir() bool {
	return f.dir
}

type memFile struct {
	memFileInfo
	data *bytes.Buffer // Contents of the file
}

func (f *memFile) Name() string {
	return f.name
}

func (f *memFile) Stat() (FileInfo, error) {
	info := &memFileInfo{
		name:    f.name,
		size:    f.size,
		dir:     f.dir,
		lastMod: f.lastMod,
		linkTo:  f.linkTo,
	}

	return info, nil
}

func (f *memFile) Read(p []byte) (int, error) {
	if f.data == nil {
		return 0, io.EOF
	}

	return f.data.Read(p)
}

func (f *memFile) Close() error {
	if f.data == nil {
		return io.EOF
	}

	f.data = nil

	return nil
}

type memFilesystem struct {
	metadata map[string]string
	metaLock sync.RWMutex

	// Mapping of path to file
	files *memStorage

	// Current size of the filesystem in bytes and its mutes
	currentSize int64
	sizeLock    sync.RWMutex

	// Logger from the config
	logger log.Logger
}

// NewMemFilesystem creates a new filesystem in memory that implements
// the Filesystem interface.
func NewMemFilesystem(config MemConfig) (Filesystem, error) {
	fs := &memFilesystem{
		metadata: make(map[string]string),
		logger:   config.Logger,
	}

	if fs.logger == nil {
		fs.logger = log.New("")
	}

	fs.logger = fs.logger.WithField("type", "mem")

	fs.files = newMemStorage()

	fs.logger.Debug().Log("Created")

	return fs, nil
}

func NewMemFilesystemFromDir(dir string, config MemConfig) (Filesystem, error) {
	mem, err := NewMemFilesystem(config)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		mode := info.Mode()
		if !mode.IsRegular() {
			return nil
		}

		if mode&os.ModeSymlink != 0 {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}

		defer file.Close()

		_, _, err = mem.WriteFileReader(path, file)
		if err != nil {
			return fmt.Errorf("can't copy %s", path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mem, nil
}

func (fs *memFilesystem) Name() string {
	return "mem"
}

func (fs *memFilesystem) Type() string {
	return "mem"
}

func (fs *memFilesystem) Metadata(key string) string {
	fs.metaLock.RLock()
	defer fs.metaLock.RUnlock()

	return fs.metadata[key]
}

func (fs *memFilesystem) SetMetadata(key, data string) {
	fs.metaLock.Lock()
	defer fs.metaLock.Unlock()

	fs.metadata[key] = data
}

func (fs *memFilesystem) Size() (int64, int64) {
	fs.sizeLock.RLock()
	defer fs.sizeLock.RUnlock()

	return fs.currentSize, -1
}

func (fs *memFilesystem) Files() int64 {
	nfiles := int64(0)

	fs.files.Range(func(key string, f *memFile) bool {
		if f.dir {
			return true
		}

		nfiles++

		return true
	})

	return nfiles
}

func (fs *memFilesystem) Open(path string) File {
	path = fs.cleanPath(path)

	file, ok := fs.files.LoadAndCopy(path)

	if !ok {
		return nil
	}

	newFile := &memFile{
		memFileInfo: memFileInfo{
			name:    file.name,
			size:    file.size,
			lastMod: file.lastMod,
			linkTo:  file.linkTo,
		},
	}

	if len(file.linkTo) != 0 {
		file.Close()

		file, ok = fs.files.LoadAndCopy(file.linkTo)
		if !ok {
			return nil
		}
	}

	if file.data != nil {
		newFile.lastMod = file.lastMod
		newFile.data = bytes.NewBuffer(file.data.Bytes())
		newFile.size = int64(newFile.data.Len())
	}

	return newFile
}

func (fs *memFilesystem) ReadFile(path string) ([]byte, error) {
	path = fs.cleanPath(path)

	file, ok := fs.files.LoadAndCopy(path)

	if !ok {
		return nil, os.ErrNotExist
	}

	if len(file.linkTo) != 0 {
		file.Close()

		file, ok = fs.files.LoadAndCopy(file.linkTo)
		if !ok {
			return nil, os.ErrNotExist
		}
	}

	defer file.Close()

	if file.data != nil {
		return file.data.Bytes(), nil
	}

	return nil, nil
}

func (fs *memFilesystem) Symlink(oldname, newname string) error {
	oldname = fs.cleanPath(oldname)
	newname = fs.cleanPath(newname)

	if fs.files.Has(newname) {
		return os.ErrExist
	}

	oldFile, ok := fs.files.Load(oldname)
	if !ok {
		return os.ErrNotExist
	}

	if len(oldFile.linkTo) != 0 {
		return fmt.Errorf("%s can't link to another link (%s)", newname, oldname)
	}

	newFile := &memFile{
		memFileInfo: memFileInfo{
			name:    newname,
			dir:     false,
			size:    0,
			lastMod: time.Now(),
			linkTo:  oldname,
		},
		data: nil,
	}

	oldFile, loaded := fs.files.Store(newname, newFile)

	fs.sizeLock.Lock()
	defer fs.sizeLock.Unlock()

	if loaded {
		oldFile.Close()
		fs.currentSize -= oldFile.size
	}

	fs.currentSize += newFile.size

	return nil
}

func (fs *memFilesystem) WriteFileReader(path string, r io.Reader) (int64, bool, error) {
	path = fs.cleanPath(path)

	newFile := &memFile{
		memFileInfo: memFileInfo{
			name:    path,
			dir:     false,
			size:    0,
			lastMod: time.Now(),
		},
		data: &bytes.Buffer{},
	}

	size, err := newFile.data.ReadFrom(r)
	if err != nil {
		fs.logger.WithFields(log.Fields{
			"path":           path,
			"filesize_bytes": size,
			"error":          err,
		}).Warn().Log("Incomplete file")
	}

	newFile.size = size

	oldFile, replace := fs.files.Store(path, newFile)

	fs.sizeLock.Lock()
	defer fs.sizeLock.Unlock()

	if replace {
		oldFile.Close()

		fs.currentSize -= oldFile.size
	}

	fs.currentSize += newFile.size

	logger := fs.logger.WithFields(log.Fields{
		"path":           newFile.name,
		"filesize_bytes": newFile.size,
		"size_bytes":     fs.currentSize,
	})

	if replace {
		logger.Debug().Log("Replaced file")
	} else {
		logger.Debug().Log("Added file")
	}

	return newFile.size, !replace, nil
}

func (fs *memFilesystem) WriteFile(path string, data []byte) (int64, bool, error) {
	return fs.WriteFileReader(path, bytes.NewBuffer(data))
}

func (fs *memFilesystem) WriteFileSafe(path string, data []byte) (int64, bool, error) {
	return fs.WriteFileReader(path, bytes.NewBuffer(data))
}

func (fs *memFilesystem) Purge(size int64) int64 {
	files := []*memFile{}

	fs.files.Range(func(_ string, file *memFile) bool {
		files = append(files, file)
		return true
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i].lastMod.Before(files[j].lastMod)
	})

	var freed int64 = 0

	for _, f := range files {
		fs.files.Delete(f.name)
		size -= f.size
		freed += f.size

		fs.sizeLock.Lock()
		fs.currentSize -= f.size
		fs.sizeLock.Unlock()

		f.Close()

		fs.logger.WithFields(log.Fields{
			"path":           f.name,
			"filesize_bytes": f.size,
			"size_bytes":     fs.currentSize,
		}).Debug().Log("Purged file")

		if size <= 0 {
			break
		}
	}

	files = nil

	return freed
}

func (fs *memFilesystem) MkdirAll(path string, perm os.FileMode) error {
	path = fs.cleanPath(path)

	fs.sizeLock.Lock()
	defer fs.sizeLock.Unlock()

	info, err := fs.stat(path)
	if err == nil {
		if info.IsDir() {
			return nil
		}

		return os.ErrExist
	}

	f := &memFile{
		memFileInfo: memFileInfo{
			name:    path,
			size:    0,
			dir:     true,
			lastMod: time.Now(),
		},
		data: nil,
	}

	fs.files.Store(path, f)

	return nil
}

func (fs *memFilesystem) Rename(src, dst string) error {
	src = filepath.Join("/", filepath.Clean(src))
	dst = filepath.Join("/", filepath.Clean(dst))

	if src == dst {
		return nil
	}

	srcFile, ok := fs.files.Load(src)
	if !ok {
		return os.ErrNotExist
	}

	dstFile, replace := fs.files.Store(dst, srcFile)
	fs.files.Delete(src)

	fs.sizeLock.Lock()
	defer fs.sizeLock.Unlock()

	if replace {
		dstFile.Close()

		fs.currentSize -= dstFile.size
	}

	return nil
}

func (fs *memFilesystem) Copy(src, dst string) error {
	src = filepath.Join("/", filepath.Clean(src))
	dst = filepath.Join("/", filepath.Clean(dst))

	if src == dst {
		return nil
	}

	if fs.isDir(dst) {
		return os.ErrInvalid
	}

	srcFile, ok := fs.files.LoadAndCopy(src)
	if !ok {
		return os.ErrNotExist
	}

	if srcFile.dir {
		srcFile.Close()
		return os.ErrNotExist
	}

	dstFile := &memFile{
		memFileInfo: memFileInfo{
			name:    dst,
			dir:     false,
			size:    srcFile.size,
			lastMod: time.Now(),
		},
		data: srcFile.data,
	}

	f, replace := fs.files.Store(dst, dstFile)

	fs.sizeLock.Lock()
	defer fs.sizeLock.Unlock()

	if replace {
		f.Close()
		fs.currentSize -= f.size
	}

	fs.currentSize += dstFile.size

	return nil
}

func (fs *memFilesystem) Stat(path string) (FileInfo, error) {
	path = fs.cleanPath(path)

	return fs.stat(path)
}

func (fs *memFilesystem) stat(path string) (FileInfo, error) {
	file, ok := fs.files.Load(path)
	if ok {
		f := &memFileInfo{
			name:    file.name,
			size:    file.size,
			dir:     file.dir,
			lastMod: file.lastMod,
			linkTo:  file.linkTo,
		}

		if len(f.linkTo) != 0 {
			file, ok := fs.files.Load(f.linkTo)
			if !ok {
				return nil, os.ErrNotExist
			}

			f.lastMod = file.lastMod
			f.size = file.size
		}

		return f, nil
	}

	// Check for directories
	if !fs.isDir(path) {
		return nil, os.ErrNotExist
	}

	f := &memFileInfo{
		name:    path,
		size:    0,
		dir:     true,
		lastMod: time.Now(),
		linkTo:  "",
	}

	return f, nil
}

func (fs *memFilesystem) isDir(path string) bool {
	file, ok := fs.files.Load(path)
	if ok {
		return file.dir
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	if path == "/" {
		return true
	}

	found := false

	fs.files.Range(func(k string, _ *memFile) bool {
		if strings.HasPrefix(k, path) {
			found = true
			return false
		}

		return true
	})

	return found
}

func (fs *memFilesystem) Remove(path string) int64 {
	file, ok := fs.files.Delete(path)
	if ok {
		file.Close()

		fs.sizeLock.Lock()
		defer fs.sizeLock.Unlock()

		fs.currentSize -= file.size
	} else {
		return -1
	}

	fs.logger.WithFields(log.Fields{
		"path":           file.name,
		"filesize_bytes": file.size,
		"size_bytes":     fs.currentSize,
	}).Debug().Log("Removed file")

	return file.size
}

func (fs *memFilesystem) RemoveAll() int64 {
	fs.sizeLock.Lock()
	defer fs.sizeLock.Unlock()

	size := fs.currentSize

	fs.files = newMemStorage()
	fs.currentSize = 0

	return size
}

func (fs *memFilesystem) List(path, pattern string) []FileInfo {
	path = fs.cleanPath(path)
	files := []FileInfo{}

	var compiledPattern glob.Glob
	var err error

	if len(pattern) != 0 {
		compiledPattern, err = glob.Compile(pattern, '/')
		if err != nil {
			return nil
		}
	}

	fs.files.Range(func(key string, file *memFile) bool {
		if file.dir {
			return true
		}

		if !strings.HasPrefix(file.name, path) {
			return true
		}

		if compiledPattern != nil {
			if !compiledPattern.Match(file.name) {
				return true
			}
		}

		files = append(files, &memFileInfo{
			name:    file.name,
			size:    file.size,
			lastMod: file.lastMod,
			linkTo:  file.linkTo,
		})

		return true
	})

	return files
}

func (fs *memFilesystem) LookPath(file string) (string, error) {
	if strings.Contains(file, "/") {
		file = fs.cleanPath(file)
		info, err := fs.Stat(file)
		if err == nil {
			if !info.Mode().IsRegular() {
				return file, os.ErrNotExist
			}
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
		if info, err := fs.Stat(path); err == nil {
			if !filepath.IsAbs(path) {
				return path, os.ErrNotExist
			}
			if !info.Mode().IsRegular() {
				return path, os.ErrNotExist
			}
			return path, nil
		}
	}
	return "", os.ErrNotExist
}

func (fs *memFilesystem) cleanPath(path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join("/", path)
	}

	return filepath.Join("/", filepath.Clean(path))
}
