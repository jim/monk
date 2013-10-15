package monk

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

// Get the asset specified by assetPath.
func Get(assetPath string) (string, error) {
	cache := NewContext(DiskFS{})

	r := &Resolution{}

	_, filepath, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filepath), "assets")

	cache.SearchPath(dir)

	if err := r.Resolve(assetPath, cache); err != nil {
		return "", err
	}

	return Build(r, cache), nil
}

func Build(r *Resolution, context *Context) string {
	contents := make([]string, len(r.Resolved))
	for _, logicalPath := range r.Resolved {
		asset, err := context.lookup(logicalPath)
		if err != nil {
			panic(err)
		}
		header := fmt.Sprintf("/* %s */\n", logicalPath)
		contents = append(contents, header, asset.Content)
	}

	return strings.Join(contents, "\n")
}
