package monk

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

/*os.FileInfo*/
type Asset struct {
	path       string
	extensions string
	edges      []string
}

type Resolution struct {
	Resolved []string
	Seen     []string
}

type File struct {
	os.FileInfo
	Content string
}

type FileCache struct {
	Store map[string]*File
}

// Return the contents of a file based on its logical path. Loads the content from
// disk if needed.
func (fc *FileCache) lookup(logicalPath string) string {
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

	content, _ := loadFile(absPath)

	newFile := &File{info, content}
	fc.Store[logicalPath] = newFile
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

func loadFile(filePath string) (string, error) {
	bytes, _ := ioutil.ReadFile(filePath)
	content := string(bytes)

	exts := strings.Split(path.Base(filePath), ".")

	// Nothing elce to do if there aren't additional extensions
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
	}
	return content, nil
}

func NewFileCache() *FileCache {
	return &FileCache{make(map[string]*File)}
}

// LogicalPath
// Path

func Build(r *Resolution, cache *FileCache) string {
	contents := make([]string, len(r.Resolved))
	for _, logicalPath := range r.Resolved {
		content := cache.lookup(logicalPath)
		header := fmt.Sprintf("/* %s */\n", logicalPath)
		contents = append(contents, header, content)
	}

	return strings.Join(contents, "\n")
}

func (r *Resolution) Resolve(assetPath string, cache *FileCache) {
	r.Seen = append(r.Seen, assetPath)

	contents := cache.lookup(assetPath)
	e := edges(string(contents))

	for _, edge := range e {
		if !contains(edge, r.Resolved) {
			if contains(edge, r.Seen) {
				log.Fatal("circular dependency detected")
			}
			r.Resolve(edge, cache)
		}
	}

	r.Resolved = append(r.Resolved, assetPath)
}

func contains(needle string, haystack []string) bool {
	found := false

	for _, item := range haystack {
		if needle == item {
			found = true
			break
		}
	}

	return found
}

func edges(fileContents string) []string {
	r, err := regexp.Compile(`//= require ([\w\.]+)`)
	if err != nil {
		panic(err)
	}

	matches := r.FindAllStringSubmatch(fileContents, -1)
	requires := make([]string, 0)

	for _, m := range matches {
		requires = append(requires, m[1])
	}

	return requires
}
