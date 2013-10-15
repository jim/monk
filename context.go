package monk

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Context struct {
	fs          fileSystem
	Store       map[string]*Asset
	SearchPaths []string
}

type Asset struct {
	os.FileInfo
	Content      string
	Dependencies []string
}

func NewContext(fs fileSystem) *Context {
	return &Context{fs, make(map[string]*Asset), []string{}}
}

// Append a path to the list of asset paths to be searched for assets.
func (c *Context) SearchPath(dirpath string) (*Context, error) {
	abs, err := filepath.Abs(dirpath)
	if err != nil {
		return c, err
	}
	c.SearchPaths = append(c.SearchPaths, abs)
	return c, nil
}

// Return the contents of a file based on its logical path. Loads the content from
// disk if needed.
func (c *Context) lookup(logicalPath string) (*Asset, error) {
	asset, ok := c.Store[logicalPath]
	if ok {
		return asset, nil
	}

	asset, err := c.findAssetInSearchPaths(logicalPath)

	if err != nil {
		/*fmt.Printf("lookup %q failed: %q\n", logicalPath, err.Error())*/
		return nil, err
	}

	c.Store[logicalPath] = asset
	return asset, nil
}

func (c *Context) findAssetInSearchPaths(logicalPath string) (*Asset, error) {

	if len(c.SearchPaths) == 0 {
		return nil, fmt.Errorf("No search paths have been defined.")
	}

	for _, searchPath := range c.SearchPaths {
		absPath := path.Join(searchPath, logicalPath)

		// Look for exact match
		info, err := c.fs.Stat(absPath)
		if os.IsNotExist(err) {

			// Search the entire directory for a matching base name
			absPath, err := c.searchDirectory("assets", logicalPath)
			if err != nil {
				continue
			}

			return c.createAsset(absPath, info)
		}

		// Found an exact match
		/*fmt.Printf("Found an exact match for %q\n", absPath)*/
		info, _ = c.fs.Stat(absPath)
		return c.createAsset(absPath, info)
	}

	return nil, fmt.Errorf("Could not find a file matching %q in %v", logicalPath, c.SearchPaths)
}

// Create and return a pointer to a new Asset. The content of the file at absPath will
// be used as the asset's contents.
//
// TODO passing both FileInfo and an absolute path here seems redundant.
func (c *Context) createAsset(absPath string, info os.FileInfo) (*Asset, error) {
	rawContent, err := c.loadAssetContent(absPath)
	if err != nil {
		fmt.Printf("failed to load asset content for %q\n", absPath)
		return nil, err
	}
	content, dependencies := extractDependencies(rawContent)
	return &Asset{info, content, dependencies}, nil
}

// Iterates over the immediate child nodes of dirPath, returning the absolute path
// to a matching file if one is found.
func (c *Context) searchDirectory(dirPath string, logicalPath string) (string, error) {
	files, err := c.fs.ReadDir("assets")
	if os.IsNotExist(err) {
		return "", err
	}

	strippedPath := strings.TrimPrefix(logicalPath, "/")
	pattern := fmt.Sprintf(`^%s[\.\w+]+`, regexp.QuoteMeta(strippedPath))
	r, _ := regexp.Compile(pattern)
	for _, fileInfo := range files {
		name := fileInfo.Name()
		if r.MatchString(name) {
			absPath := path.Join(dirPath, name)
			return absPath, nil
		}
	}
	return "", fmt.Errorf("Could not find a file matching %s/%s", dirPath, logicalPath)
}

// Loads a file from filePath, filtering its contents through a series filters based
// on the additional extensions in the filename. The first extension is assumed to
// be the final type of the file.
func (c *Context) loadAssetContent(filePath string) (string, error) {
	bytes, err := c.fs.ReadFile(filePath)
	if err != nil {
		return "", err
	}

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
		filtered, err := ApplyFilter(content, ext)
		if err != nil {
			return "", err
		}
		content = filtered
	}

	return content, nil
}
