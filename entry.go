package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Entry struct {
	Path string
	Type string
}

// alias to wditems
func Entries(path string) []Entry {

	// list directory entries
	readdir, err := ioutil.ReadDir(path)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	entries := []Entry{}
	entries = append(entries, Entry{
		Path: "../",
		Type: FS_TYPE_DIR,
	})
	for _, fi := range readdir {
		switch fi.IsDir() {
		case true:
			entries = append(entries, Entry{
				Path: fi.Name(),
				Type: FS_TYPE_DIR,
			})
		case false:
			entries = append(entries, Entry{
				Path: fi.Name(),
				Type: FS_TYPE_FILE,
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
