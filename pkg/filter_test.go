package fod

import "testing"

var _ Filter = (*DirectoryFilter)(nil)
var _ Filter = (*FileFilter)(nil)
var _ Filter = (*DotfileFilter)(nil)
var _ Filter = (*FilenameFilter)(nil)

func TestFilters(t *testing.T) {
	entries := []*Entry{
		{Path: "../", Type: FsTypeDir},
		{Path: "/tmp/dir", Type: FsTypeDir},
		{Path: "/tmp/file.txt", Type: FsTypeFile},
		{Path: ".hidden", Type: FsTypeFile},
	}

	cases := []struct {
		name string
		f    Filter
		want int
	}{
		{"file", &FileFilter{}, 4},
		{"directory", &DirectoryFilter{}, 2},
		{"dotfile", &DotfileFilter{enable: true}, 3},
		{"filename", &FilenameFilter{filterString: "file"}, 2},
	}

	for _, tc := range cases {
		got := tc.f.filter(entries)
		if len(got) != tc.want {
			t.Fatalf("%s.filter() len = %d, want %d", tc.name, len(got), tc.want)
		}
	}
}
