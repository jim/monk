package monk

import (
	"fmt"
	"strings"
)

func Build(r *Resolution, cache *AssetCache) string {
	contents := make([]string, len(r.Resolved))
	for _, logicalPath := range r.Resolved {
		asset, err := cache.lookup(logicalPath)
		if err != nil {
			panic(err)
		}
		header := fmt.Sprintf("/* %s */\n", logicalPath)
		contents = append(contents, header, asset.Content)
	}

	return strings.Join(contents, "\n")
}
