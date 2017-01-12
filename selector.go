package fod

import (
	"errors"
	"os"
)

// Selector interface of "*Selector"
type Selector interface {
	result() ([]string, ResultCode)

	markItem()
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
