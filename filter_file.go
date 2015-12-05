package main

// file filter
type FileFilter struct {
}

// singleton instance
var fileFilter *FileFilter = &FileFilter{}

// get singleton instance
func FileFilterSingleton() *FileFilter {
	return fileFilter
}

// filter function
func (self *FileFilter) Filter(entries []*Entry) []*Entry {
	return entries // NOP
}
