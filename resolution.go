package monk

import (
	"fmt"
	"regexp"
)

type Resolution struct {
	Resolved []string
	Seen     []string
}

// Resolve the asset at assetPath and its dependencies.
//
// TODO should return error
func (r *Resolution) Resolve(assetPath string, cache *AssetCache) error {
	r.Seen = append(r.Seen, assetPath)

	contents, _ := cache.lookup(assetPath)
	e := findRequires(string(contents))

	for _, edge := range e {
		if !contains(edge, r.Resolved) {
			if contains(edge, r.Seen) {
				return fmt.Errorf("circular dependency detected: %s <-> %s", assetPath, edge)
			}
			r.Resolve(edge, cache)
		}
	}

	r.Resolved = append(r.Resolved, assetPath)
	return nil
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

func findRequires(fileContents string) []string {
	r, err := regexp.Compile(`//= require ['"]?([\w\.]+)["']?`)
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
