package upspinfs

import (
	"fmt"
	"os"
	"strings"

	"upspin.io/upspin"

	billy "gopkg.in/src-d/go-billy.v2"
	"gopkg.in/src-d/go-billy.v2/subdirfs"
)

const (
	// Separator used between userNames and file names in upspin.PathNames
	sep = "/"
	// upspin directories has no mode
	dirMode = 0
)

type Upspin struct {
	client   upspin.Client
	userName upspin.UserName
	append   bool
}

func New(c upspin.Client, userName upspin.UserName) *Upspin {
	return &Upspin{
		client:   c,
		userName: userName,
	}
}

// Implements billy.Filesystem.
func (fs *Upspin) Create(path string) (billy.File, error) {
	pathName := fs.pathName(path)

	alreadyADir, err := fs.isDir(pathName)
	if err != nil {
		return nil, err
	}
	if alreadyADir {
		return nil,
			fmt.Errorf("cannot create file: already a directory: %s", pathName)
	}

	fs.createSubDirs(path)

	f, err := fs.client.Create(pathName)
	if err != nil {
		return nil, err
	}

	return newFile(f, fs.userName), nil
}

func (fs *Upspin) createSubDirs(path string) error {
	dirs := strings.Split(path, sep)
	dirs = dirs[:len(dirs)-1] // strip the basename
	return fs.MkdirAll(fs.Join(dirs...), dirMode)
}

func (fs *Upspin) pathName(path string) upspin.PathName {
	if !strings.HasPrefix(path, sep) {
		path = sep + path
	}
	return upspin.PathName(string(fs.userName) + path)
}

func (fs *Upspin) isDir(path upspin.PathName) (bool, error) {
	dirEntry, found, err := fs.lookup(path)
	if err != nil {
		return false, err
	}
	if found {
		return dirEntry.IsDir(), nil
	}
	return false, nil
}

func (fs *Upspin) lookup(path upspin.PathName) (*upspin.DirEntry, bool, error) {
	followLinks := true
	dirEntry, err := fs.client.Lookup(path, followLinks)
	if err != nil {
		if strings.Contains(err.Error(), "item does not exist") {
			return nil, false, nil
		}
		return nil, false, err
	}
	return dirEntry, true, nil
}

// Implements billy.Filesystem.
func (fs *Upspin) Open(path string) (billy.File, error) {
	f, err := fs.client.Open(fs.pathName(path))
	if err != nil {
		return nil, err
	}
	return newFile(f, fs.userName), nil
}

// Implements billy.Filesystem.
func (fs *Upspin) OpenFile(filename string, flags int, _ os.FileMode) (billy.File, error) {
	if err := checkFlags(flags); err != nil {
		return nil, err
	}

	fn := fs.Open
	if flags&os.O_CREATE != 0 {
		fn = fs.Create
	}

	f, err := fn(filename)
	if err != nil {
		return nil, err
	}

	fs.append = flags&os.O_APPEND != 0

	return f, nil
}

const (
	rwMask = os.O_RDONLY | os.O_WRONLY | os.O_RDWR
)

func checkFlags(f int) error {
	switch f & rwMask {
	case os.O_RDONLY:
	case os.O_WRONLY:
	case os.O_RDWR:
	default:
		return fmt.Errorf(
			"invalid access mode: more than one O_RDONLY, O_WRONLY or O_RDWR detected")
	}
	return nil
}

// Implements billy.Filesystem.
func (fs *Upspin) Stat(filename string) (billy.FileInfo, error) {
	return nil, fmt.Errorf("Stat TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) ReadDir(base string) (entries []billy.FileInfo, err error) {
	return nil, fmt.Errorf("ReadDir TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) MkdirAll(path string, _ os.FileMode) error {
	dirs := strings.Split(path, sep)
	for _, d := range dirs {
		pathName := fs.pathName(d)
		dirEntry, found, err := fs.lookup(pathName)
		if err != nil {
			return err
		}
		if found {
			if dirEntry.IsDir() {
				continue
			}
			return fmt.Errorf("cannot create dir, it is already a file: %s",
				pathName)
		}
		_, err = fs.client.MakeDirectory(fs.pathName(path))
		if err != nil {
			return err
		}
	}
	return nil
}

// Implements billy.Filesystem.
func (fs *Upspin) TempFile(dir, prefix string) (billy.File, error) {
	return nil, fmt.Errorf("TempFile TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) Rename(from, to string) error {
	return fmt.Errorf("Rename TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) Remove(filename string) error {
	return fmt.Errorf("Remove TODO")
}

// Implements billy.Filesystem.
func (fs *Upspin) Join(elem ...string) string {
	return strings.Join(elem, sep)
}

// Implements billy.Filesystem.
func (fs *Upspin) Dir(path string) billy.Filesystem {
	return subdirfs.New(fs, path)
}

// Implements billy.Filesystem.
func (fs *Upspin) Base() string {
	return "Base TODO"
}
