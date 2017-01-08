package fod

import (
	"sync"

	termbox "github.com/nsf/termbox-go"
)

// Dialog execute dialog
func Dialog(opt Option) ([]string, ResultCode) {

	// create selector
	var selector *SelectorFramework
	if _selector, err := newSelectorFramework(opt.Mode, opt.Multi); err == nil {
		selector = _selector
	} else {
		return []string{}, ResultNone
	}

	// init termbox
	if err := termbox.Init(); err != nil {
		return []string{}, ResultNone
	}
	termbox.SetInputMode(termbox.InputEsc)

	// select loop
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		selector.Select(opt.Base)
		wg.Done()
		return
	}()
	wg.Wait()
	termbox.Close()

	return selector.result()
}
