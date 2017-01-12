package fod

// DirectorySelector : directory selector
type DirectorySelector struct {
	*SelectorCommon
}

func (selector *DirectorySelector) getMode() string {
	return "D"
}

func (selector *DirectorySelector) markItem() {
	if selector.currentItemType(FsTypeDir) {
		selector.SelectorCommon.markItem()
	}
}
