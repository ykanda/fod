package main

// function mode
type Mode int

// mode definition
const (
	MODE_FILE Mode = iota
	MODE_DIRECTORY
)

// Mode to string
func (mode Mode) String() (str string) {
	switch mode {
	case MODE_FILE:
		str = "file"
	case MODE_DIRECTORY:
		str = "directory"
	default:
		str = "unknown"
	}
	return
}

// dialog result
type ResultCode int

// dialog result
const (
	RESULT_NONE ResultCode = iota
	RESULT_CANCEL
	RESULT_OK
)
