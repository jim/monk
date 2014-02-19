package monk

import (
	"bytes"
	"os"
	"path"
	"strings"
	"time"
)

type TestFS struct {
	files map[string]*TestFSFile
}

type TestFSFile struct {
	info    os.FileInfo
	content []byte
}

type TestFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (f TestFileInfo) Name() string       { return f.name }
func (f TestFileInfo) Size() int64        { return f.size }
func (f TestFileInfo) Mode() os.FileMode  { return f.mode }
func (f TestFileInfo) ModTime() time.Time { return f.modTime }
func (f TestFileInfo) IsDir() bool        { return f.isDir }
func (f TestFileInfo) Sys() interface{}   { return nil }

func NewTestFS() *TestFS {
	return &TestFS{make(map[string]*TestFSFile)}
}

func (f TestFSFile) Close() error {
	return nil
}

func (f TestFSFile) Read(b []byte) (n int, err error) {
	reader := bytes.NewReader(f.content)
	return reader.Read(b)
}

func (fs *TestFS) File(name string, content string) {
	base := path.Base(name)
	info := TestFileInfo{
		name:    base,
		size:    int64(len(content)),
		mode:    0777,
		modTime: time.Now(),
		isDir:   false,
	}
	fs.files[name] = &TestFSFile{info, []byte(content)}
}

func (fs TestFS) Stat(name string) (os.FileInfo, error) {
	if file, ok := fs.files[name]; ok {
		return file.info, nil
	}
	return nil, &os.PathError{"stat", name, os.ErrNotExist}
}

func (fs TestFS) ReadFile(name string) ([]byte, error) {
	if file, ok := fs.files[name]; ok {
		return file.content, nil
	}
	return nil, &os.PathError{"stat", name, os.ErrNotExist}
}

func (fs TestFS) Open(name string) (file, error) {
	if file, ok := fs.files[name]; ok {
		return file, nil
	}
	return nil, &os.PathError{"stat", name, os.ErrNotExist}
}

func (fs TestFS) ReadDir(name string) ([]os.FileInfo, error) {
	files := []os.FileInfo{}
	for path, file := range fs.files {
		if strings.HasPrefix(path, name) {
			files = append(files, file.info)
		}
	}
	if len(files) == 0 {
		err := &os.PathError{"stat", name, os.ErrNotExist}
		return files, err
	}
	return files, nil
}
