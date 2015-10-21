package main

// directory selector
type DirectorySelector struct {
	*SelectorCommon
}

func (self *DirectorySelector) GetMode() string {
	return "D"
}

// set result
func (self *DirectorySelector) Choice() (chose bool) {
	if self.CurrentItemType(FS_TYPE_FILE) {
		chose = false
	} else if path, err := self.CurrentItem(); err == nil {
		self.result = path
		self.resultCode = RESULT_OK
		chose = true
	}
	return
}
