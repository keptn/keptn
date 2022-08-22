package middleware

import (
	"net/http"
)

// EnforceMaxEventSize returns a handler that uses a MaxBytesReader to limit the bytes that
// are read from the request body
func EnforceMaxEventSize(maxSizeBytes int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req.Body = http.MaxBytesReader(w, req.Body, maxSizeBytes)
			h.ServeHTTP(w, req)
		})
	}
}
