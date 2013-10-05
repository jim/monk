package main

import "os"
import "fmt"
import "io/ioutil"
import "path"
import "log"
import "strings"
import "regexp"

/*os.FileInfo*/
type Asset struct {
	path       string
	extensions string
	edges      []string
}

type Resolution struct {
	resolved []string
	seen     []string
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

	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatal(err)
	}
	info, _ := os.Stat(absPath)
	stringified := string(bytes)
	newFile := &File{info, stringified}
	fc.Store[logicalPath] = newFile
	return stringified
}

func NewFileCache() *FileCache {
	return &FileCache{make(map[string]*File)}
}

// LogicalPath
// Path

func main() {

	cache := NewFileCache()

	r := &Resolution{}
	r.resolve("a.js", cache)

	for _, resolved := range r.resolved {
		fmt.Printf("%s ", resolved)
	}
	fmt.Println()

	built := build(r)
	fmt.Println(built)
}

func build(r *Resolution) string {
	contents := make([]string, len(r.resolved))
	for _, assetPath := range r.resolved {
		absPath := path.Join("assets", assetPath)
		bytes, _ := ioutil.ReadFile(absPath)
		header := fmt.Sprintf("/* %s */\n", absPath)
		contents = append(contents, header, string(bytes))
	}

	return strings.Join(contents, "\n")
}

func (r *Resolution) resolve(assetPath string, cache *FileCache) {
	fmt.Println(assetPath)
	r.seen = append(r.seen, assetPath)

	contents := cache.lookup(assetPath)
	e := edges(string(contents))

	for _, edge := range e {
		if !contains(edge, r.resolved) {
			if contains(edge, r.seen) {
				log.Fatal("circular dependency detected")
			}
			r.resolve(edge, cache)
		}
	}

	r.resolved = append(r.resolved, assetPath)
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
