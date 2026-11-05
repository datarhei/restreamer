package value

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/datarhei/core/v16/io/fs"
)

// must directory

type MustDir struct {
	p  *string
	fs fs.Filesystem
}

func NewMustDir(p *string, val string, fs fs.Filesystem) *MustDir {
	v := &MustDir{
		p:  p,
		fs: fs,
	}

	*p = val

	return v
}

func (u *MustDir) Set(val string) error {
	*u.p = val
	return nil
}

func (u *MustDir) String() string {
	return *u.p
}

func (u *MustDir) Validate() error {
	val := *u.p

	if len(strings.TrimSpace(val)) == 0 {
		return fmt.Errorf("path name must not be empty")
	}

	if err := u.fs.MkdirAll(val, 0750); err != nil {
		return fmt.Errorf("%s can't be created (%w)", val, err)
	}

	finfo, err := u.fs.Stat(val)
	if err != nil {
		return fmt.Errorf("%s does not exist", val)
	}

	if !finfo.IsDir() {
		return fmt.Errorf("%s is not a directory", val)
	}

	return nil
}

func (u *MustDir) IsEmpty() bool {
	return len(*u.p) == 0
}

// directory

type Dir struct {
	p  *string
	fs fs.Filesystem
}

func NewDir(p *string, val string, fs fs.Filesystem) *Dir {
	v := &Dir{
		p:  p,
		fs: fs,
	}

	*p = val

	return v
}

func (u *Dir) Set(val string) error {
	*u.p = val
	return nil
}

func (u *Dir) String() string {
	return *u.p
}

func (u *Dir) Validate() error {
	val := *u.p

	if len(strings.TrimSpace(val)) == 0 {
		return nil
	}

	finfo, err := u.fs.Stat(val)
	if err != nil {
		return fmt.Errorf("%s does not exist", val)
	}

	if !finfo.IsDir() {
		return fmt.Errorf("%s is not a directory", val)
	}

	return nil
}

func (u *Dir) IsEmpty() bool {
	return len(*u.p) == 0
}

// executable

type Exec struct {
	p  *string
	fs fs.Filesystem
}

func NewExec(p *string, val string, fs fs.Filesystem) *Exec {
	v := &Exec{
		p:  p,
		fs: fs,
	}

	*p = val

	return v
}

func (u *Exec) Set(val string) error {
	*u.p = val
	return nil
}

func (u *Exec) String() string {
	return *u.p
}

func (u *Exec) Validate() error {
	val := *u.p

	_, err := u.fs.LookPath(val)
	if err != nil {
		return fmt.Errorf("%s not found or is not executable", val)
	}

	return nil
}

func (u *Exec) IsEmpty() bool {
	return len(*u.p) == 0
}

// regular file

type File struct {
	p  *string
	fs fs.Filesystem
}

func NewFile(p *string, val string, fs fs.Filesystem) *File {
	v := &File{
		p:  p,
		fs: fs,
	}

	*p = val

	return v
}

func (u *File) Set(val string) error {
	*u.p = val
	return nil
}

func (u *File) String() string {
	return *u.p
}

func (u *File) Validate() error {
	val := *u.p

	if len(val) == 0 {
		return nil
	}

	finfo, err := u.fs.Stat(val)
	if err != nil {
		return fmt.Errorf("%s does not exist", val)
	}

	if !finfo.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", val)
	}

	return nil
}

func (u *File) IsEmpty() bool {
	return len(*u.p) == 0
}

// absolute path

type AbsolutePath string

func NewAbsolutePath(p *string, val string) *AbsolutePath {
	*p = filepath.Clean(val)

	return (*AbsolutePath)(p)
}

func (s *AbsolutePath) Set(val string) error {
	*s = AbsolutePath(filepath.Clean(val))
	return nil
}

func (s *AbsolutePath) String() string {
	return string(*s)
}

func (s *AbsolutePath) Validate() error {
	path := string(*s)

	if !filepath.IsAbs(path) {
		return fmt.Errorf("%s is not an absolute path", path)
	}

	return nil
}

func (s *AbsolutePath) IsEmpty() bool {
	return len(string(*s)) == 0
}
