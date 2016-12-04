package main

// [todo] - symlink
// [todo] - support windows
// [todo] - bookmark

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/panicwrap"
	"github.com/nsf/termbox-go"
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
		Name:  "multi, m",
		Usage: "multiple select mode",
	},
}

// entry point
func main() {

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	os.Exit(run(os.Args))
}

func run(args []string) int {

	var exitStatus int
	if os.Getenv("FOD_ENABLE_PANIC_LOG") != "" {
		_exitStatus, _ := panicwrap.BasicWrap(panicHandler)
		exitStatus = _exitStatus
	}

	fod.InitDebug()
	app := cli.NewApp()
	app.Name = name
	app.Version = versionStr()
	app.Usage = "interactive file/directory selector"
	app.Author = "Yasuhiro KANDA"
	app.Email = "yasuhiro.kanda@gmail.com"
	app.Action = doMain
	app.Flags = Flags
	app.Run(args)
	fod.CloseDebug()

	if os.Getenv("FOD_ENABLE_PANIC_LOG") != "" {
		if exitStatus > 0 {
			os.Exit(exitStatus)
		}
	}
	return ExitCodeOK
}

func panicHandler(output string) {
	f, _ := os.Create(fmt.Sprintf("crash_%d.log", time.Now().Unix()))
	fmt.Fprintf(f, "The child panicked!\n\n%s", output)
	os.Exit(1)
}

// main function
func doMain(context *cli.Context) {

	// extends cli.Context
	var appContext *fod.AppContext = &fod.AppContext{context}

	// get working directory
	base := appContext.String("base")

	// create selector
	var selector *fod.SelectorFramework
	if _selector, err := fod.NewSelectorFramework(appContext.Mode(), appContext.Multi()); err == nil {
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

	if result, resultCode := selector.Result(); resultCode == fod.RESULT_OK {
		fmt.Fprintln(os.Stdout, result)
	}
}
