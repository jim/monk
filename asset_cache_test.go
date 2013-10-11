package monk

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
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

func NewTestFS() *TestFS {
	return &TestFS{make(map[string]*TestFSFile)}
}

func (f TestFileInfo) Name() string       { return f.name }
func (f TestFileInfo) Size() int64        { return f.size }
func (f TestFileInfo) Mode() os.FileMode  { return f.mode }
func (f TestFileInfo) ModTime() time.Time { return f.modTime }
func (f TestFileInfo) IsDir() bool        { return f.isDir }
func (f TestFileInfo) Sys() interface{}   { return nil }

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
	return nil, &os.PathError{"stat", name, fmt.Errorf("could not stat file")}
}

func (fs TestFS) ReadFile(name string) ([]byte, error) {
	if file, ok := fs.files[name]; ok {
		return file.content, nil
	}
	return nil, &os.PathError{"stat", name, fmt.Errorf("%s: no such file or directory", name)}
}

func (fs TestFS) ReadDir(name string) ([]os.FileInfo, error) {
	files := []os.FileInfo{}
	for path, file := range fs.files {
		if strings.HasPrefix(path, name) {
			files = append(files, file.info)
		}
	}
	if len(files) == 0 {
		err := &os.PathError{"stat", name, fmt.Errorf("%s: no such file or directory", name)}
		return files, err
	}
	return files, nil
}

func TestLoadAsset(t *testing.T) {
	fs := NewTestFS()
	ac := NewAssetCache(fs)
	assetContent := "//= require a"
	assetPath := "assets/simple.js"

	fs.File(assetPath, assetContent)
	ac.SearchPath("assets")

	if content, err := ac.loadAsset(assetPath); err == nil {
		if content != assetContent {
			t.Errorf("requiring %q, want %q, got %q", assetPath, assetContent, content)
		}
	} else {
		t.Error(err)
	}
}
