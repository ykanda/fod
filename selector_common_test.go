package fod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSelectorCommonMoveCursor(t *testing.T) {
	sel := &SelectorCommon{
		Entries: []*Entry{{Path: "a"}, {Path: "b"}},
		Cursor:  0,
	}

	sel.moveCursorUp()
	if sel.Cursor != 0 {
		t.Fatalf("moveCursorUp() cursor = %d, want 0", sel.Cursor)
	}

	sel.moveCursorDown()
	if sel.Cursor != 1 {
		t.Fatalf("moveCursorDown() cursor = %d, want 1", sel.Cursor)
	}

	sel.moveCursorDown()
	if sel.Cursor != 1 {
		t.Fatalf("moveCursorDown() at end cursor = %d, want 1", sel.Cursor)
	}
}

func TestSelectorCommonChangeDirectory(t *testing.T) {
	root, err := os.MkdirTemp("", "fod-selector-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	defer os.RemoveAll(root)

	sub := filepath.Join(root, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	sel := &SelectorCommon{CurrentDir: root}
	if err := sel.changeDirectory("sub"); err != nil {
		t.Fatalf("changeDirectory: %v", err)
	}
	absSub, _ := filepath.Abs(sub)
	if sel.CurrentDir != absSub {
		t.Fatalf("CurrentDir = %q, want %q", sel.CurrentDir, absSub)
	}
	if sel.Cursor != 0 {
		t.Fatalf("Cursor = %d, want 0", sel.Cursor)
	}
	if len(sel.Entries) == 0 {
		t.Fatalf("Entries should not be empty")
	}

	sel.changeDirectoryUp()
	absRoot, _ := filepath.Abs(root)
	if sel.CurrentDir != absRoot {
		t.Fatalf("CurrentDir after up = %q, want %q", sel.CurrentDir, absRoot)
	}
}

func TestSelectorCommonMarkItemToggle(t *testing.T) {
	root, err := os.MkdirTemp("", "fod-selector-mark-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	defer os.RemoveAll(root)

	filePath := filepath.Join(root, "file.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	sel := &SelectorCommon{
		Multi:   true,
		Cursor:  0,
		Entries: []*Entry{{Path: filePath, Type: FsTypeFile}},
	}

	sel.markItem()
	if len(sel.marked) != 1 || sel.marked[0] != filePath {
		t.Fatalf("markItem() marked = %#v, want %q", sel.marked, filePath)
	}
	if !sel.Entries[0].Marked {
		t.Fatalf("markItem() should mark entry")
	}

	sel.markItem()
	if len(sel.marked) != 0 {
		t.Fatalf("markItem() toggle marked = %#v, want empty", sel.marked)
	}
	if sel.Entries[0].Marked {
		t.Fatalf("markItem() should unmark entry")
	}
}

func TestSelectorCommonResult(t *testing.T) {
	sel := &SelectorCommon{}
	if res, code := sel.result(); res != nil || code != ResultCancel {
		t.Fatalf("result() = (%v,%v), want (nil,%v)", res, code, ResultCancel)
	}
	want := []string{"/tmp/file.txt"}
	sel.marked = want
	res, code := sel.result()
	if code != ResultOK || len(res) != 1 || res[0] != want[0] {
		t.Fatalf("result() = (%v,%v), want (%v,%v)", res, code, want, ResultOK)
	}
}

func TestSelectorCommonGetFilterString(t *testing.T) {
	filter := filenameFilterSingleton()
	prev := filter.filterString
	prevIgnore := filter.ignoreCase
	defer func() {
		filter.filterString = prev
		filter.ignoreCase = prevIgnore
	}()

	filter.filterString = "abc"
	sel := &SelectorCommon{}
	if got := sel.getFilterString(); got != "abc" {
		t.Fatalf("getFilterString() = %q, want %q", got, "abc")
	}
}
