package main

// [todo] - symlink
// [todo] - support windows
// [todo] - bookmark

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/urfave/cli"
	"github.com/ykanda/fod"
)

var (
	name     string
	version  string
	revision string
)

func versionStr() string {
	return fmt.Sprintf(
		"%s version %s revision %s",
		name,
		version,
		revision,
	)
}

// ExitCode
const (
	ExitCodeOK    int = iota
	ExitCodeError     = 1
)

// DefaultPathSeparator default path separator
const DefaultPathSeparator = ":"

// entry point
func main() {

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	os.Exit(run(os.Args))
}

func run(args []string) int {

	fod.InitDebug()
	defer fod.CloseDebug()

	flags, err := flags()
	if err != nil {
		return ExitCodeError
	}
	app := cli.NewApp()
	app.Name = name
	app.Version = versionStr()
	app.Usage = "file open dialog"
	app.Author = "Yasuhiro KANDA"
	app.Email = "yasuhiro.kanda@gmail.com"
	app.Action = action
	app.Flags = flags
	if err := app.Run(args); err != nil {
		return ExitCodeError
	}

	return ExitCodeOK
}

// main function
func action(context *cli.Context) error {

	mode, err := fod.StringToMode(context.String("mode"))
	if err != nil {
		return err
	}
	output, result := fod.Dialog(
		fod.Option{
			Base:  context.String("base"),
			Multi: context.Bool("multi"),
			Mode:  mode,
		},
	)

	if result != fod.ResultOK {
		return errors.New("unexpected result code")
	}

	sep := context.String("separator")
	fmt.Fprintln(os.Stdout, strings.Join(output, sep))
	return nil
}

func flags() ([]cli.Flag, error) {

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Flags : options for urfave/cli
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "mode, m",
			Value: fod.ModeFile.String(),
		},
		cli.StringFlag{
			Name:  "base, b",
			Value: dir,
			Usage: "base dir",
		},
		cli.BoolFlag{
			Name:  "multi",
			Usage: "multiple select flag",
		},
		cli.StringFlag{
			Name:  "separator, s",
			Usage: "path separator character (string)",
			Value: DefaultPathSeparator,
		},
	}
	return flags, nil
}
