package main

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/ryantate13/walker"
)

//go:embed example_fs
var exampleFS embed.FS

func filePaths(w walker.Walker) []string {
	paths := make([]string, 0)
	w.MustWalk("example_fs", func(path string, _ fs.FileInfo) bool {
		paths = append(paths, path)
		return true
	})
	return paths
}

func countFileTypes(w walker.Walker) (int, int) {
	numDirs := 0
	numFiles := 0
	w.MustWalk("example_fs", func(path string, e fs.FileInfo) bool {
		if e.IsDir() {
			numDirs++
		} else {
			numFiles++
		}
		return true
	})
	return numDirs, numFiles
}

func totalFileSize(w walker.Walker) int64 {
	size := int64(0)
	w.MustWalk("example_fs", func(path string, e fs.FileInfo) bool {
		if !e.IsDir() {
			size += e.Size()
		}
		return true
	})
	return size
}

func main() {
	w := walker.New(exampleFS)
	numDirs, numFiles := countFileTypes(w)
	fmt.Printf(
		"directories: %d\nfiles: %d\ntotal size: %d bytes\nfile paths:\n%s\n",
		numDirs,
		numFiles,
		totalFileSize(w),
		strings.Join(filePaths(w), "\n - "),
	)
}
