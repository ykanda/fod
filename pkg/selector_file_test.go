package fod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileSelectorGetMode(t *testing.T) {
	sel := &FileSelector{}
	if sel.getMode() != "F" {
		t.Fatalf("getMode() = %q, want %q", sel.getMode(), "F")
	}
}

func TestFileSelectorMarkItem(t *testing.T) {
	root, err := os.MkdirTemp("", "fod-filesel-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	defer os.RemoveAll(root)

	sub := filepath.Join(root, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	filePath := filepath.Join(root, "file.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	common := &SelectorCommon{
		Multi: true,
		Entries: []*Entry{
			{Path: sub, Type: FsTypeDir},
			{Path: filePath, Type: FsTypeFile},
		},
	}
	selector := &FileSelector{SelectorCommon: common}

	selector.Cursor = 0
	selector.markItem()
	if len(selector.marked) != 0 {
		t.Fatalf("markItem() should not mark directory, marked=%#v", selector.marked)
	}

	selector.Cursor = 1
	selector.markItem()
	if len(selector.marked) != 1 || selector.marked[0] != filePath {
		t.Fatalf("markItem() marked = %#v, want %q", selector.marked, filePath)
	}
	if !selector.Entries[1].Marked {
		t.Fatalf("markItem() should mark file entry")
	}
}

func TestFileSelectorSelectAll(t *testing.T) {
	common := &SelectorCommon{
		Multi: true,
		Entries: []*Entry{
			{Path: "../", Type: FsTypeDir},
			{Path: "/tmp/sub", Type: FsTypeDir},
			{Path: "/tmp/a.txt", Type: FsTypeFile},
			{Path: "/tmp/b.txt", Type: FsTypeFile},
		},
	}
	selector := &FileSelector{SelectorCommon: common}

	selector.selectAll()
	if len(selector.marked) != 2 {
		t.Fatalf("selectAll() marked size = %d, want 2", len(selector.marked))
	}
	if selector.Entries[0].Marked || selector.Entries[1].Marked {
		t.Fatalf("selectAll() should not mark directory entries in file mode")
	}
	if !selector.Entries[2].Marked || !selector.Entries[3].Marked {
		t.Fatalf("selectAll() should mark only file entries in file mode")
	}
}
