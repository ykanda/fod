package main

import (
	"errors"
	"os"
	"path/filepath"
)
import (
	"github.com/nsf/termbox-go"
)

type Selector interface {
	Result() (string, ResultCode)
	Choice() bool
	Cancel()

	ChangeDirectoryToCurrentItem()
	ChangeDirectoryUp()
	ChangeDirectory(path string)

	MoveCursorUp()
	MoveCursorDown()
	Refresh()
}

//
type SelectorFramework struct {
	concrete Selector
}

// create new SelectorFramework
func NewSelectorFramework(mode Mode) (*SelectorFramework, error) {

	// create concrete selector
	var selector Selector
	if _selector, err := NewSelector(mode); err == nil {
		selector = _selector
	} else {
		return nil, err
	}

	// create framework
	return &SelectorFramework{
		concrete: selector,
	}, nil
}

// select
func (self *SelectorFramework) Select(base string) {
	self.concrete.ChangeDirectory(base)
Loop:
	for {
		Draw(interface{}(self.concrete).(DrawContext))

		event := termbox.PollEvent()
		DebugLog(event)

		// [todo] - key map config
		switch {
		case event.Type == termbox.EventKey && event.Key == termbox.KeyEnter:
			self.concrete.ChangeDirectoryToCurrentItem()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowUp:
			self.concrete.MoveCursorUp()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowDown:
			self.concrete.MoveCursorDown()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowRight:
			self.concrete.ChangeDirectoryToCurrentItem()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowLeft:
			self.concrete.ChangeDirectoryUp()

		// filename filter
		case event.Ch >= 0x20 && event.Ch <= 0x7E:
			FilenameFilterSingleton().AddCharacter(event.Ch)
			self.concrete.Refresh()
		case event.Type == termbox.EventKey && event.Key == 0x007F:
			FilenameFilterSingleton().DelCharacter()
			self.concrete.Refresh()

		// Chose
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlO:
			if self.concrete.Choice() {
				break Loop
			}

		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlQ:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlC:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeyEsc:
			self.concrete.Cancel()
			break Loop
		}
	}
}

// get result
func (self *SelectorFramework) Result() (string, ResultCode) {
	return self.concrete.Result()
}

// ----------------------------------------------------------------------------

// create a new selector
func NewSelector(mode Mode) (selector Selector, err error) {
	switch mode {
	case MODE_FILE:
		selector = &FileSelector{
			&SelectorCommon{
				CurrentDir: "./",
				result:     "",
				resultCode: RESULT_NONE,
				Cursor:     0,
				Filters: []Filter{
					FileFilterSingleton(),
					FilenameFilterSingleton(),
					DotfileFilterSinleton(),
				},
			},
		}
	case MODE_DIRECTORY:
		selector = &DirectorySelector{
			&SelectorCommon{
				CurrentDir: "./",
				result:     "",
				resultCode: RESULT_NONE,
				Cursor:     0,
				Filters: []Filter{
					DirectoryFilterSingleton(),
					FilenameFilterSingleton(),
					DotfileFilterSinleton(),
				},
			},
		}
	default:
		err = errors.New("invalid selector mode")
	}
	return
}

//-----------------------------------------------------------------------------

type SelectorCommon struct {
	CurrentDir string
	Cursor     int
	Entries    []Entry
	Filters    []Filter
	result     string
	resultCode ResultCode
}

// move cursor up
func (self *SelectorCommon) MoveCursorUp() {
	if self.Cursor > 0 {
		self.Cursor--
	}
}

// move cursor down
func (self *SelectorCommon) MoveCursorDown() {
	if self.Cursor < (len(self.GetEntries()) - 1) {
		self.Cursor++
	}
}

// set result
func (self *SelectorCommon) Cancel() {
	self.result = ""
	self.resultCode = RESULT_CANCEL
}

// refresh
func (self *SelectorCommon) Refresh() {
	self.Cursor = 0
}

// get focused item (absolute) path
func (self *SelectorCommon) CurrentItem() (string, error) {
	entries := self.GetEntries()
	path := self.CurrentDir + "/" + entries[self.Cursor].Path
	return filepath.Abs(path)
}

// get focused item type (file or directory)
func (self *SelectorCommon) CurrentItemType(fsType string) (is bool) {

	var path string = ""
	if _path, err := self.CurrentItem(); err == nil {
		path = _path
	} else {
		DebugLog(err)
	}

	if fi, err := os.Stat(path); err == nil {
		DebugLog(fi)
		switch {
		case (FS_TYPE_DIR == fsType) && (fi.IsDir() == true):
			is = true
		case (FS_TYPE_FILE == fsType) && (fi.IsDir() == false):
			is = true
		}
	} else {
		DebugLog(err)
	}
	return
}

// current dir is root?
func (self *SelectorCommon) CurrentDirIsRoot() bool {
	return self.CurrentDir == "/"
}

// change current dir
func (self *SelectorCommon) ChangeDirectory(path string) {
	if filepath.IsAbs(path) {
		self.CurrentDir = path
		self.Cursor = 0
		self.Entries = Entries(self.CurrentDir)
	} else if abs, err := filepath.Abs(self.CurrentDir + "/" + path); err == nil {
		self.CurrentDir = abs
		self.Cursor = 0
		self.Entries = Entries(self.CurrentDir)
	}

	DebugLog(self)
}

// change directory up
func (self *SelectorCommon) ChangeDirectoryUp() {
	self.ChangeDirectory("../")
}

// change current dir
func (self *SelectorCommon) ChangeDirectoryToCurrentItem() {
	if self.CurrentItemType(FS_TYPE_DIR) == false {
		DebugLog("current item is not directory")
		return
	}
	if targetDir, err := self.CurrentItem(); err == nil {
		DebugLog(err)
		self.ChangeDirectory(targetDir)
	}
}

// get result
func (self *SelectorCommon) Result() (result string, resultCode ResultCode) {
	return self.result, self.resultCode
}

// entries
func (self *SelectorCommon) GetEntries() []Entry {
	entries := self.Entries
	for _, f := range self.Filters {
		entries = f.Filter(entries)
	}
	return entries
}

// return cursor pos
func (self *SelectorCommon) GetCurrentItemIndex() int {
	return self.Cursor
}

// return entries num
func (self *SelectorCommon) GetTotalItems() int {
	return len(self.GetEntries())
}

// get current directory path
func (self *SelectorCommon) GetPwd() string {
	if path, err := self.CurrentItem(); err == nil {
		return path
	}
	return ""
}

// get filter string
func (self *SelectorCommon) GetFilterString() string {
	return FilenameFilterSingleton().GetFilterString()
}
