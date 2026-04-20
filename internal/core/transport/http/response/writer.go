package core_http_response

import "net/http"

var (
	StatusCodeUninitializationed = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitializationed,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitializationed {
		panic("no status code set")
	}
	return rw.statusCode
}
