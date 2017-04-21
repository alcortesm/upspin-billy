package upspinfs

import "upspin.io/upspin"

type File struct {
	upspin.File
	userName upspin.UserName
	isClosed bool
}

func newFile(file upspin.File, userName upspin.UserName) *File {
	return &File{
		File:     file,
		userName: userName,
		isClosed: false,
	}
}

func (f *File) Filename() string {
	// strip the userName and the separator from the pathName of the upspin file.
	return string(f.File.Name())[len(userName)+len(sep):]
}

func (f *File) IsClosed() bool {
	return f.isClosed
}

func (f *File) Write(p []byte) (n int, err error) {
	return f.File.Write(p)
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.File.Read(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	panic("File.Seek not implemented")
}

func (f *File) Close() (err error) {
	f.isClosed = true
	return f.File.Close()
}
