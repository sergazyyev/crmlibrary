package ocrmserver

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rs *responseWriter) WriteHeader(statusCode int) {
	rs.statusCode = statusCode
	rs.ResponseWriter.WriteHeader(statusCode)
}
