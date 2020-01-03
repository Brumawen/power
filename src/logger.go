package main

import (
	"net/http"
	"time"
)

// Logger will create a Logger Handler wrapper for the specified handler.
func Logger(c Controller, inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		c.LogInfo(r.Method, r.RequestURI, "from", r.RemoteAddr, "tool", time.Since(start))
	})
}
