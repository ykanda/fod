package fod

import (
	"reflect"
	"testing"
)

func TestFilenameFilterSingleton(t *testing.T) {
	if filenameFilterSingleton() != filenameFilterSingleton() {
		t.Fatalf("filenameFilterSingleton() should return the same instance")
	}
}

func TestFilenameFilterFilter(t *testing.T) {
	entries := []*Entry{
		{Path: "../", Type: FsTypeDir},
		{Path: "/tmp/alpha.txt", Type: FsTypeFile},
		{Path: "/tmp/BRAVO.txt", Type: FsTypeFile},
		{Path: "/tmp/charlie.txt", Type: FsTypeFile},
	}

	filter := &FilenameFilter{filterString: ""}
	got := filter.filter(entries)
	if !reflect.DeepEqual(got, entries) {
		t.Fatalf("FilenameFilter.filter() empty = %#v, want %#v", got, entries)
	}

	filter = &FilenameFilter{filterString: "alpha BRAVO", ignoreCase: true}
	got = filter.filter(entries)
	if len(got) != 4 {
		t.Fatalf("FilenameFilter.filter() ignoreCase len = %d, want 4", len(got))
	}
	if got[0].Path != "../" {
		t.Fatalf("FilenameFilter.filter() should keep ../ entry, got %q", got[0].Path)
	}
}

func TestFilenameFilterMutators(t *testing.T) {
	filter := &FilenameFilter{}
	filter.addCharacter('a')
	filter.addCharacter('b')
	if filter.getFilterString() != "ab" {
		t.Fatalf("getFilterString() = %q, want %q", filter.getFilterString(), "ab")
	}
	filter.delCharacter()
	if filter.getFilterString() != "a" {
		t.Fatalf("getFilterString() after del = %q, want %q", filter.getFilterString(), "a")
	}
	filter.delCharacter()
	filter.delCharacter()
	if filter.getFilterString() != "" {
		t.Fatalf("getFilterString() after extra del = %q, want empty", filter.getFilterString())
	}
}
