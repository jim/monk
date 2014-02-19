package main

import (
	"flag"
	"fmt"
	"github.com/jim/monk"
	"os"
	"path"
	"strings"
)

type searchPaths []string

var searchPathsFlag searchPaths

func (sp *searchPaths) Set(searchPath string) error {
	if !strings.HasPrefix(searchPath, "/") {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		searchPath = path.Join(dir, searchPath)
	}
	*sp = append(*sp, searchPath)
	return nil
}

func (sp *searchPaths) String() string {
	return strings.Join(*sp, ",")
}

func init() {
	flag.Var(&searchPathsFlag, "s", "path to search for assets when building")
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		printUsage()
		return
	}

	r := &monk.Resolution{}
	context := monk.NewContext(monk.DiskFS{})

	if len(searchPathsFlag) == 0 {
		panic("You must specify at least one path using -s")
	}

	for _, sPath := range searchPathsFlag {
		context.SearchPath(sPath)
	}

	if flag.NArg() == 0 {
		panic("You must specify the asset to build.")
	}
	asset := flag.Arg(0)

	err := r.Resolve(asset, context)

	if err != nil {
		panic(err)
	}

	built := monk.Build(r, context)
	fmt.Println(built)
}

func printUsage() {
	fmt.Println("monk, a tool to build assets")
	fmt.Println("  usage: monk [OPTIONS] asset_to_build.ext\n")
	flag.PrintDefaults()
}
