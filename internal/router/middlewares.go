package router

import (
	"log"
	"net/http"
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
		next(lrw, r)
		w.Header().Get("")
		log.Printf("Request handled: | %s\t| %s\t|X| %d", r.URL, r.Method, lrw.StatusCode)
	}
}
