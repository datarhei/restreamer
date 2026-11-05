package fs

import (
	"io"
	"os"
)

type readOnlyFilesystem struct {
	Filesystem
}

func NewReadOnlyFilesystem(fs Filesystem) (Filesystem, error) {
	r := &readOnlyFilesystem{
		Filesystem: fs,
	}

	return r, nil
}

func (r *readOnlyFilesystem) Symlink(oldname, newname string) error {
	return os.ErrPermission
}

func (r *readOnlyFilesystem) WriteFileReader(path string, rd io.Reader) (int64, bool, error) {
	return -1, false, os.ErrPermission
}

func (r *readOnlyFilesystem) WriteFile(path string, data []byte) (int64, bool, error) {
	return -1, false, os.ErrPermission
}

func (r *readOnlyFilesystem) WriteFileSafe(path string, data []byte) (int64, bool, error) {
	return -1, false, os.ErrPermission
}

func (r *readOnlyFilesystem) MkdirAll(path string, perm os.FileMode) error {
	return os.ErrPermission
}

func (r *readOnlyFilesystem) Remove(path string) int64 {
	return -1
}

func (r *readOnlyFilesystem) RemoveAll() int64 {
	return 0
}

func (r *readOnlyFilesystem) Purge(size int64) int64 {
	return 0
}

func (r *readOnlyFilesystem) Resize(size int64, purge bool) error {
	return os.ErrPermission
}
