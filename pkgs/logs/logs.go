package logs

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
)

func init() {
	infoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(msg string, args ...any) {
	infoLogger.Printf(msg, args...)
}

func Warn(msg string, args ...any) {
	warnLogger.Printf(msg, args...)
}

func Error(msg string, args ...any) {
	errorLogger.Printf(msg, args...)
}

func Fatal(msg string, args ...any) {
	fatalLogger.Printf(msg, args...)
	os.Exit(1)
}
