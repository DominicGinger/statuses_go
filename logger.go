package main

import (
	"log"
	"os"
)

type loggerStruct struct {
	debug *log.Logger
	err   *log.Logger
	info  *log.Logger
	warn  *log.Logger
}

const enableLog = true
const CLR_R = "\x1b[31;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"
const CLR_N = "\x1b[0m"

func createLogger(prefix string, out *os.File) *log.Logger {
	if !enableLog {
		out, _ = os.Create("")
	}
	return log.New(out, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}

var logger = loggerStruct{
	debug: createLogger(CLR_B+"[DEBUG] ", os.Stdout),
	err:   createLogger(CLR_R+"[ERROR] ", os.Stderr),
	info:  createLogger(CLR_N+"[INFO] ", os.Stdout),
	warn:  createLogger(CLR_Y+"[WARN] ", os.Stdout),
}
