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
	client      upspin.Client
	userName    upspin.UserName
	append      bool
	followLinks bool
}

func New(c upspin.Client, userName upspin.UserName) *Upspin {
	return &Upspin{
		client:      c,
		userName:    userName,
		append:      false,
		followLinks: true,
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
// O_SYNC is ignored.
func (fs *Upspin) OpenFile(filename string, flags int, _ os.FileMode) (billy.File, error) {
	if err := checkFlags(flags); err != nil {
		return nil, err
	}

	var ctor func(string) (billy.File, error)
	if isSet(flags, os.O_CREATE) {
		ctor = fs.Create
		if isSet(flags, os.O_EXCL) {
			found, err := fs.exists(fs.pathName(filename))
			if err != nil {
				return nil, err
			}
			if found {
				return nil, fmt.Errorf("file already exists: %s", filename)
			}

			// From open(2): When these two flags are specified [O_CREATE and
			// O_EXCL], symbolic links are not followed: if pathname is
			// a symbolic link, then open() fails regardless of where the
			// symbolic link points to.
			fs.followLinks = false
		}
	} else {
		ctor = fs.Open
		if canWrite(flags) && isSet(flags, os.O_TRUNC) {
			// delete the file contents
			fs.client.Put(fs.pathName(filename), []byte{})
		}
	}

	f, err := ctor(filename)
	if err != nil {
		return nil, err
	}

	fs.append = isSet(flags, os.O_APPEND)

	return f, nil
}

func isSet(flags int, bit int) bool {
	return flags&bit != 0
}

func canWrite(flags int) bool {
	return isSet(flags, os.O_WRONLY) || isSet(flags, os.O_RDWR)
}

func (fs *Upspin) exists(filename upspin.PathName) (bool, error) {
	_, found, err := fs.lookup(filename)
	return found, err
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
