package fod

import (
	"sync"

	termbox "github.com/nsf/termbox-go"
)

// Dialog execute dialog
func Dialog(opt Option) ([]string, ResultCode) {

	// create selector
	var selector Selector
	var err error
	if selector, err = newSelector(opt.Mode, opt.Multi); err != nil {
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
		runDialog(selector, opt.Base)
		wg.Done()
	}()
	wg.Wait()
	termbox.Close()

	return selector.result()
}

func runDialog(selector Selector, base string) {
	selector.changeDirectory(base)
Loop:
	for {
		draw(interface{}(selector).(DrawContext))

		event := termbox.PollEvent()
		// logger.Printf("%#v\n", event)

		// [todo] - key map config
		switch {
		case event.Type == termbox.EventKey && event.Key == termbox.KeyEnter:
			selector.changeDirectoryToCurrentItem()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowUp:
			selector.moveCursorUp()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowDown:
			selector.moveCursorDown()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowRight:
			selector.changeDirectoryToCurrentItem()
		case event.Type == termbox.EventKey && event.Key == termbox.KeyArrowLeft:
			selector.changeDirectoryUp()

		// filename filter
		case event.Ch >= 0x20 && event.Ch <= 0x7E:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeySpace:
			filenameFilterSingleton().addCharacter(event.Ch)
			selector.refresh()

		case event.Type == termbox.EventKey && event.Key == termbox.KeyBackspace:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeyBackspace2:
			filenameFilterSingleton().delCharacter()
			selector.refresh()

		// done
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlO:
			if selector.decide() {
				break Loop
			}

		// mark an item
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlS:
			selector.markItem()

		// toggle dotfile-filter
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlH:
			dotfileFilterSinleton().toggle()
			selector.refresh()

		// cancel and exit
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlQ:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeyCtrlC:
			fallthrough
		case event.Type == termbox.EventKey && event.Key == termbox.KeyEsc:
			selector.cancel()
			break Loop
		}
	}
}
