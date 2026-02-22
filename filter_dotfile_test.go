package fod

import "testing"

func TestDotfileFilterSingleton(t *testing.T) {
	if dotfileFilterSinleton() != dotfileFilterSinleton() {
		t.Fatalf("dotfileFilterSinleton() should return the same instance")
	}
}

func TestDotfileFilterFilter(t *testing.T) {
	entries := []*Entry{
		{Path: "../", Type: FsTypeDir},
		{Path: ".hidden", Type: FsTypeFile},
		{Path: "/tmp/file.txt", Type: FsTypeFile},
	}
	filter := &DotfileFilter{enable: true}
	got := filter.filter(entries)
	if len(got) != 2 {
		t.Fatalf("DotfileFilter.filter() len = %d, want 2", len(got))
	}
	for _, entry := range got {
		if entry.Path == ".hidden" {
			t.Fatalf("DotfileFilter.filter() unexpectedly included dotfile")
		}
	}
	filter.toggle()
	got = filter.filter(entries)
	if len(got) != 3 {
		t.Fatalf("DotfileFilter.filter() after toggle len = %d, want 3", len(got))
	}
}
