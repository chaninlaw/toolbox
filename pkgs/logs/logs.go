package logs

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Debug *log.Logger
)

func init() {
	if err := os.MkdirAll("log", os.ModePerm); err != nil {
		log.Printf("error during create log directory")
		os.Exit(1)
	}

	logFile, err := os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Printf("error during open app.log file")
		os.Exit(1)
	}

	infoWriter := io.MultiWriter(os.Stdout, logFile)
	warnWriter := io.MultiWriter(os.Stdout, logFile)
	errorWriter := io.MultiWriter(os.Stderr, logFile)
	debugWriter := io.MultiWriter(os.Stderr, logFile)

	flags := log.Ldate | log.Ltime | log.Lshortfile

	Info = log.New(infoWriter, "[INFO] ", flags)
	Warn = log.New(warnWriter, "[WARN] ", flags)
	Error = log.New(errorWriter, "[ERROR] ", flags)
	Debug = log.New(debugWriter, "[DEBUG] ", flags)
}
