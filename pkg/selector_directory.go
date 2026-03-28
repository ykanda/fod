package fod

// DirectorySelector : directory selector
type DirectorySelector struct {
	*SelectorCommon
}

func (selector *DirectorySelector) getMode() string {
	return "Dir"
}

func (selector *DirectorySelector) markItem() {
	if selector.currentItemType(FsTypeDir) {
		selector.SelectorCommon.markItem()
	}
}

func (selector *DirectorySelector) selectAll() {
	selector.SelectorCommon.selectAllByType(FsTypeDir)
}
