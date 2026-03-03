package fod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectorySelectorGetMode(t *testing.T) {
	sel := &DirectorySelector{}
	if sel.getMode() != "D" {
		t.Fatalf("getMode() = %q, want %q", sel.getMode(), "D")
	}
}

func TestDirectorySelectorMarkItem(t *testing.T) {
	root, err := os.MkdirTemp("", "fod-dirsel-*")
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
	selector := &DirectorySelector{SelectorCommon: common}

	selector.Cursor = 1
	selector.markItem()
	if len(selector.marked) != 0 {
		t.Fatalf("markItem() should not mark file, marked=%#v", selector.marked)
	}

	selector.Cursor = 0
	selector.markItem()
	if len(selector.marked) != 1 || selector.marked[0] != sub {
		t.Fatalf("markItem() marked = %#v, want %q", selector.marked, sub)
	}
	if !selector.Entries[0].Marked {
		t.Fatalf("markItem() should mark directory entry")
	}
}
