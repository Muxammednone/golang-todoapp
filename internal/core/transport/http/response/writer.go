package core_http_response

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

var (
	StatusCodeUninitialized = -1
)

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) GetStatusCode() int {
	return rw.statusCode
}

func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("no status code set")
	}

	return rw.statusCode
}
