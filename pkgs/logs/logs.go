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

	stdoutWriter := io.MultiWriter(os.Stdout, logFile)
	stderrWriter := io.MultiWriter(os.Stderr, logFile)

	flags := log.Ldate | log.Ltime | log.Lshortfile

	Debug = log.New(stdoutWriter, "[DEBUG] ", flags)
	Info = log.New(stdoutWriter, "[INFO] ", flags)
	Warn = log.New(stdoutWriter, "[WARN] ", flags)
	Error = log.New(stderrWriter, "[ERROR] ", flags)
}
