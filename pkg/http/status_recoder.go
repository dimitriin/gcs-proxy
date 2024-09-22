package http

import (
	"net/http"
)

type StatusRecorder struct {
	http.ResponseWriter

	statusCode int
}

func NewStatusRecorder(responseWriter http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{ResponseWriter: responseWriter, statusCode: 0}
}

func (rec *StatusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *StatusRecorder) GetStatusCode() int {
	return rec.statusCode
}
