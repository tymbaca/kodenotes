package log

import (
	"io"
	"log"
	"net/http"
	"os"
)

type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	fatal *log.Logger
}

var (
	logger = &Logger{}
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
	logger = NewLogger(os.Stdout, "[INFO]\t", "[WARN]\t", "[ERROR]\t", "[FATAL]\t")
}

func SetOutput(out io.Writer) {
	logger = NewLogger(out, "[INFO]\t", "[WARN]\t", "[ERROR]\t", "[FATAL]\t")
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

func RequestInfo(r *http.Request, msg string, respStatus int) {
	Info("| %s\t| %s\t| %s\t|X| %s\t| %d", r.Host, r.URL, r.Method, msg, respStatus)
}
