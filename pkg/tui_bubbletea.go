package fod

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
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

type dialogModel struct {
	selector Selector
	width    int
	height   int
}

func newDialogModel(selector Selector, base string) dialogModel {
	selector.changeDirectory(base)
	return dialogModel{
		selector: selector,
		width:    defaultWidth,
		height:   defaultHeight,
	}
}

func (m dialogModel) Init() tea.Cmd {
	return nil
}

func (m dialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > 0 {
			m.width = msg.Width
		}
		if msg.Height > 0 {
			m.height = msg.Height
		}
		return m, nil
	case tea.KeyPressMsg:
		return m.handleKey(msg.Key())
	}
	return m, nil
}

func (m dialogModel) handleKey(key tea.Key) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "enter", "right":
		m.selector.changeDirectoryToCurrentItem()
		return m, nil
	case "up":
		m.selector.moveCursorUp()
		return m, nil
	case "down":
		m.selector.moveCursorDown()
		return m, nil
	case "left":
		m.selector.changeDirectoryUp()
		return m, nil
	case "ctrl+o":
		if m.selector.decide() {
			return m, tea.Quit
		}
		return m, nil
	case "ctrl+s":
		m.selector.markItem()
		return m, nil
	case "ctrl+h":
		dotfileFilterSinleton().toggle()
		m.selector.refresh()
		return m, nil
	case "ctrl+q", "ctrl+c", "esc":
		m.selector.cancel()
		return m, tea.Quit
	case "backspace":
		filenameFilterSingleton().delCharacter()
		m.selector.refresh()
		return m, nil
	}

	if key.Text != "" && key.Mod == 0 {
		for _, r := range []rune(key.Text) {
			filenameFilterSingleton().addCharacter(r)
		}
		m.selector.refresh()
	}

	return m, nil
}

func (m dialogModel) View() tea.View {
	dc, ok := interface{}(m.selector).(DrawContext)
	if !ok {
		return tea.NewView("")
	}
	view := tea.NewView(buildView(dc, m.width, m.height))
	view.AltScreen = true
	return view
}

func buildView(dc DrawContext, width int, height int) string {
	width, height = normalizeSize(width, height)

	linePerPage := height - 3 // 1 (top) + 2 (bottom)
	if linePerPage < 1 {
		linePerPage = 1
	}

	entries := dc.getEntries()
	currentIndex := dc.getCurrentItemIndex()
	currentPage := 0
	if linePerPage > 0 {
		currentPage = currentIndex / linePerPage
	}
	pageTop := currentPage * linePerPage
	pageEnd := pageTop + linePerPage
	if pageEnd > len(entries) {
		pageEnd = len(entries)
	}
	cursorIndex := currentIndex - pageTop

	var lines []string
	lines = append(lines, truncateLine(fmt.Sprintf("> %s", dc.getFilterString()), width))

	highlight := lipgloss.NewStyle().Reverse(true).Width(width)
	normal := lipgloss.NewStyle().Width(width)

	for index, entry := range entries[pageTop:pageEnd] {
		line := fmt.Sprintf("[%s] %s %s", entry.typeCharcter(), marked(entry.Marked), entry.Path)
		line = truncateLine(line, width)
		if index == cursorIndex {
			lines = append(lines, highlight.Render(line))
		} else {
			lines = append(lines, normal.Render(line))
		}
	}

	status1 := fmt.Sprintf(
		"select %3d of %3d items %s",
		dc.getCurrentItemIndex()+1,
		dc.getTotalItems(),
		dc.getPwd(),
	)
	status2 := fmt.Sprintf("[%s] Ctrl+(O)K / Ctrl+(C)ancel, Ctrl+(Q)uit, Esc to exit", dc.getMode())
	lines = append(lines, truncateLine(status1, width))
	lines = append(lines, truncateLine(status2, width))

	return strings.Join(lines, "\n")
}

func marked(m bool) string {
	if m == true {
		return "*"
	}
	return " "
}

func normalizeSize(width int, height int) (int, int) {
	if width <= 0 {
		width = defaultWidth
	}
	if height <= 0 {
		height = defaultHeight
	}
	if width < 20 {
		width = 20
	}
	if height < 5 {
		height = 5
	}
	return width, height
}

func truncateLine(text string, width int) string {
	if width <= 0 {
		return text
	}
	runes := []rune(text)
	if len(runes) <= width {
		return text
	}
	return string(runes[:width])
}
