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
)

import (
	"github.com/codegangsta/cli"
	"github.com/mitchellh/panicwrap"
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
		Name:  "multi, m",
		Usage: "multiple select mode",
	},
}

// entry point
func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	var exitStatus int
	if os.Getenv("FOD_ENABLE_PANIC_LOG") != "" {
		_exitStatus, _ := panicwrap.BasicWrap(panicHandler)
		exitStatus = _exitStatus
	}

	InitDebug()
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

	if os.Getenv("FOD_ENABLE_PANIC_LOG") != "" {
		if exitStatus > 0 {
			os.Exit(exitStatus)
		}
	}
}

func panicHandler(output string) {
	f, _ := os.Create(fmt.Sprintf("crash_%d.log", time.Now().Unix()))
	fmt.Fprintf(f, "The child panicked!\n\n%s", output)
	os.Exit(1)
}

// main function
func doMain(context *cli.Context) {

	// extends cli.Context
	var appContext *AppContext = &AppContext{context}

	// get working directory
	base := appContext.String("base")
	logger.Printf("%#v\n", base)

	// create selector
	var selector *SelectorFramework
	if _selector, err := NewSelectorFramework(appContext.Mode(), appContext.Multi()); err == nil {
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
