package fod

import "testing"

func TestNewSelectorInvalidMode(t *testing.T) {
	if _, err := newSelector(ModeInvalid, false); err == nil {
		t.Fatalf("newSelector(ModeInvalid) expected error")
	}
}

func TestNewSelectorFile(t *testing.T) {
	sel, err := newSelector(ModeFile, true)
	if err != nil {
		t.Fatalf("newSelector(ModeFile) error: %v", err)
	}
	if _, ok := sel.(*FileSelector); !ok {
		t.Fatalf("newSelector(ModeFile) = %T, want *FileSelector", sel)
	}
}

func TestNewSelectorDirectory(t *testing.T) {
	sel, err := newSelector(ModeDirectory, false)
	if err != nil {
		t.Fatalf("newSelector(ModeDirectory) error: %v", err)
	}
	if _, ok := sel.(*DirectorySelector); !ok {
		t.Fatalf("newSelector(ModeDirectory) = %T, want *DirectorySelector", sel)
	}
}
