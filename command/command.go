package main

import (
	"fmt"
	"github.com/jim/monk"
	"path"
	"runtime"
)

func main() {
	context := monk.NewContext(monk.DiskFS{})

	r := &monk.Resolution{}

	_, filepath, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filepath), "../assets")

	context.SearchPath(dir)

	err := r.Resolve("e.js", context)

	if err != nil {
		panic(err)
	}

	for _, resolved := range r.Resolved {
		fmt.Printf("%s ", resolved)
	}
	fmt.Println()

	built := monk.Build(r, context)
	fmt.Println(built)
}
