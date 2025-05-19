package net

import "net/http"

func GetFlyClientIP(r *http.Request) string {
	if clientIP := r.Header.Get("Fly-Client-IP"); clientIP != "" {
		return clientIP
	}

	// If not, return the remote address
	return "fly.io reverse proxy at: " + r.RemoteAddr
}
