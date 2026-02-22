package fod

import "testing"

func TestDirectoryFilterSingleton(t *testing.T) {
	if directoryFilterSingleton() != directoryFilterSingleton() {
		t.Fatalf("directoryFilterSingleton() should return the same instance")
	}
}

func TestDirectoryFilterFilter(t *testing.T) {
	entries := []*Entry{
		{Path: "../", Type: FsTypeDir},
		{Path: "/tmp/dir", Type: FsTypeDir},
		{Path: "/tmp/file.txt", Type: FsTypeFile},
	}
	filter := &DirectoryFilter{}
	got := filter.filter(entries)
	if len(got) != 2 {
		t.Fatalf("DirectoryFilter.filter() len = %d, want 2", len(got))
	}
	for _, entry := range got {
		if entry.Type != FsTypeDir {
			t.Fatalf("DirectoryFilter.filter() returned non-dir entry: %+v", entry)
		}
	}
}
