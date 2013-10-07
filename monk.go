package monk

import (
	"fmt"
	"os"
	"strings"
)

/*os.FileInfo*/
type Asset struct {
	path       string
	extensions string
	edges      []string
}
// LogicalPath
// Path

type Resolution struct {
	Resolved []string
	Seen     []string
}

type File struct {
	os.FileInfo
	Content string
}

func Build(r *Resolution, cache *FileCache) string {
	contents := make([]string, len(r.Resolved))
	for _, logicalPath := range r.Resolved {
		content := cache.lookup(logicalPath)
		header := fmt.Sprintf("/* %s */\n", logicalPath)
		contents = append(contents, header, content)
	}

	return strings.Join(contents, "\n")
}

