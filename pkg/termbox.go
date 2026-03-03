package fod

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// DrawContext : draw context
type DrawContext interface {
	// get entries
	getEntries() []*Entry

	// get current item index in entries
	getCurrentItemIndex() int

	// get number of items
	getTotalItems() int

	// get current directory path
	getPwd() string

	// get selector mode
	getMode() string

	// get filter string
	getFilterString() string
}

// draw a string
func drawString(x int, y int, text string, fgColor termbox.Attribute, bgColor termbox.Attribute) {
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		termbox.SetCell(x+i, y, runes[i], fgColor, bgColor)
	}
}

func marked(m bool) string {
	if m == true {
		return "*"
	}
	return " "
}

// draw menu
func drawEntries(dc DrawContext) {

	entries := dc.getEntries()
	currentIndex := dc.getCurrentItemIndex()

	_, h := termbox.Size()
	linePerPage := h - 3 // 1 (top) + 2 (bottom)

	currentPage := currentIndex / linePerPage
	pageTop := (currentPage * linePerPage)
	pageEnd := (currentPage * linePerPage) + linePerPage
	if pageEnd > len(entries) {
		pageEnd = len(entries)
	}
	cursorIndex := currentIndex - pageTop

	for index, entry := range entries[pageTop:pageEnd] {
		line := fmt.Sprintf(
			"[%s] %s %s",
			entry.typeCharcter(), marked(entry.Marked),
			entry.Path,
		)
		switch index == cursorIndex {
		case true:
			drawString(
				0, 1+index, line,
				termbox.ColorBlack,
				termbox.ColorWhite,
			)
		case false:
			drawString(
				0, 1+index, line,
				termbox.ColorDefault,
				termbox.ColorDefault,
			)
		}
	}
}

func drawStatusLineTop(dc DrawContext) {
	text := fmt.Sprintf("> %s", dc.getFilterString())
	drawString(
		0,
		0, text,
		termbox.ColorDefault,
		termbox.ColorDefault,
	)
}

// draw status line
func drawStatusLineBottom(dc DrawContext) {
	_, h := termbox.Size()
	var x int
	y1 := h - 2
	y2 := h - 1

	text1 := fmt.Sprintf(
		"select %3d of %3d items %s",
		dc.getCurrentItemIndex()+1,
		dc.getTotalItems(),
		dc.getPwd(),
	)
	drawString(
		x, y1, text1,
		termbox.ColorDefault,
		termbox.ColorDefault,
	)
	text2 := fmt.Sprintf("[%s] Ctrl+(O)K / Ctrl+(C)ancel, Ctrl+(Q)uit, Esc to exit", dc.getMode())
	drawString(
		x, y2, text2,
		termbox.ColorDefault,
		termbox.ColorDefault,
	)
}

// draw screen
func draw(dc DrawContext) {
	termbox.Clear(
		termbox.ColorDefault,
		termbox.ColorDefault,
	)
	drawStatusLineTop(dc)
	drawEntries(dc)
	drawStatusLineBottom(dc)
	termbox.Flush()
}
