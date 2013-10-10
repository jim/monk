package monk

import (
	"fmt"
	"log"
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
	Content string
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
func (fc *AssetCache) lookup(logicalPath string) (string, error) {
	file, ok := fc.Store[logicalPath]
	if ok {
		return file.Content, nil
	}

	asset, err := fc.findAssetInSearchPaths(logicalPath)

	if err != nil {
		return "", err
	}

	fc.Store[logicalPath] = asset
	return asset.Content, nil
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
			asset, err := ac.searchDirectory("assets", logicalPath)
			if err != nil {
				continue
			}
			return asset, nil
		}

		// handle exact match
		content, err := ac.loadAsset(absPath)
		if err != nil {
			log.Fatal(err)
		}

		return &Asset{info, content}, nil
	}

	return nil, fmt.Errorf("Could not find a file matching %q in %v", logicalPath, ac.SearchPaths)
}

func (ac *AssetCache) searchDirectory(dirPath string, logicalPath string) (*Asset, error) {
	files, err := ac.fs.ReadDir("assets")
	if os.IsNotExist(err) {
		return nil, err
	}

	pattern := fmt.Sprintf(`^%s[\.\w+]+`, regexp.QuoteMeta(logicalPath))
	r, _ := regexp.Compile(pattern)
	for _, fileInfo := range files {
		name := fileInfo.Name()
		if r.MatchString(name) {

			absPath := path.Join(dirPath, name)
			content, err := ac.loadAsset(absPath)
			if err != nil {
				log.Fatal(err)
			}

			return &Asset{fileInfo, content}, nil
		}
	}
	return nil, fmt.Errorf("Could not find a file matching %s/%s", dirPath, logicalPath)
}

// Loads a file from filePath, filtering its contents through a series filters based
// on the additional extensions in the filename. The first extension is assumed to
// be the final type of the file.
func (ac *AssetCache) loadAsset(filePath string) (string, error) {
	bytes, _ := ac.fs.ReadFile(filePath)
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
