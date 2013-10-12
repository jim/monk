package monk

import (
	"fmt"
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

	asset, _ := cache.lookup(assetPath)

	for _, edge := range asset.Dependencies {
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
