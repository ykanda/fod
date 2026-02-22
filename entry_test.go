package fod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEntries(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fod-entries-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	filePath := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got := entries(tmpDir)
	if len(got) != 3 {
		t.Fatalf("entries(%q) len = %d, want 3", tmpDir, len(got))
	}

	if got[0].Path != "../" || got[0].Type != FsTypeDir || got[0].Marked {
		t.Fatalf("entries(%q)[0] = %+v, want Path=../ Type=%s Marked=false", tmpDir, got[0], FsTypeDir)
	}

	absSub, err := filepath.Abs(subDir)
	if err != nil {
		t.Fatalf("Abs(subDir): %v", err)
	}
	absFile, err := filepath.Abs(filePath)
	if err != nil {
		t.Fatalf("Abs(filePath): %v", err)
	}

	byPath := map[string]*Entry{}
	for _, e := range got {
		byPath[e.Path] = e
	}

	if e := byPath[absSub]; e == nil || e.Type != FsTypeDir || e.Marked {
		t.Fatalf("entry for %q = %+v, want Type=%s Marked=false", absSub, e, FsTypeDir)
	}
	if e := byPath[absFile]; e == nil || e.Type != FsTypeFile || e.Marked {
		t.Fatalf("entry for %q = %+v, want Type=%s Marked=false", absFile, e, FsTypeFile)
	}
}

func TestEntryTypeCharcter(t *testing.T) {
	cases := []struct {
		entry Entry
		want  string
	}{
		{Entry{Type: FsTypeDir}, "d"},
		{Entry{Type: FsTypeFile}, "-"},
		{Entry{Type: "unknown"}, ""},
	}
	for _, tc := range cases {
		if got := tc.entry.typeCharcter(); got != tc.want {
			t.Fatalf("Entry{Type:%q}.typeCharcter() = %q, want %q", tc.entry.Type, got, tc.want)
		}
	}
}
