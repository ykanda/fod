package main

import (
	"io/ioutil"
	"os"
	"runtime"
)

import "github.com/k0kubun/pp"

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
		ioutil.WriteFile("./fod.log", []byte(output), os.ModePerm)
	}
}
