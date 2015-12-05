package main

// directory selector
type DirectorySelector struct {
	*SelectorCommon
}

func (self *DirectorySelector) GetMode() string {
	return "D"
}

func (self *DirectorySelector) Mark() {
	if self.CurrentItemType(FS_TYPE_DIR) {
		self.MarkItem()
	}
}
