package fod

import "github.com/urfave/cli"

// AppContext extends cli.Context
type AppContext struct {
	*cli.Context
}

// Mode get mode
func (ctx *AppContext) Mode() (mode Mode) {
	switch {
	case ctx.Bool("file"):
		mode = ModeFile
	case ctx.Bool("directory"):
		mode = ModeDirectory
	default:
		mode = ModeDirectory
	}
	// DebugPrint("mode", mode)
	return
}

// Multi return multiple selection mode
func (ctx *AppContext) Multi() bool {
	return ctx.Bool("multi")
}
