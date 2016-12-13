package fod

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)
import (
	"github.com/nsf/termbox-go"
)

type Selector interface {
	Result() (string, ResultCode)
	Mark()
	Decide() bool
	Cancel()

	ChangeDirectoryToCurrentItem()
	ChangeDirectoryUp()
	ChangeDirectory(path string) error

	MoveCursorUp()
	MoveCursorDown()
	Refresh()
}

//
type SelectorFramework struct {
	concrete Selector
}

// create new SelectorFramework
func NewSelectorFramework(mode Mode, multi bool) (*SelectorFramework, error) {

	// create concrete selector
	var selector Selector
	if _selector, err := NewSelector(mode, multi); err == nil {
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
		// logger.Printf("%#v\n", event)

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
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeySpace:
			FilenameFilterSingleton().AddCharacter(event.Ch)
			self.concrete.Refresh()

		case event.Type == termbox.EventKey && event.Key == termbox.KeyBackspace:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeyBackspace2:
			FilenameFilterSingleton().DelCharacter()
			self.concrete.Refresh()

		// done
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlO:
			if self.concrete.Decide() {
				break Loop
			}

		// mark an item
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlS:
			self.concrete.Mark()

		// toggle dotfile-filter
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlH:
			DotfileFilterSinleton().toggle()
			self.concrete.Refresh()

		// cancel and exit
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
func NewSelector(mode Mode, multi bool) (selector Selector, err error) {

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	switch mode {
	case MODE_FILE:
		selector = &FileSelector{
			&SelectorCommon{
				Multi:      multi,
				CurrentDir: dir,
				marked:     []string{},
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
				Multi:      multi,
				CurrentDir: dir,
				marked:     []string{},
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
	Multi      bool
	CurrentDir string
	Cursor     int
	Entries    []*Entry
	Filters    []Filter
	marked     []string
	result     string
	resultCode ResultCode
}

// MoveCursorUp move cursor up
func (selector *SelectorCommon) MoveCursorUp() {
	if selector.Cursor > 0 {
		selector.Cursor--
	}
}

// MoveCursorDown move cursor down
func (selector *SelectorCommon) MoveCursorDown() {
	if selector.Cursor < (len(selector.GetEntries()) - 1) {
		selector.Cursor++
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
	return entries[self.Cursor].Path, nil
}

// get focused item type (file or directory)
func (self *SelectorCommon) CurrentItemType(fsType string) (is bool) {

	var path string = ""
	if _path, err := self.CurrentItem(); err == nil {
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
func (self *SelectorCommon) CurrentDirIsRoot() bool {
	return self.CurrentDir == "/"
}

// change current dir
func (self *SelectorCommon) ChangeDirectory(path string) error {
	if filepath.IsAbs(path) {
		self.CurrentDir = path
		self.Cursor = 0
		self.Entries = Entries(self.CurrentDir)
		return nil
	}
	abs, err := filepath.Abs(filepath.Join(self.CurrentDir, path))
	if err == nil {
		self.CurrentDir = abs
		self.Cursor = 0
		self.Entries = Entries(self.CurrentDir)
	}
	return err
}

// change directory up
func (self *SelectorCommon) ChangeDirectoryUp() {
	self.ChangeDirectory("../")
}

// change current dir
func (self *SelectorCommon) ChangeDirectoryToCurrentItem() {
	if self.CurrentItemType(FsTypeDir) == false {
		logger.Println("current item is not directory")
		return
	}
	if targetDir, err := self.CurrentItem(); err == nil {
		logger.Printf("%s\n", err)
		self.ChangeDirectory(targetDir)
	}
}

// get result
func (self *SelectorCommon) Result() (result string, resultCode ResultCode) {
	return self.result, self.resultCode
}

// entries
func (self *SelectorCommon) GetEntries() []*Entry {
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

func (self *SelectorCommon) MarkItem() {
	var path string = ""
	if _path, err := self.CurrentItem(); err != nil {
		return
	} else {
		path = _path
	}
	switch self.Multi {
	case true:
		self.ToggleItem(path, self.GetCurrentItemIndex())
	case false:
		self.SetItem(path, self.GetCurrentItemIndex())
	}
	logger.Printf("%#v\n", self.marked)
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

func (self *SelectorCommon) ToggleItem(path string, index int) {
	exists := 0 < len(
		contains(self.marked, func(elem string) bool {
			return (elem == path)
		}))

	entries := self.GetEntries()
	if exists {
		entries[index].Marked = false
		self.marked = filter(self.marked, func(elem string) bool {
			return (elem != path)
		})
	} else {
		entries[index].Marked = true
		self.marked = append(self.marked, path)
	}
}

func (selector *SelectorCommon) SetItem(path string, index int) {
	for i := range selector.Entries {
		selector.Entries[i].Marked = false
	}
	entries := selector.GetEntries()
	entries[index].Marked = true
	selector.marked = []string{path}
	logger.Printf("%#v\n", selector.Entries)
}

func (self *SelectorCommon) Decide() (selected bool) {
	if len(self.marked) > 0 {
		selected = true
		self.result = strings.Join(self.marked, " ")
		self.resultCode = RESULT_OK
	} else {
		self.result = ""
		self.resultCode = RESULT_CANCEL
	}
	return selected
}
