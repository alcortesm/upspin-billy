package upspinfs

import (
	"fmt"
	"os"

	"upspin.io/upspin"

	billy "gopkg.in/src-d/go-billy.v2"
)

const separator = '/'

type Upspin struct {
	client   upspin.Client
	userName upspin.UserName
}

func New(c upspin.Client, u upspin.UserName) *Upspin {
	return &Upspin{
		client:   c,
		userName: u,
	}
}

// Implements billy.Filesystem.
func (fs *Upspin) Create(path string) (billy.File, error) {
	return &File{}, nil
}

func (fs *Upspin) pathName(path string) upspin.PathName {
	return upspin.PathName(string(fs.userName) + path)
}

// Implements billy.Filesystem.
func (fs *Upspin) Open(filename string) (billy.File, error) {
	return nil, fmt.Errorf("TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) OpenFile(filename string, flag int, perm os.FileMode) (billy.File, error) {
	return nil, fmt.Errorf("TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) Stat(filename string) (billy.FileInfo, error) {
	return nil, fmt.Errorf("TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) ReadDir(base string) (entries []billy.FileInfo, err error) {
	return nil, fmt.Errorf("TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) MkdirAll(path string, perm os.FileMode) error {
	return fmt.Errorf("TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) TempFile(dir, prefix string) (billy.File, error) {
	return nil, fmt.Errorf("TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) Rename(from, to string) error {
	return fmt.Errorf("TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) Remove(filename string) error {
	return fmt.Errorf("TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) Join(elem ...string) string {
	return "TODO"
}

// Implements billy.Filesystem.
func (fs *Upspin) Dir(path string) billy.Filesystem {
	return nil
}

// Implements billy.Filesystem.
func (fs *Upspin) Base() string {
	return "TODO"
}
