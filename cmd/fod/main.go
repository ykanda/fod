package main

// [todo] - symlink
// [todo] - support windows
// [todo] - bookmark

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/nsf/termbox-go"
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
	ExitCodeOK    int = 0
	ExitCodeError     = 1
)

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

	// extends cli.Context
	var appContext = &fod.AppContext{context}

	// get working directory
	base := appContext.String("base")

	// create selector
	var selector *fod.SelectorFramework
	if _selector, err := fod.NewSelectorFramework(appContext.Mode(), appContext.Multi()); err == nil {
		selector = _selector
	} else {
		return err
	}

	if err := termbox.Init(); err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputEsc)

	// select loop
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		selector.Select(base)
		wg.Done()
		return
	}()
	wg.Wait()
	termbox.Close()

	if result, resultCode := selector.Result(); resultCode == fod.RESULT_OK {
		fmt.Fprintln(os.Stdout, result)
	}
	return nil
}

func flags() ([]cli.Flag, error) {

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Flags : options for urfave/cli
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "directory, d",
			Usage: "directory selecting mode",
		},
		cli.BoolFlag{
			Name:  "file, f",
			Usage: "file selecting mode",
		},
		cli.StringFlag{
			Name:  "base, b",
			Value: dir,
			Usage: "base dir",
		},
		cli.BoolFlag{
			Name:  "multi, m",
			Usage: "multiple select mode",
		},
	}
	return flags, nil
}
