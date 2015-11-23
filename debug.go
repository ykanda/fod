package main

import (
	"io/ioutil"
	"os"
	"runtime"
)

import "github.com/k0kubun/pp"

func InitDebug() {
	if os.Getenv("FOD_ENABLE_DEBUG") != "" {
	}
}

func CloseDebug() {
	if os.Getenv("FOD_ENABLE_DEBUG") != "" {
	}
}

func DebugLog(args ...interface{}) {
	if os.Getenv("FOD_ENABLE_DEBUG") != "" {
		_, file, line, _ := runtime.Caller(1)
		output := pp.Sprintf(
			"%s line %s\n%v\n",
			file,
			line,
			args,
		)
		ioutil.WriteFile("/tmp/vcd.log", []byte(output), os.ModePerm)
	}
}
