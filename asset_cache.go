package monk

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type AssetCache struct {
	fs          fileSystem
	Store       map[string]*Asset
	SearchPaths []string
}

type Asset struct {
	os.FileInfo
	Content      string
	Dependencies []string
}

func NewAssetCache(fs fileSystem) *AssetCache {
	return &AssetCache{fs, make(map[string]*Asset), []string{}}
}

// Append a path to the list of asset paths to be searched for assets.
func (cache *AssetCache) SearchPath(dirpath string) (*AssetCache, error) {
	abs, err := filepath.Abs(dirpath)
	if err != nil {
		return cache, err
	}
	cache.SearchPaths = append(cache.SearchPaths, abs)
	return cache, nil
}

// Return the contents of a file based on its logical path. Loads the content from
// disk if needed.
func (fc *AssetCache) lookup(logicalPath string) (*Asset, error) {
	asset, ok := fc.Store[logicalPath]
	if ok {
		return asset, nil
	}

	asset, err := fc.findAssetInSearchPaths(logicalPath)

	if err != nil {
		/*fmt.Printf("lookup %q failed: %q\n", logicalPath, err.Error())*/
		return nil, err
	}

	fc.Store[logicalPath] = asset
	return asset, nil
}

func (ac *AssetCache) findAssetInSearchPaths(logicalPath string) (*Asset, error) {

	if len(ac.SearchPaths) == 0 {
		return nil, fmt.Errorf("No search paths have been defined.")
	}

	for _, searchPath := range ac.SearchPaths {
		absPath := path.Join(searchPath, logicalPath)

		// Look for exact match
		info, err := ac.fs.Stat(absPath)
		if os.IsNotExist(err) {

			// Search the entire directory for a matching base name
			absPath, err := ac.searchDirectory("assets", logicalPath)
			if err != nil {
				continue
			}

			return ac.createAsset(absPath, info)
		}

		// Found an exact match
		/*fmt.Printf("Found an exact match for %q\n", absPath)*/
		info, _ = ac.fs.Stat(absPath)
		return ac.createAsset(absPath, info)
	}

	return nil, fmt.Errorf("Could not find a file matching %q in %v", logicalPath, ac.SearchPaths)
}

// Create and return a pointer to a new Asset. The content of the file at absPath will
// be used as the asset's contents.
//
// TODO passing both FileInfo and an absolute path here seems redundant.
func (ac *AssetCache) createAsset(absPath string, info os.FileInfo) (*Asset, error) {
	rawContent, err := ac.loadAssetContent(absPath)
	if err != nil {
		fmt.Printf("failed to load asset content for %q\n", absPath)
		return nil, err
	}
	content, dependencies := extractDependencies(rawContent)
	return &Asset{info, content, dependencies}, nil
}

// Iterates over the immediate child nodes of dirPath, returning the absolute path
// to a matching file if one is found.
func (ac *AssetCache) searchDirectory(dirPath string, logicalPath string) (string, error) {
	files, err := ac.fs.ReadDir("assets")
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
func (ac *AssetCache) loadAssetContent(filePath string) (string, error) {
	bytes, err := ac.fs.ReadFile(filePath)
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
