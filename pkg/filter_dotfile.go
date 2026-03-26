package fod

import (
	"path/filepath"
	"strings"
)

// DotfileFilter : dot file fileter (exclude dotfiles)
type DotfileFilter struct {
	enable bool
}

// singleton instance
var dotfileFilter = &DotfileFilter{
	enable: true,
}

// get singleton instance
func dotfileFilterSinleton() *DotfileFilter {
	return dotfileFilter
}

// filter function
func (filter *DotfileFilter) filter(entries []*Entry) (result []*Entry) {
	for _, entry := range entries {
		f := true
		if entry.Path != "../" {
			name := filepath.Base(entry.Path)
			if filter.enable == true && strings.HasPrefix(name, ".") {
				f = false
			}
		}
		if f {
			result = append(result, entry)
		}
	}
	return
}

// toggle enable or disable
func (filter *DotfileFilter) toggle() {
	filter.enable = !filter.enable
}
