package fod

import (
	"io/ioutil"
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
	readdir, err := ioutil.ReadDir(path)
	if nil != err {
		os.Exit(1)
	}

	entries := []*Entry{}
	entries = append(entries, &Entry{
		Path: "../",
		Type: FsTypeDir,
	})

	for _, fi := range readdir {
		var abs string // [todo] - FileInfo を Entry に含める
		if _abs, err := filepath.Abs(fi.Name()); err == nil {
			abs = _abs
		} else {
			continue
		}
		switch fi.IsDir() {
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
