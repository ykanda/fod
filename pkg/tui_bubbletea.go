package fod

import (
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mattn/go-runewidth"
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
	mode     inputMode
	showHelp bool
}

type inputMode int

const (
	modeNormal inputMode = iota
	modeFilter
)

func newDialogModel(selector Selector, base string) dialogModel {
	selector.changeDirectory(base)
	return dialogModel{
		selector: selector,
		width:    defaultWidth,
		height:   defaultHeight,
		mode:     modeNormal,
		showHelp: false,
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
	if key.String() == "?" {
		m.showHelp = !m.showHelp
		return m, nil
	}
	switch m.mode {
	case modeNormal:
		switch key.String() {
		case "ctrl+f":
			m.mode = modeFilter
			return m, nil
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
		case "shift+enter", "ctrl+o":
			if m.selector.decide() {
				return m, tea.Quit
			}
			return m, nil
		case "space":
			m.selector.markItem()
			return m, nil
		case "ctrl+a":
			m.selector.selectAll()
			return m, nil
		case "ctrl+d":
			m.selector.clearAll()
			return m, nil
		case "ctrl+h":
			dotfileFilterSinleton().toggle()
			m.selector.refresh()
			return m, nil
		case "ctrl+q", "ctrl+c", "esc":
			m.selector.cancel()
			return m, tea.Quit
		}
	case modeFilter:
		switch key.String() {
		case "esc":
			m.mode = modeNormal
			return m, nil
		case "ctrl+q", "ctrl+c":
			m.selector.cancel()
			return m, tea.Quit
		case "shift+enter", "ctrl+o":
			if m.selector.decide() {
				return m, tea.Quit
			}
			return m, nil
		case "backspace":
			filenameFilterSingleton().delCharacter()
			m.selector.refresh()
			return m, nil
		}
		if key.Text != "" && (key.Mod == 0 || key.Mod == tea.ModShift) {
			for _, r := range []rune(key.Text) {
				filenameFilterSingleton().addCharacter(r)
			}
			m.selector.refresh()
		}
	}

	return m, nil
}

func (m dialogModel) View() tea.View {
	dc, ok := interface{}(m.selector).(DrawContext)
	if !ok {
		return tea.NewView("")
	}
	view := tea.NewView(buildView(dc, m.width, m.height, m.mode, m.showHelp))
	view.AltScreen = true
	return view
}

func buildView(dc DrawContext, width int, height int, mode inputMode, showHelp bool) string {
	width, height = normalizeSize(width, height)

	helpLines := []string{
		truncateLine(fmt.Sprintf("[%s] ? help", dc.getMode()), width),
	}
	if showHelp {
		helpLines = []string{
			truncateLine(formatHelpLine2("Ctrl+O", "quit, output selected items (or Shift + Enter)"), width),
			truncateLine(formatHelpLine2("Ctrl+Q", "quit, no output (or Esc)"), width),
			truncateLine(formatHelpLine2("Ctrl+F", "filter mode (Esc to exit)"), width),
			truncateLine(formatHelpLine2("Ctrl+H", "toggle hidden file filter"), width),
			truncateLine(formatHelpLine2("Ctrl+A", "select all displayed items (--multi only)"), width),
			truncateLine(formatHelpLine2("Ctrl+D", "clear all selected items (--multi only)"), width),
			truncateLine(formatHelpLine2("space", "select/unselect item"), width),
			truncateLine(formatHelpLine2("↑", "move cursor up"), width),
			truncateLine(formatHelpLine2("↓", "move cursor down"), width),
			truncateLine(formatHelpLine2("←", "go to parent directory"), width),
			truncateLine(formatHelpLine2("→", "enter the directory at the cursor"), width),
			truncateLine(formatHelpLine2("?", "toggle Help"), width),
		}
	}
	availableBottom := height - 2 // 1 (top) + 1 (status)
	if availableBottom < 0 {
		availableBottom = 0
	}
	linePerPage := availableBottom - len(helpLines)
	if linePerPage < 1 {
		linePerPage = 1
	}
	helpLinesCount := availableBottom - linePerPage
	if helpLinesCount < 0 {
		helpLinesCount = 0
	}
	if helpLinesCount > len(helpLines) {
		helpLinesCount = len(helpLines)
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
	filterRaw := dc.getFilterString()
	filterDisplay := filterRaw
	if mode == modeFilter {
		filterDisplay = fmt.Sprintf("[%s]", filterDisplay)
	}
	lines = append(lines, truncateLine(fmt.Sprintf("> %s", filterDisplay), width))

	highlight := lipgloss.NewStyle().Underline(true).Width(width)
	normal := lipgloss.NewStyle().Width(width)

	for index, entry := range entries[pageTop:pageEnd] {
		prefix := fmt.Sprintf("[%s] %s ", entry.typeCharcter(), marked(entry.Marked))
		available := width - runewidth.StringWidth(prefix)
		var line string
		if available <= 0 {
			line = truncateRunes(prefix, width)
			if index == cursorIndex {
				lines = append(lines, highlight.Render(line))
			} else {
				lines = append(lines, normal.Render(line))
			}
			continue
		}

		path := truncateRunes(entry.Path, available)
		if filterRaw != "" {
			path = highlightMatches(path, filterRaw, filenameFilterSingleton().ignoreCase)
		}
		line = prefix + path
		if index == cursorIndex {
			lines = append(lines, highlight.Render(line))
		} else {
			lines = append(lines, normal.Render(line))
		}
	}

	for i := pageEnd - pageTop; i < linePerPage; i++ {
		lines = append(lines, "")
	}

	status1 := fmt.Sprintf(
		"select %3d of %3d items %s",
		dc.getCurrentItemIndex()+1,
		dc.getTotalItems(),
		dc.getPwd(),
	)
	statusLine := lipgloss.NewStyle().Reverse(true).Width(width).Render(truncateLine(status1, width))
	lines = append(lines, statusLine)
	for i := 0; i < helpLinesCount; i++ {
		lines = append(lines, helpLines[i])
	}

	return strings.Join(lines, "\n")
}

func formatHelpLine2(left string, right string) string {
	const leftWidth = 12
	return fmt.Sprintf("%-*s %s", leftWidth, left, right)
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
	if runewidth.StringWidth(text) <= width {
		return text
	}
	return runewidth.Truncate(text, width, "")
}

func truncateRunes(text string, width int) string {
	if width <= 0 {
		return ""
	}
	if runewidth.StringWidth(text) <= width {
		return text
	}
	return runewidth.Truncate(text, width, "")
}

const (
	sgrReverseOn  = "\x1b[7m"
	sgrReverseOff = "\x1b[27m"
)

func highlightMatches(text string, filter string, ignoreCase bool) string {
	words := strings.Fields(filter)
	if len(words) == 0 || text == "" {
		return text
	}
	ranges := findMatchRanges(text, words, ignoreCase)
	if len(ranges) == 0 {
		return text
	}

	runes := []rune(text)
	var b strings.Builder
	pos := 0
	for _, r := range ranges {
		if r[0] > pos {
			b.WriteString(string(runes[pos:r[0]]))
		}
		if r[1] > r[0] {
			b.WriteString(sgrReverseOn)
			b.WriteString(string(runes[r[0]:r[1]]))
			b.WriteString(sgrReverseOff)
		}
		pos = r[1]
		logger.Printf("b.String(): %#v, pos: %d\n", b.String(), pos)
	}
	if pos < len(runes) {
		b.WriteString(string(runes[pos:]))
	}
	logger.Printf("b.String(): %#v\n", b.String())
	return b.String()
}

func findMatchRanges(text string, words []string, ignoreCase bool) [][2]int {
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}
	compareRunes := runes
	if ignoreCase {
		compareRunes = []rune(strings.ToLower(text))
	}

	var ranges [][2]int
	for _, word := range words {
		if word == "" {
			continue
		}
		wordRunes := []rune(word)
		if ignoreCase {
			wordRunes = []rune(strings.ToLower(word))
		}
		if len(wordRunes) == 0 || len(wordRunes) > len(compareRunes) {
			continue
		}
		for i := 0; i <= len(compareRunes)-len(wordRunes); i++ {
			match := true
			for j := 0; j < len(wordRunes); j++ {
				if compareRunes[i+j] != wordRunes[j] {
					match = false
					break
				}
			}
			if match {
				ranges = append(ranges, [2]int{i, i + len(wordRunes)})
				i += len(wordRunes) - 1
			}
		}
	}

	if len(ranges) == 0 {
		return nil
	}
	sort.Slice(ranges, func(i, j int) bool {
		if ranges[i][0] == ranges[j][0] {
			return ranges[i][1] < ranges[j][1]
		}
		return ranges[i][0] < ranges[j][0]
	})

	merged := make([][2]int, 0, len(ranges))
	for _, r := range ranges {
		if len(merged) == 0 {
			merged = append(merged, r)
			continue
		}
		last := &merged[len(merged)-1]
		if r[0] > (*last)[1] {
			merged = append(merged, r)
			continue
		}
		if r[1] > (*last)[1] {
			(*last)[1] = r[1]
		}
	}

	return merged
}
