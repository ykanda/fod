package fod

import (
	"strings"
)

// file name filter
type FilenameFilter struct {
	filterString string
	ignoreCase   bool
}

// singleton instance
var filenameFilter *FilenameFilter = &FilenameFilter{
	filterString: "",
	ignoreCase:   false,
}

func (self *FilenameFilter) AddCharacter(char rune) {
	self.filterString += string(char)
}

func (self *FilenameFilter) DelCharacter() {
	s := self.filterString
	if len(s) > 0 {
		self.filterString = s[:len(s)-1]
	}
}

// get singleton instance
func FilenameFilterSingleton() *FilenameFilter {
	return filenameFilter
}

// filter function
func (self *FilenameFilter) Filter(entries []*Entry) (result []*Entry) {
	switch self.filterString == "" {
	case true:
		result = entries

	case false:
		for _, entry := range entries {
			for _, word := range strings.Split(self.filterString, " ") {
				path := entry.Path
				if path == "../" {
					result = append(result, entry)
					continue
				}
				if self.ignoreCase {
					path = strings.ToLower(path)
					word = strings.ToLower(word)
				}
				if strings.Index(path, word) != -1 {
					result = append(result, entry)
				}
			}
		}
	}
	return
}

func (self *FilenameFilter) SetFilterString(filterString string) {
	self.filterString = filterString
}

func (self *FilenameFilter) GetFilterString() string {
	return self.filterString
}
