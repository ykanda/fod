package fod

// DirectorySelector : directory selector
type DirectorySelector struct {
	*SelectorCommon
}

func (selector *DirectorySelector) getMode() string {
	return "D"
}

func (selector *DirectorySelector) mark() {
	if selector.currentItemType(FsTypeDir) {
		selector.markItem()
	}
}
