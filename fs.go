package main

import (
	"io/fs"
	"path/filepath"
	"slices"
	"time"
)

type File struct {
	Path    string
	LastMod time.Time
}

// walkDir walks the given path recursively and returns a list of contained paths.
//
// Ignores:
// - all directories
// - all files in directories with a leading .
func walkDir(root string, extensions []string) []File {
	r := make([]File, 0, 32)
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || path == root || path[0] == '.' {
			return nil
		}
		if len(extensions) > 0 && !slices.Contains(extensions, filepath.Ext(path)[1:]) {
			return nil
		}
		r = append(r, File{path, info.ModTime()})
		return nil
	})
	return r
}
