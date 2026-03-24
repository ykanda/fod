package fod

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/mattn/go-runewidth"
)

type keyTestSelector struct {
	decideReturn bool
	decideCalled int
}

func (s *keyTestSelector) result() ([]string, ResultCode) { return nil, ResultNone }
func (s *keyTestSelector) markItem()                      {}
func (s *keyTestSelector) markedItem() []string           { return nil }
func (s *keyTestSelector) decide() bool {
	s.decideCalled++
	return s.decideReturn
}
func (s *keyTestSelector) cancel()                        {}
func (s *keyTestSelector) changeDirectoryToCurrentItem()  {}
func (s *keyTestSelector) changeDirectoryUp()             {}
func (s *keyTestSelector) changeDirectory(path string) error {
	return nil
}
func (s *keyTestSelector) moveCursorUp()   {}
func (s *keyTestSelector) moveCursorDown() {}
func (s *keyTestSelector) refresh()        {}

type drawContextForHelp struct{}

func (d drawContextForHelp) getEntries() []*Entry {
	return []*Entry{{Path: "/tmp/example", Marked: true}}
}
func (d drawContextForHelp) getCurrentItemIndex() int { return 0 }
func (d drawContextForHelp) getTotalItems() int       { return 1 }
func (d drawContextForHelp) getPwd() string           { return "/tmp/example" }
func (d drawContextForHelp) getMode() string          { return "normal" }
func (d drawContextForHelp) getFilterString() string  { return "" }

func TestTruncateLine_DisplayWidth(t *testing.T) {
	text := "あああ"

	if got := truncateLine(text, 6); got != text {
		t.Fatalf("truncateLine(%q, 6) = %q, want %q", text, got, text)
	}

	if got := truncateLine(text, 5); got != "ああ" {
		t.Fatalf("truncateLine(%q, 5) = %q, want %q", text, got, "ああ")
	}

	if got := truncateLine(text, 5); runewidth.StringWidth(got) > 5 {
		t.Fatalf("truncateLine(%q, 5) width = %d, want <= 5", text, runewidth.StringWidth(got))
	}
}

func TestTruncateRunes_DisplayWidth(t *testing.T) {
	text := "あい"

	if got := truncateRunes(text, 4); got != text {
		t.Fatalf("truncateRunes(%q, 4) = %q, want %q", text, got, text)
	}

	if got := truncateRunes(text, 3); got != "あ" {
		t.Fatalf("truncateRunes(%q, 3) = %q, want %q", text, got, "あ")
	}

	if got := truncateRunes(text, 1); got != "" {
		t.Fatalf("truncateRunes(%q, 1) = %q, want %q", text, got, "")
	}
}

func TestHighlightMatches(t *testing.T) {
	text := "あああ.md"
	filter := "あああ"

	if got := highlightMatches(text, filter, false); got != sgrReverseOn+"あああ"+sgrReverseOff+".md" {
		t.Fatalf("highlightMatches(%q, %q) = %q", text, filter, got)
	}

	if got := highlightMatches(text, "zzz", false); got != text {
		t.Fatalf("highlightMatches(%q, %q) = %q, want %q", text, "zzz", got, text)
	}
}

func TestHandleKey_ShiftEnter_InNormalMode(t *testing.T) {
	selector := &keyTestSelector{decideReturn: true}
	model := dialogModel{
		selector: selector,
		mode:     modeNormal,
	}

	_, cmd := model.handleKey(tea.Key{Code: tea.KeyEnter, Mod: tea.ModShift})
	if selector.decideCalled != 1 {
		t.Fatalf("decide() called %d times, want 1", selector.decideCalled)
	}
	if cmd == nil {
		t.Fatal("cmd is nil, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("cmd() = %T, want tea.QuitMsg", cmd())
	}
}

func TestHandleKey_ShiftEnter_InFilterMode(t *testing.T) {
	selector := &keyTestSelector{decideReturn: true}
	model := dialogModel{
		selector: selector,
		mode:     modeFilter,
	}

	_, cmd := model.handleKey(tea.Key{Code: tea.KeyEnter, Mod: tea.ModShift})
	if selector.decideCalled != 1 {
		t.Fatalf("decide() called %d times, want 1", selector.decideCalled)
	}
	if cmd == nil {
		t.Fatal("cmd is nil, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("cmd() = %T, want tea.QuitMsg", cmd())
	}
}

func TestHandleKey_CtrlO_InNormalMode(t *testing.T) {
	selector := &keyTestSelector{decideReturn: true}
	model := dialogModel{
		selector: selector,
		mode:     modeNormal,
	}

	_, cmd := model.handleKey(tea.Key{Code: 'o', Mod: tea.ModCtrl})
	if selector.decideCalled != 1 {
		t.Fatalf("decide() called %d times, want 1", selector.decideCalled)
	}
	if cmd == nil {
		t.Fatal("cmd is nil, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("cmd() = %T, want tea.QuitMsg", cmd())
	}
}

func TestHandleKey_CtrlO_InFilterMode(t *testing.T) {
	selector := &keyTestSelector{decideReturn: true}
	model := dialogModel{
		selector: selector,
		mode:     modeFilter,
	}

	_, cmd := model.handleKey(tea.Key{Code: 'o', Mod: tea.ModCtrl})
	if selector.decideCalled != 1 {
		t.Fatalf("decide() called %d times, want 1", selector.decideCalled)
	}
	if cmd == nil {
		t.Fatal("cmd is nil, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("cmd() = %T, want tea.QuitMsg", cmd())
	}
}

func TestBuildView_HelpIncludesShiftEnterAndCtrlO(t *testing.T) {
	view := buildView(drawContextForHelp{}, 120, 20, modeNormal, true)
	if want := "Shift+Enter, Ctrl+O"; !strings.Contains(view, want) {
		t.Fatalf("view does not include %q", want)
	}
}
