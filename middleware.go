package main

import (
	"net/http"
)

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			method = r.Method
			url    = r.URL.String()
		)

		logger.Debug("Request %s %s", method, url)

		next.ServeHTTP(w, r)
	})
}
