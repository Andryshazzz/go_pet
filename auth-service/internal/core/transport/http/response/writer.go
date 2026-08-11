package httpresponse

import "net/http"

// StatusCodeUninitialized indicates that no status code has been set yet.
// Used as the initial value for ResponseWriter.statusCode.
var (
	StatusCodeUninitialized = -1
)

// ResponseWriter wraps http.ResponseWriter to track the written status code.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates a new ResponseWriter wrapping the given writer.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

// WriteHeader writes the HTTP status code and records it for later retrieval.
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

// GetStatusCodeOrPanic returns the status code that was written.
func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("No status code set")
	}
	return rw.statusCode
}
