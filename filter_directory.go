package main

// directory filter
type DirectoryFilter struct {
}

// singleton instance
var directoryFilter *DirectoryFilter = &DirectoryFilter{}

// get singleton instance
func DirectoryFilterSingleton() *DirectoryFilter {
	return directoryFilter
}

// filter function
func (self *DirectoryFilter) Filter(entries []*Entry) []*Entry {
	temp := []*Entry{}
	for _, entry := range entries {
		if entry.Type == FS_TYPE_DIR {
			temp = append(temp, entry)
		}
	}
	return temp
}
