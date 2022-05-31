package walker

import (
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/ryantate13/walker/stack"
)

// WalkFunc is the type of callback passed to Walk.
//
// It will receive three arguments, the full path as a string, the file info and any error encountered
// In the case of directories the callback may be invoked twice, once to stat the directory and again
// when attempting to read it. Returning false on either attempt will prevent further iteration
type WalkFunc func(path string, entry fs.FileInfo, err error) bool

// MustWalkFunc is identical to WalkFunc except passed to MustWalk and any error will induce a panic
type MustWalkFunc func(path string, entry fs.FileInfo) bool

// A Walker walks an fs.FS processing directories before files and invoking a callback for each entry encountered.
// Each path will be prefixed by the dir argument passed in, and each entry will is guarantee to be below it.
// Note that zero escape analysis is performed with regard to symlinks, etc.
type Walker interface {
	// Walk walks dir allowing user-defined error handling
	Walk(dir string, fn WalkFunc)
	// MustWalk walks dir and panics if errors are encountered
	MustWalk(dir string, fn MustWalkFunc)
}

type walker struct {
	fs fs.FS
}

// New instantiates a Walker
func New(fs fs.FS) Walker {
	return &walker{fs}
}

func (w *walker) Walk(dir string, fn WalkFunc) {
	s := stack.New[string]()
	s.Push(dir)
	for !s.Empty() {
		path := s.Pop()
		entry, err := fs.Stat(w.fs, path)
		if !fn(path, entry, err) {
			return
		}
		if entry != nil && entry.IsDir() {
			entries, err := fs.ReadDir(w.fs, path)
			if err != nil {
				if fn(path, entry, err) {
					continue
				} else {
					return
				}
			}
			sort.Slice(entries, func(i, j int) bool {
				return entries[i].IsDir() && !entries[j].IsDir()
			})
			for i := len(entries) - 1; i >= 0; i-- {
				p := filepath.Join(path, entries[i].Name())
				s.Push(p)
			}
		}
	}
}

func (w *walker) MustWalk(dir string, fn MustWalkFunc) {
	w.Walk(dir, func(path string, entry fs.FileInfo, err error) bool {
		if err != nil {
			panic(err)
		}
		return fn(path, entry)
	})
}
