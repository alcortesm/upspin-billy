package upspinfs

import "upspin.io/upspin"

type File struct{}

func newFile(f upspin.File) *File {
	return &File{}
}

func (f *File) Filename() string {
	panic("not implemented")
}

func (f *File) IsClosed() bool {
	panic("not implemented")
}

func (f *File) Write(p []byte) (n int, err error) {
	panic("not implemented")
}

func (f *File) Read(p []byte) (n int, err error) {
	panic("not implemented")
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	panic("not implemented")
}

func (f *File) Close() error {
	panic("not implemented")
}
