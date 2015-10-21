package main

// [todo] - symlink
// [todo] - support windows
// [todo] - multiple select
// [todo] - bookmark

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

import (
	"github.com/codegangsta/cli"
	"github.com/nsf/termbox-go"
)

// flags
var Flags = []cli.Flag{
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
		Value: "./",
		Usage: "base dir",
	},
	cli.BoolFlag{
		Name:  "multiple, multi, m",
		Usage: "multiple select mode",
	},
}

// entry point
func main() {
	InitDebug()

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	app := cli.NewApp()
	app.Name = "vcd"
	app.Version = Version()
	app.Usage = "interactive file/directory selector"
	app.Author = "Yasuhiro KANDA"
	app.Email = "yasuhiro.kanda@gmail.com"
	app.Action = doMain
	app.Flags = Flags
	app.Run(os.Args)

	CloseDebug()
}

// main function
func doMain(context *cli.Context) {

	// extends cli.Context
	var appContext *AppContext = &AppContext{context}

	// get working directory
	base := appContext.String("base")
	DebugLog(base)

	// create selector
	var selector *SelectorFramework
	if _selector, err := NewSelectorFramework(appContext.Mode()); err == nil {
		selector = _selector
	} else {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if err := termbox.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		termbox.SetInputMode(termbox.InputEsc)
	}

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

	if result, resultCode := selector.Result(); resultCode == RESULT_OK {
		fmt.Fprintln(os.Stdout, result)
	}
}
