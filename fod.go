package fod

import (
	"sync"

	termbox "github.com/nsf/termbox-go"
)

// Dialog execute dialog
func Dialog(opt Option) ([]string, ResultCode) {

	// create selector
	var selector Selector
	if selector, err := newSelector(opt.Mode, opt.Multi); err != nil {
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
		selector.run(opt.Base)
		wg.Done()
		return
	}()
	wg.Wait()
	termbox.Close()

	return selector.result()
}
