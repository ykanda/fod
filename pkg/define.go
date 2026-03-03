package fod

import "fmt"

// Mode select mode
type Mode int

// mode definition
const (
	ModeInvalid Mode = iota
	ModeFile
	ModeDirectory
)

// Entry type character
const (
	FsTypeFile    = "f"
	FsTypeDir     = "d"
	FsTypeSymlink = "s"
)

// Mode to string
func (mode Mode) String() (str string) {
	switch mode {
	case ModeFile:
		str = "file"
	case ModeDirectory:
		str = "directory"
	default:
		str = "invalid"
	}
	return
}

// StringToMode convert string to Mode value
func StringToMode(s string) (Mode, error) {
	switch s {
	case "file", "f":
		return ModeFile, nil
	case "directory", "dir", "d":
		return ModeDirectory, nil
	}
	return ModeInvalid, fmt.Errorf("unexpected mode string: %s", s)
}

// ResultCode dialog result code
type ResultCode int

// dialog result
const (
	ResultNone ResultCode = iota
	ResultCancel
	ResultOK
)
