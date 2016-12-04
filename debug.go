package fod

import (
	"errors"
	"log"
	"os"
	"runtime"
)

var logger *log.Logger
var logfile *os.File

// InitDebug : initialize debug function
func InitDebug() error {

	logfile, err := createLogFile()
	if err != nil {
		return err
	}

	logger = log.New(logfile, "fod: ", log.Lshortfile)
	if logger == nil {
		return errors.New("failed to create debug log file")
	}

	return nil
}

// CloseDebug : finalize debug function
func CloseDebug() {
	if logfile != nil {
		logfile.Close()
		logfile = nil
	}
}

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

func createLogFile() (*os.File, error) {
	if os.Getenv("FOD_ENABLE_DEBUG") == "" {
		return nullDevice()
	}
	return os.Create("./fod.log")
}
