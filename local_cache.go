package monk

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"
)

type LocalCache struct{}

type CachedFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (f CachedFileInfo) Name() string       { return f.name }
func (f CachedFileInfo) Size() int64        { return f.size }
func (f CachedFileInfo) Mode() os.FileMode  { return f.mode }
func (f CachedFileInfo) ModTime() time.Time { return f.modTime }
func (f CachedFileInfo) IsDir() bool        { return f.isDir }
func (f CachedFileInfo) Sys() interface{}   { return nil }

func (lc *LocalCache) Open(name string) (file http.File, err error) {
	content, err := Get(name)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	info := &CachedFileInfo{
		name:    name,
		size:    int64(len(content)),
		mode:    0777,
		modTime: time.Now(),
		isDir:   false,
	}
	file = CachedFile{info: info, buffer: bytes.NewBufferString(content)}

	return
}

type CachedFile struct {
	info   *CachedFileInfo
	buffer *bytes.Buffer
}

func (cf CachedFile) Close() error {
	return nil
}

func (cf CachedFile) Stat() (os.FileInfo, error) {
	return cf.info, nil
}

func (cf CachedFile) Readdir(count int) ([]os.FileInfo, error) {
	return []os.FileInfo{}, fmt.Errorf("not supported")
}

func (cf CachedFile) Read(p []byte) (int, error) {
	return cf.buffer.Read(p)
}

func (cf CachedFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
