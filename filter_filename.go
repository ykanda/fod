package fod

import (
	"strings"
)

// FilenameFilter : file name filter
type FilenameFilter struct {
	filterString string
	ignoreCase   bool
}

// singleton instance
var filenameFilter = &FilenameFilter{
	filterString: "",
	ignoreCase:   false,
}

func (filter *FilenameFilter) addCharacter(char rune) {
	filter.filterString += string(char)
}

func (filter *FilenameFilter) delCharacter() {
	s := filter.filterString
	if len(s) > 0 {
		filter.filterString = s[:len(s)-1]
	}
}

// get singleton instance
func filenameFilterSingleton() *FilenameFilter {
	return filenameFilter
}

// filter function
func (filter *FilenameFilter) filter(entries []*Entry) (result []*Entry) {
	switch filter.filterString == "" {
	case true:
		result = entries

	case false:
		for _, entry := range entries {
			for _, word := range strings.Split(filter.filterString, " ") {
				path := entry.Path
				if path == "../" {
					result = append(result, entry)
					continue
				}
				if filter.ignoreCase {
					path = strings.ToLower(path)
					word = strings.ToLower(word)
				}
				if strings.Contains(path, word) {
					result = append(result, entry)
				}
			}
		}
	}
	return
}

func (filter *FilenameFilter) setFilterString(filterString string) {
	filter.filterString = filterString
}

func (filter *FilenameFilter) getFilterString() string {
	return filter.filterString
}
