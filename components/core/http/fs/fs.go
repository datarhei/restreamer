package fs

import (
	"github.com/datarhei/core/v16/http/cache"
	"github.com/datarhei/core/v16/io/fs"
)

type FS struct {
	Name       string
	Mountpoint string

	AllowWrite bool

	EnableAuth bool
	Username   string
	Password   string

	DefaultFile        string
	DefaultContentType string
	Gzip               bool

	Filesystem fs.Filesystem

	Cache cache.Cacher
}
