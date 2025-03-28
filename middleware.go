package main

import (
	"log/slog"
	"net/http"
)

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			method = r.Method
			url    = r.URL.String()
		)

		slog.Debug("Request", "method", method, "url", url)

		next.ServeHTTP(w, r)
	})
}
