package monk

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
  "bytes"
  "os/exec"
)

type AssetCache struct {
	Store map[string]*Asset
}

type Asset struct {
	os.FileInfo
	Content string
}

func NewAssetCache() *AssetCache {
	return &AssetCache{make(map[string]*Asset)}
}

// Return the contents of a file based on its logical path. Loads the content from
// disk if needed.
func (fc *AssetCache) lookup(logicalPath string) string {
	file, ok := fc.Store[logicalPath]
	if ok {
		return file.Content
	}
	absPath := path.Join("assets", logicalPath)

	info, err := os.Stat(absPath)

	if os.IsNotExist(err) {
		p, err := searchDirectory("assets", logicalPath)
		if err != nil {
			log.Fatal(err)
		}
		absPath = p
	}

	content, _ := loadAsset(absPath)

	newAsset := &Asset{info, content}
	fc.Store[logicalPath] = newAsset
	return content
}

// Find if a matching path is in a directory. If so, returns the full path to the
// file on disk.
func searchDirectory(dirPath string, logicalPath string) (string, error) {
	files, err := ioutil.ReadDir("assets")
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
func loadAsset(filePath string) (string, error) {
	bytes, _ := ioutil.ReadFile(filePath)
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

func filter(content string, extension string) (string, error) {
	switch extension {
	case "bs":
		return strings.Replace(content, "a", "b", -1), nil
	case "fs":
		return strings.Replace(content, "f", "x", -1), nil
  case "coffee":
    return coffeeFilter(content)
	}
	return content, nil
}

func coffeeFilter(content string) (string, error) {
	cmd := exec.Command("coffee", "-s", "-c")
	cmd.Stdin = strings.NewReader(content)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
  return out.String(), err
}
