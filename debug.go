package main

import (
	"log"
	"os"
	"runtime"
)

import "github.com/k0kubun/pp"

var logger *log.Logger = nil
var logfile *os.File = nil

func nullDevice() (*os.File, error) {
	switch runtime.GOOS {
	case "windows":
		return os.Create("nul")
	case "darwin":
		fallthrough
	case "linux":
		fallthrough
	case "freebsd":
		fallthrough
	default:
		return os.Create("/dev/null")
	}
}

func InitDebug() {
	if os.Getenv("FOD_ENABLE_DEBUG") != "" {
		scheme := pp.ColorScheme{
			Bool:            pp.NoColor,
			Integer:         pp.NoColor,
			Float:           pp.NoColor,
			String:          pp.NoColor,
			StringQuotation: pp.NoColor,
			EscapedChar:     pp.NoColor,
			FieldName:       pp.NoColor,
			PointerAdress:   pp.NoColor,
			Nil:             pp.NoColor,
			Time:            pp.NoColor,
			StructName:      pp.NoColor,
			ObjectLength:    pp.NoColor,
		}
		pp.SetColorScheme(scheme)

		// open log file
		if _logfile, _err := os.Create("./fod.log"); _err != nil {
			panic(_err)
		} else {
			logfile = _logfile
		}
	} else {
		// open log file
		if _logfile, _err := nullDevice(); _err != nil {
			panic(_err)
		} else {
			logfile = _logfile
		}
	}
	logger = log.New(logfile, "fod: ", log.Lshortfile)
	if logger == nil {
	}
}

func CloseDebug() {
	if logfile != nil {
		logfile.Close()
		logfile = nil
	}
}
