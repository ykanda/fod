package fod

import (
	"os"
	"path/filepath"
)

// Entry represent file|directory entry
type Entry struct {
	Path   string
	Type   string
	Marked bool
}

// alias to wditems
func entries(path string) []*Entry {

	// list directory entries
	readdir, err := os.ReadDir(path)
	if nil != err {
		os.Exit(1)
	}

	entries := []*Entry{}
	entries = append(entries, &Entry{
		Path: "../",
		Type: FsTypeDir,
	})

	for _, de := range readdir {
		name := de.Name()
		var abs string
		if _abs, err := filepath.Abs(filepath.Join(path, name)); err == nil {
			abs = _abs
		} else {
			continue
		}
		switch de.IsDir() {
		case true:
			entries = append(entries, &Entry{
				Path:   abs,
				Type:   FsTypeDir,
				Marked: false,
			})
		case false:
			entries = append(entries, &Entry{
				Path:   abs,
				Type:   FsTypeFile,
				Marked: false,
			})
		}
	}
	return entries
}

// get type character
// FsTypeDir  -> d
// FsTypeFile -> -
func (entry *Entry) typeCharcter() (tc string) {
	switch entry.Type {
	case FsTypeDir:
		tc = "d"
	case FsTypeFile:
		tc = "-"
	}
	return
}
