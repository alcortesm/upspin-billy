package upspinfs

import (
	"fmt"
	"os"

	billy "gopkg.in/src-d/go-billy.v2"
)

const separator = '/'

type Upspin struct{}

func New() *Upspin {
	return &Upspin{}
}

// Implements billy.Filesystem.
func (fs *Upspin) Create(filename string) (billy.File, error) {
	return nil, fmt.Errorf("TODO")
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
