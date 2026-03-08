package fod

import (
	"reflect"
	"testing"
)

func TestFindMatchRanges(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		text       string
		words      []string
		ignoreCase bool
		want       [][2]int
	}{
		{
			name:  "single match",
			text:  "/tmp/alpha.txt",
			words: []string{"alpha"},
			want:  [][2]int{{5, 10}},
		},
		{
			name:  "multiple matches same word",
			text:  "foo_bar_foo",
			words: []string{"foo"},
			want:  [][2]int{{0, 3}, {8, 11}},
		},
		{
			name:  "multiple words merge overlap",
			text:  "foobar",
			words: []string{"foo", "oob"},
			want:  [][2]int{{0, 4}},
		},
		{
			name:       "ignore case",
			text:       "HelloWorld",
			words:      []string{"hello", "WORLD"},
			ignoreCase: true,
			want:       [][2]int{{0, 10}},
		},
		{
			name:  "no match",
			text:  "abc",
			words: []string{"zzz"},
			want:  nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := findMatchRanges(tc.text, tc.words, tc.ignoreCase)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("findMatchRanges() = %#v, want %#v", got, tc.want)
			}
		})
	}
}
