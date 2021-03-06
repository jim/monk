package monk

import (
	"io"
	"io/ioutil"
	"os"
)

type fileSystem interface {
	Stat(name string) (os.FileInfo, error)
	ReadDir(name string) ([]os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
	Open(name string) (file, error)
}

type file interface {
	io.Closer
	io.Reader
	/*io.ReaderAt*/
	/*io.Seeker*/
	/*Stat() (os.FileInfo, error)*/
}

// DiskFS implements fileSystem using the local disk.
type DiskFS struct{}

func (DiskFS) Stat(name string) (os.FileInfo, error)      { return os.Stat(name) }
func (DiskFS) ReadFile(name string) ([]byte, error)       { return ioutil.ReadFile(name) }
func (DiskFS) ReadDir(name string) ([]os.FileInfo, error) { return ioutil.ReadDir(name) }
func (DiskFS) Open(name string) (file, error)             { return os.Open(name) }
