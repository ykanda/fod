package fod

// file selector
type FileSelector struct {
	*SelectorCommon
}

func (self *FileSelector) GetMode() string {
	return "F"
}

func (self *FileSelector) Mark() {
	if self.CurrentItemType(FsTypeFile) {
		self.MarkItem()
	}
}
