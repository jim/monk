package monk

import (
	"log"
	"regexp"
)

type Resolution struct {
	Resolved []string
	Seen     []string
}

// Resolve the asset at assetPath and its dependencies.
//
// TODO should return error
func (r *Resolution) Resolve(assetPath string, cache *AssetCache) {
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

