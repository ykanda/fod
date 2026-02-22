package fod

import (
	"reflect"
	"testing"
)

func TestFileFilterSingleton(t *testing.T) {
	if fileFilterSingleton() != fileFilterSingleton() {
		t.Fatalf("fileFilterSingleton() should return the same instance")
	}
}

func TestFileFilterFilter(t *testing.T) {
	entries := []*Entry{
		{Path: "../", Type: FsTypeDir},
		{Path: "/tmp/file.txt", Type: FsTypeFile},
	}
	filter := &FileFilter{}
	got := filter.filter(entries)
	if !reflect.DeepEqual(got, entries) {
		t.Fatalf("FileFilter.filter() = %#v, want %#v", got, entries)
	}
}
