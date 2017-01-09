package fod

import (
	"errors"
	"os"
	"path/filepath"
)
import (
	"github.com/nsf/termbox-go"
)

// Selector interface of "*Selector"
type Selector interface {
	Result() (string, ResultCode)

	mark()
	markedItem() []string
	decide() bool
	cancel()

	changeDirectoryToCurrentItem()
	changeDirectoryUp()
	changeDirectory(path string) error

	moveCursorUp()
	moveCursorDown()
	refresh()
}

// SelectorFramework base type of selector
type SelectorFramework struct {
	concrete Selector
}

// newSelectorFramework create a instance of new selector
func newSelectorFramework(mode Mode, multi bool) (*SelectorFramework, error) {

	// create concrete selector
	var selector Selector
	var err error
	if selector, err = newSelector(mode, multi); err != nil {
		return nil, err
	}

	// create framework
	return &SelectorFramework{
		concrete: selector,
	}, nil
}

// Select select item
func (selector *SelectorFramework) Select(base string) {
	selector.concrete.changeDirectory(base)
Loop:
	for {
		draw(interface{}(selector.concrete).(DrawContext))

		event := termbox.PollEvent()
		// logger.Printf("%#v\n", event)

		// [todo] - key map config
		switch {
		case event.Type == termbox.EventKey && event.Key == termbox.KeyEnter:
			selector.concrete.changeDirectoryToCurrentItem()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowUp:
			selector.concrete.moveCursorUp()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowDown:
			selector.concrete.moveCursorDown()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowRight:
			selector.concrete.changeDirectoryToCurrentItem()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowLeft:
			selector.concrete.changeDirectoryUp()

		// filename filter
		case event.Ch >= 0x20 && event.Ch <= 0x7E:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeySpace:
			filenameFilterSingleton().addCharacter(event.Ch)
			selector.concrete.refresh()

		case event.Type == termbox.EventKey && event.Key == termbox.KeyBackspace:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeyBackspace2:
			filenameFilterSingleton().delCharacter()
			selector.concrete.refresh()

		// done
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlO:
			if selector.concrete.decide() {
				break Loop
			}

		// mark an item
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlS:
			selector.concrete.mark()

		// toggle dotfile-filter
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlH:
			dotfileFilterSinleton().toggle()
			selector.concrete.refresh()

		// cancel and exit
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlQ:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlC:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeyEsc:
			selector.concrete.cancel()
			break Loop
		}
	}
}

// Result get result
func (selector *SelectorFramework) result() ([]string, ResultCode) {
	m := selector.concrete.markedItem()
	if len(m) > 0 {
		return m, ResultOK
	}
	return nil, ResultCancel
}

// create a new selector
func newSelector(mode Mode, multi bool) (selector Selector, err error) {

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	switch mode {
	case ModeFile:
		selector = &FileSelector{
			&SelectorCommon{
				Multi:      multi,
				CurrentDir: dir,
				marked:     []string{},
				result:     "",
				resultCode: ResultNone,
				Cursor:     0,
				Filters: []Filter{
					fileFilterSingleton(),
					filenameFilterSingleton(),
					dotfileFilterSinleton(),
				},
			},
		}
	case ModeDirectory:
		selector = &DirectorySelector{
			&SelectorCommon{
				Multi:      multi,
				CurrentDir: dir,
				marked:     []string{},
				result:     "",
				resultCode: ResultNone,
				Cursor:     0,
				Filters: []Filter{
					directoryFilterSingleton(),
					filenameFilterSingleton(),
					dotfileFilterSinleton(),
				},
			},
		}
	default:
		err = errors.New("invalid selector mode")
	}
	return
}

// SelectorCommon : common implementation of selector
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
	selector.result = ""
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

// Result : get result
func (selector *SelectorCommon) Result() (result string, resultCode ResultCode) {
	return selector.result, selector.resultCode
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
