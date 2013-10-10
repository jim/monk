package main

import (
	"fmt"
	"github.com/jim/monk"
	"path"
	"runtime"
)

func main() {
	cache := monk.NewAssetCache(monk.DiskFS{})

	r := &monk.Resolution{}

	_, filepath, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filepath), "../assets")

	cache.SearchPath(dir)

	r.Resolve("a.js", cache)

	for _, resolved := range r.Resolved {
		fmt.Printf("%s ", resolved)
	}
	fmt.Println()

	built := monk.Build(r, cache)
	fmt.Println(built)
}
