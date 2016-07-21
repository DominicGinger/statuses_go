package main

import (
	"log"
	"os"
)

type loggerStruct struct {
	info *log.Logger
	warn *log.Logger
	err  *log.Logger
}

const enableLog = true

func createLogger(prefix string, out *os.File) *log.Logger {
	if !enableLog {
		out, _ = os.Create("")
	}
	return log.New(out, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}

var logger = loggerStruct{
	info: createLogger("[INFO] ", os.Stdout),
	warn: createLogger("[WARN] ", os.Stdout),
	err:  createLogger("[ERROR] ", os.Stderr),
}
