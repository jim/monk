package main

import (
	"fmt"
	"github.com/jim/monk"
)

func main() {
	cache := monk.NewAssetCache()

	r := &monk.Resolution{}
	r.Resolve("a.js", cache)

	for _, resolved := range r.Resolved {
		fmt.Printf("%s ", resolved)
	}
	fmt.Println()

	built := monk.Build(r, cache)
	fmt.Println(built)
}
