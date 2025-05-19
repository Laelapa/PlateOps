package middleware

import (
	"net/http"
	"strings"
)

func SecurityResponseHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// TODO: Consider adding a CSP header
		next.ServeHTTP(w, r)
	})
}

func CacheControlHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			// Cache static assets for 1 week
			w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		}
		next.ServeHTTP(w, r)
	})
}
