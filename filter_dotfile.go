package fod

import (
	"strings"
)

// dot file fileter (exclude dotfiles)
type DotfileFilter struct {
	enable bool
}

// singleton instance
var dotfileFilter *DotfileFilter = &DotfileFilter{
	enable: true,
}

// get singleton instance
func DotfileFilterSinleton() *DotfileFilter {
	return dotfileFilter
}

// filter function
func (self *DotfileFilter) Filter(entries []*Entry) (result []*Entry) {
	for _, entry := range entries {
		f := true
		if entry.Path != "../" {
			if self.enable == true && strings.HasPrefix(entry.Path, ".") {
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
func (self *DotfileFilter) Toggle() {
	self.enable = !self.enable
}
