package fod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectorySelectorGetMode(t *testing.T) {
	sel := &DirectorySelector{}
	if sel.getMode() != "Dir" {
		t.Fatalf("getMode() = %q, want %q", sel.getMode(), "Dir")
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

func TestDirectorySelectorSelectAll(t *testing.T) {
	common := &SelectorCommon{
		Multi: true,
		Entries: []*Entry{
			{Path: "../", Type: FsTypeDir},
			{Path: "/tmp/sub", Type: FsTypeDir},
			{Path: "/tmp/a.txt", Type: FsTypeFile},
		},
	}
	selector := &DirectorySelector{SelectorCommon: common}

	selector.selectAll()
	if len(selector.marked) != 1 {
		t.Fatalf("selectAll() marked size = %d, want 1", len(selector.marked))
	}
	if selector.Entries[0].Marked {
		t.Fatalf("selectAll() should not mark parent entry")
	}
	if !selector.Entries[1].Marked {
		t.Fatalf("selectAll() should mark directory entries in directory mode")
	}
	if selector.Entries[2].Marked {
		t.Fatalf("selectAll() should not mark file entries in directory mode")
	}
}
