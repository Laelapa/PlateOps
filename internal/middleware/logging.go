package middleware

import (
	"net/http"

	"github.com/Laelapa/GoHome/logging"
)

func RequestLogger(next http.Handler, logger *logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.LogRequestInfo("HTTP Request", r)
		next.ServeHTTP(w, r)
	})
}
