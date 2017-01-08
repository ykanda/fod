package fod

// FileSelector selector for file
type FileSelector struct {
	*SelectorCommon
}

func (selector *FileSelector) getMode() string {
	return "F"
}

func (selector *FileSelector) mark() {
	if selector.currentItemType(FsTypeFile) {
		selector.markItem()
	}
}
