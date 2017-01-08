package fod

// Mode select mode
type Mode int

// mode definition
const (
	ModeFile Mode = iota
	ModeDirectory
)

// Mode to string
func (mode Mode) String() (str string) {
	switch mode {
	case ModeFile:
		str = "file"
	case ModeDirectory:
		str = "directory"
	default:
		str = "unknown"
	}
	return
}

// ResultCode dialog result code
type ResultCode int

// dialog result
const (
	ResultNone ResultCode = iota
	ResultCancel
	ResultOK
)
