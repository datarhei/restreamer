package fs

import (
	"bytes"
	"fmt"
	"io"
)

type SizedFilesystem interface {
	Filesystem

	// Resize resizes the filesystem to the new size. Files may need to be deleted.
	Resize(size int64, purge bool) error
}

type PurgeFilesystem interface {
	// Purge will free up at least size number of bytes and returns the actual
	// freed space in bytes.
	Purge(size int64) int64
}

type sizedFilesystem struct {
	Filesystem

	// Size is the capacity of the filesystem in bytes
	maxSize int64

	// Set true to automatically delete the oldest files until there's
	// enough space to store a new file
	purge bool
}

var _ PurgeFilesystem = &sizedFilesystem{}

func NewSizedFilesystem(fs Filesystem, maxSize int64, purge bool) (SizedFilesystem, error) {
	r := &sizedFilesystem{
		Filesystem: fs,
		maxSize:    maxSize,
		purge:      purge,
	}

	return r, nil
}

func (r *sizedFilesystem) Size() (int64, int64) {
	currentSize, _ := r.Filesystem.Size()

	return currentSize, r.maxSize
}

func (r *sizedFilesystem) Resize(size int64, purge bool) error {
	r.purge = purge

	currentSize, _ := r.Size()
	if size >= currentSize {
		// If the new size is the same or larger than the current size,
		// nothing to do.
		r.maxSize = size
		return nil
	}

	// If the new size is less than the current size, purge some files.
	r.Purge(currentSize - size)

	r.maxSize = size

	return nil
}

func (r *sizedFilesystem) WriteFileReader(path string, rd io.Reader) (int64, bool, error) {
	currentSize, maxSize := r.Size()
	if maxSize <= 0 {
		return r.Filesystem.WriteFileReader(path, rd)
	}

	data := bytes.Buffer{}
	size, err := data.ReadFrom(rd)
	if err != nil {
		return -1, false, err
	}

	// reject if the new file is larger than the available space
	if size > maxSize {
		return -1, false, fmt.Errorf("File is too big")
	}

	newSize := currentSize + size

	// Calculate the new size of the filesystem
	finfo, err := r.Filesystem.Stat(path)
	if err == nil {
		// If the file already exist, take it's size into account
		newSize -= finfo.Size()
	}

	// If the the new size is larger than the allowed size, we have to free
	// some space.
	if newSize > maxSize {
		if !r.purge {
			return -1, false, fmt.Errorf("not enough space on device")
		}

		if r.Purge(size) < size {
			return -1, false, fmt.Errorf("not enough space on device")
		}
	}

	return r.Filesystem.WriteFileReader(path, &data)
}

func (r *sizedFilesystem) WriteFile(path string, data []byte) (int64, bool, error) {
	return r.WriteFileReader(path, bytes.NewBuffer(data))
}

func (r *sizedFilesystem) WriteFileSafe(path string, data []byte) (int64, bool, error) {
	currentSize, maxSize := r.Size()
	if maxSize <= 0 {
		return r.Filesystem.WriteFile(path, data)
	}

	size := int64(len(data))

	// reject if the new file is larger than the available space
	if size > maxSize {
		return -1, false, fmt.Errorf("File is too big")
	}

	newSize := currentSize + size

	// Calculate the new size of the filesystem
	finfo, err := r.Filesystem.Stat(path)
	if err == nil {
		// If the file already exist, take it's size into account
		newSize -= finfo.Size()
	}

	// If the the new size is larger than the allowed size, we have to free
	// some space.
	if newSize > maxSize {
		if !r.purge {
			return -1, false, fmt.Errorf("not enough space on device")
		}

		if r.Purge(size) < size {
			return -1, false, fmt.Errorf("not enough space on device")
		}
	}

	return r.Filesystem.WriteFileSafe(path, data)
}

func (r *sizedFilesystem) Purge(size int64) int64 {
	if purger, ok := r.Filesystem.(PurgeFilesystem); ok {
		return purger.Purge(size)
	}

	return 0
	/*
	   files := r.Filesystem.List("/", "")

	   	sort.Slice(files, func(i, j int) bool {
	   		return files[i].ModTime().Before(files[j].ModTime())
	   	})

	   var freed int64 = 0

	   	for _, f := range files {
	   		r.Filesystem.Remove(f.Name())
	   		size -= f.Size()
	   		freed += f.Size()
	   		r.currentSize -= f.Size()

	   		if size <= 0 {
	   			break
	   		}
	   	}

	   files = nil

	   return freed
	*/
}
