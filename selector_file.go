package main

// file selector
type FileSelector struct {
	*SelectorCommon
}

func (self *FileSelector) GetMode() string {
	return "F"
}

func (self *FileSelector) Mark() {
	if self.CurrentItemType(FS_TYPE_FILE) {
		self.MarkItem()
	}
}
