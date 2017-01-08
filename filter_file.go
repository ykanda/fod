package fod

// FileFilter : file filter
type FileFilter struct {
}

// singleton instance
var fileFilter = &FileFilter{}

// fileFilterSingleton : get singleton instance
func fileFilterSingleton() *FileFilter {
	return fileFilter
}

// filter : filter function
func (filter *FileFilter) filter(entries []*Entry) []*Entry {
	return entries // NOP
}
