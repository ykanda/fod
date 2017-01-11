package fod

import (
	"os"
	"path/filepath"
)

// SelectorCommon : common implementation of selector
type SelectorCommon struct {
	Multi      bool
	CurrentDir string
	Cursor     int
	Entries    []*Entry
	Filters    []Filter
	marked     []string
	resultCode ResultCode
}

func (selector *SelectorCommon) getMode() string {
	return "hogehoge"
}

// Result get result
func (selector *SelectorCommon) result() ([]string, ResultCode) {
	m := selector.markedItem()
	if len(m) > 0 {
		return m, ResultOK
	}
	return nil, ResultCancel
}

// moveCursorUp : move cursor up
func (selector *SelectorCommon) moveCursorUp() {
	if selector.Cursor > 0 {
		selector.Cursor--
	}
}

// moveCursorDown : move cursor down
func (selector *SelectorCommon) moveCursorDown() {
	if selector.Cursor < (len(selector.getEntries()) - 1) {
		selector.Cursor++
	}
}

// cancel : set result
func (selector *SelectorCommon) cancel() {
	selector.marked = nil
	selector.resultCode = ResultCancel
}

// refresh : refresh cursor position
func (selector *SelectorCommon) refresh() {
	selector.Cursor = 0
}

// currentItem : get focused item (absolute) path
func (selector *SelectorCommon) currentItem() (string, error) {
	entries := selector.getEntries()
	return entries[selector.Cursor].Path, nil
}

// currentItemType : get focused item type (file or directory)
func (selector *SelectorCommon) currentItemType(fsType string) (is bool) {

	var path string
	if _path, err := selector.currentItem(); err == nil {
		path = _path
	}

	if fi, err := os.Stat(path); err == nil {
		switch {
		case (FsTypeDir == fsType) && (fi.IsDir() == true):
			is = true
		case (FsTypeFile == fsType) && (fi.IsDir() == false):
			is = true
		}
	} else {
		logger.Println(err)
	}
	return
}

// [todo] - ターゲットとするOSごとに処理を分ける
// current dir is root?
func (selector *SelectorCommon) currentDirIsRoot() bool {
	return selector.CurrentDir == "/"
}

// changeDirectory : change current dir
func (selector *SelectorCommon) changeDirectory(path string) error {
	if filepath.IsAbs(path) {
		selector.CurrentDir = path
		selector.Cursor = 0
		selector.Entries = entries(selector.CurrentDir)
		return nil
	}
	abs, err := filepath.Abs(filepath.Join(selector.CurrentDir, path))
	if err == nil {
		selector.CurrentDir = abs
		selector.Cursor = 0
		selector.Entries = entries(selector.CurrentDir)
	}
	return err
}

// changeDirectoryUp : change current directory to parent
func (selector *SelectorCommon) changeDirectoryUp() {
	selector.changeDirectory("../")
}

// changeDirectoryToCurrentItem : change current dir
func (selector *SelectorCommon) changeDirectoryToCurrentItem() {
	if selector.currentItemType(FsTypeDir) == false {
		logger.Println("current item is not directory")
		return
	}
	if targetDir, err := selector.currentItem(); err == nil {
		logger.Printf("%s\n", err)
		selector.changeDirectory(targetDir)
	}
}

// getEntries : get entries
func (selector *SelectorCommon) getEntries() []*Entry {
	entries := selector.Entries
	for _, f := range selector.Filters {
		entries = f.filter(entries)
	}
	return entries
}

// getCurrentItemIndex : return cursor position
func (selector *SelectorCommon) getCurrentItemIndex() int {
	return selector.Cursor
}

// getTotalItems : return entries num
func (selector *SelectorCommon) getTotalItems() int {
	return len(selector.getEntries())
}

// getPwd get current directory path
func (selector *SelectorCommon) getPwd() string {
	if path, err := selector.currentItem(); err == nil {
		return path
	}
	return ""
}

// getFilterString get filter string
func (selector *SelectorCommon) getFilterString() string {
	return filenameFilterSingleton().getFilterString()
}

func (selector *SelectorCommon) markItem() {
	var path string
	if p, err := selector.currentItem(); err == nil {
		path = p
	}

	switch selector.Multi {
	case true:
		selector.toggleItem(path, selector.getCurrentItemIndex())
	case false:
		selector.setItem(path, selector.getCurrentItemIndex())
	}
	logger.Printf("%#v\n", selector.marked)
}

func filter(a []string, f func(string) bool) []string {
	n := []string{}
	for _, x := range a {
		if f(x) {
			n = append(n, x)
		}
	}
	return n
}

func contains(a []string, f func(string) bool) []int {
	var n = []int{}
	for i, x := range a {
		if f(x) {
			n = append(n, i)
		}
	}
	return n
}

func (selector *SelectorCommon) toggleItem(path string, index int) {
	exists := 0 < len(
		contains(selector.marked, func(elem string) bool {
			return (elem == path)
		}))

	entries := selector.getEntries()
	if exists {
		entries[index].Marked = false
		selector.marked = filter(selector.marked, func(elem string) bool {
			return (elem != path)
		})
	} else {
		entries[index].Marked = true
		selector.marked = append(selector.marked, path)
	}
}

func (selector *SelectorCommon) setItem(path string, index int) {
	for i := range selector.Entries {
		selector.Entries[i].Marked = false
	}
	entries := selector.getEntries()
	entries[index].Marked = true
	selector.marked = []string{path}
	logger.Printf("%#v\n", selector.Entries)
}

func (selector *SelectorCommon) decide() bool {
	return len(selector.marked) > 0
}

func (selector *SelectorCommon) markedItem() []string {
	return selector.marked
}
