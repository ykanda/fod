package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Entry struct {
	Path   string
	Type   string
	Marked bool
}

// alias to wditems
func Entries(path string) []*Entry {

	// list directory entries
	readdir, err := ioutil.ReadDir(path)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	entries := []*Entry{}
	entries = append(entries, &Entry{
		Path: "../",
		Type: FS_TYPE_DIR,
	})
	for _, fi := range readdir {
		var abs string = ""
		if _abs, err := filepath.Abs(fi.Name()); err == nil {
			abs = _abs
		} else {
			continue
		}
		switch fi.IsDir() {
		case true:
			entries = append(entries, &Entry{
				Path:   abs,
				Type:   FS_TYPE_DIR,
				Marked: false,
			})
		case false:
			entries = append(entries, &Entry{
				Path:   abs,
				Type:   FS_TYPE_FILE,
				Marked: false,
			})
		}
	}
	return entries
}

// get type character
// FS_TYPE_DIR  -> d
// FS_TYPE_FILE -> -
func (self *Entry) TypeCharcter() (tc string) {
	switch self.Type {
	case FS_TYPE_DIR:
		tc = "d"
	case FS_TYPE_FILE:
		tc = "-"
	}
	return
}
