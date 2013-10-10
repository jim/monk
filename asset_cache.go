package monk

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

type AssetCache struct {
  fs fileSystem
	Store map[string]*Asset
}

type Asset struct {
	os.FileInfo
	Content string
}

func NewAssetCache(fs fileSystem) *AssetCache {
	return &AssetCache{fs, make(map[string]*Asset)}
}

// Return the contents of a file based on its logical path. Loads the content from
// disk if needed.
func (fc *AssetCache) lookup(logicalPath string) string {
	file, ok := fc.Store[logicalPath]
	if ok {
		return file.Content
	}
	absPath := path.Join("assets", logicalPath)

	info, err := fc.fs.Stat(absPath)

	if os.IsNotExist(err) {
		p, err := fc.searchDirectory("assets", logicalPath)
		if err != nil {
			log.Fatal(err)
		}
		absPath = p
	}

  content, err := fc.loadAsset(absPath)
  if err != nil {
    log.Fatal(err)
  }

	newAsset := &Asset{info, content}
	fc.Store[logicalPath] = newAsset
	return content
}

// Find if a matching path is in a directory. If so, returns the full path to the
// file on disk.
func (fc *AssetCache) searchDirectory(dirPath string, logicalPath string) (string, error) {
	files, err := fc.fs.ReadDir("assets")
	if os.IsNotExist(err) {
		return "", err
	}

	pattern := fmt.Sprintf(`^%s[\.\w+]+`, regexp.QuoteMeta(logicalPath))
	r, _ := regexp.Compile(pattern)
	for _, fileInfo := range files {
		name := fileInfo.Name()
		if r.MatchString(name) {
			return path.Join(dirPath, name), nil
		}
	}
	return "", fmt.Errorf("Could not find a file matching %s/%s", dirPath, logicalPath)
}

// Loads a file from filePath, filtering its contents through a series filters based
// on the additional extensions in the filename. The first extension is assumed to
// be the final type of the file.
func (fc *AssetCache) loadAsset(filePath string) (string, error) {
	bytes, _ := fc.fs.ReadFile(filePath)
	content := string(bytes)

	exts := strings.Split(path.Base(filePath), ".")

	// Nothing else to do if there aren't additional extensions
	if len(exts) < 3 {
		return content, nil
	}

	exts = exts[2:]

	// Reverse the order of remaining extensions
	for i, j := 0, len(exts)-1; i < j; i, j = i+1, j-1 {
		exts[i], exts[j] = exts[j], exts[i]
	}

	for _, ext := range exts {
		filtered, err := filter(content, ext)
		if err != nil {
			return "", err
		}
		content = filtered
	}

	return content, nil
}

