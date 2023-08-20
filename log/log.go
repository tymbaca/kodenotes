package log

import (
	"errors"
	"io"
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	fatal *log.Logger
}

var (
	logger     = &Logger{}
)

func NewLogger(out io.Writer, infoPref, warnPref, errorPref, fatalPref string) *Logger {
	logger := &Logger{}
	logger.info = log.New(out, infoPref, 0)
	logger.warn = log.New(out, warnPref, 0)
	logger.error = log.New(out, errorPref, 0)
	logger.fatal = log.New(out, fatalPref, 0)
	return logger
}

func init() {
	path := "server.log"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		_, err2 := os.Create(path)
		if err2 != nil {
			panic(err2)
		}
	}
	logFile, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	multiOut := io.MultiWriter(os.Stdout, logFile)
	logger = NewLogger(multiOut, "[INFO]\t", "[WARN]\t", "[ERROR]\t", "[FATAL]\t")
}

func Info(msg string, args ...any) {
	logger.info.Printf(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.warn.Printf(msg, args...)
}

func Error(msg string, args ...any) {
	logger.error.Printf(msg, args...)
}

func Fatal(v ...any) {
	logger.error.Fatal(v...)
}
