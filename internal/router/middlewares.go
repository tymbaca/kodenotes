package router

import (
	"log"
	"net/http"
	"time"
)

// LogResponseWriter is a wrapper for statndart http.ResponseWriter with
// additional StatusCode field. Overloads the http.ResponseWriter.WriteHeader
// in the way that when underlying handler uses that method, then writen status
// code is dublicated to StatusCode field. Later it can be used inside
// middlewares (e.g. logger middleware).
type LogResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// NewLogResponseWriter creates new instance of LogResponseWriter with w
// and http.StatusOK by default.
func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{
		w,
		http.StatusOK,
	}
}

// WriteHeader calls standart http.ResponseWriter.WriteHeader and also writes
// status code to lrw.StatusCode field.
func (lrw *LogResponseWriter) WriteHeader(statusCode int) {
	lrw.StatusCode = statusCode
	lrw.WriteHeader(statusCode)
}

// Logger is a Middleware that writes request URL, method and response status code to standart log
// You can implement your own logger by using NewLogResponseWriter and passing it to next
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLogResponseWriter(w)
		startTime := time.Now()

		next(lrw, r)

		endTime := time.Now()
		duration := endTime.Sub(startTime)
		log.Printf(
			"Request handled: | %s \t| %s \t|X| %d \t| %s",
			r.URL, r.Method, lrw.StatusCode, duration.String(),
		)
	}
}

// WARN: doesnt work
func Recover(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			p := recover()
			if p != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// Doesn't work currently
func Duration(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next(w, r)

		endTime := time.Now()
		duration := startTime.Sub(endTime)
		w.Header().Set("X-Duration", duration.String())
	}
}
