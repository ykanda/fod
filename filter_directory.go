package fod

// DirectoryFilter : directory filter
type DirectoryFilter struct {
}

// singleton instance
var directoryFilter = &DirectoryFilter{}

// directoryFilterSingleton return singleton instance
func directoryFilterSingleton() *DirectoryFilter {
	return directoryFilter
}

// filter : filter function
func (selector *DirectoryFilter) filter(entries []*Entry) []*Entry {
	temp := []*Entry{}
	for _, entry := range entries {
		if entry.Type == FsTypeDir {
			temp = append(temp, entry)
		}
	}
	return temp
}
