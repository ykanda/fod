package fod

import "github.com/urfave/cli"

// extends cli.Context
type AppContext struct {
	*cli.Context
}

// get mode
func (self *AppContext) Mode() (mode Mode) {
	switch {
	case self.Bool("file"):
		mode = MODE_FILE
	case self.Bool("directory"):
		mode = MODE_DIRECTORY
	default:
		mode = MODE_DIRECTORY
	}
	// DebugPrint("mode", mode)
	return
}

func (self *AppContext) Multi() bool {
	return self.Bool("multi")
}
