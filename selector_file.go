package fod

// FileSelector selector for file
type FileSelector struct {
	*SelectorCommon
}

func (selector *FileSelector) getMode() string {
	return "F"
}

func (selector *FileSelector) markItem() {
	if selector.currentItemType(FsTypeFile) {
		selector.SelectorCommon.markItem()
	}
}
