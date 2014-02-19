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
	c.SearchPaths = append(c.SearchPaths, filepath.Clean(dirpath))
	return c, nil
}

// Return the contents of a file based on its logical path. Loads the content from
// disk if needed.
//
// logicalPath must have at least one extension.
func (c *Context) lookup(logicalPath string) (*Asset, error) {
	asset, ok := c.Store[logicalPath]
	if ok {
		return asset, nil
	}

	asset, err := c.findAssetInSearchPaths(logicalPath)
	// this should return just the asset path, which can then be passed to createAsset

	if err != nil {
		/*fmt.Printf("lookup %q failed: %q\n", logicalPath, err.Error())*/
		return nil, err
	}

	c.Store[logicalPath] = asset
	return asset, nil
}

// TODO this should return a Match object that includes absPath and logicalPath
func (c *Context) findPathInSearchPaths(logicalPath string) (string, os.FileInfo, error) {
	if len(c.SearchPaths) == 0 {
		return "", nil, fmt.Errorf("No search paths have been defined.")
	}

	// logicalPath must have at least one extension.
	if path.Ext(logicalPath) == "" {
		return "", nil, fmt.Errorf("Can not find '%s'. An extension is required to find an asset.", logicalPath)
	}

	for _, searchPath := range c.SearchPaths {
		absPath := path.Join(searchPath, logicalPath)

		// Look for exact match
		info, err := c.fs.Stat(absPath)
		if os.IsNotExist(err) {

			// Search the entire directory for a matching base name
			absPath, err := c.searchDirectory(searchPath, logicalPath)
			if err != nil {
				continue
			}

			return absPath, info, nil
		}

		// Found an exact match
		/*fmt.Printf("Found an exact match for %q\n", absPath)*/
		info, _ = c.fs.Stat(absPath)
		return absPath, info, nil
	}

	return "", nil, fmt.Errorf("Could not find a file matching %q in %v", logicalPath, c.SearchPaths)
}

func (c *Context) findAssetInSearchPaths(logicalPath string) (*Asset, error) {
	absPath, info, err := c.findPathInSearchPaths(logicalPath)
	if err != nil {
		return nil, err
	}
	return c.createAsset(absPath, info)
}

// Create and return a pointer to a new Asset. The content of the file at absPath will
// be used as the asset's contents.
//
// TODO passing both FileInfo and an absolute path here seems redundant.
func (c *Context) createAsset(absPath string, info os.FileInfo) (*Asset, error) {
	rawContent, err := c.loadAssetContent(absPath)
	if err != nil {
		/*fmt.Printf("failed to load asset content for %q\n", absPath)*/
		return nil, err
	}
	content, dependencies := extractDependencies(rawContent)

	for i, dep := range dependencies {
		if path.Ext(dep) == "" {
			ext := strings.Split(path.Base(absPath), ".")[1]
			dependencies[i] = fmt.Sprintf("%s.%s", dep, ext)
		}
	}

	return &Asset{info, content, dependencies}, nil
}

// Converts the wildcard/directory dependencies such as foo/* into an
// explicit list of dependencies based on what files are currently
// located at the provided path.
func explodeDependencies(absPath string, dependencies []string, fs fileSystem) []string {
	result := []string{}
	for _, req := range dependencies {
		if strings.HasSuffix(req, "/*") {
			dirName := req[0 : len(req)-2]

			infos, err := fs.ReadDir(dirName)

			// TODO return err instead of panicing.
			if err != nil {
				panic(err)
			}

			// get contents of directory, and add them (minus extensions) to result
			for _, info := range infos {
				if !info.IsDir() { // recursive requires not currently supported
					fullPath := path.Join(dirName, info.Name())
					result = append(result, fullPath)
				}
			}
		} else {
			result = append(result, req)
		}
	}
	return result
}

// Iterates over the immediate child nodes of dirPath, returning the absolute path
// to a matching file if one is found.
func (c *Context) searchDirectory(dirPath string, logicalPath string) (string, error) {
	files, err := c.fs.ReadDir(dirPath)
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
		filtered, err := ApplyFilter(c, content, ext)
		if err != nil {
			return "", err
		}
		content = filtered
	}

	return content, nil
}
