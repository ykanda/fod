package main

// file selector
type FileSelector struct {
	*SelectorCommon
}

func (self *FileSelector) GetMode() string {
	return "F"
}

// set result
func (self *FileSelector) Choice() (chose bool) {
	if self.CurrentItemType(FS_TYPE_DIR) {
		chose = false
	} else if path, err := self.CurrentItem(); err == nil {
		self.result = path
		self.resultCode = RESULT_OK
		chose = true
	}
	return
}
