package net

import (
	"errors"
	"net"
	"net/http"
)

// TODO: Enough for a controlled environment but not bulletproof.
// GetFlyClientIP retrieves the client's IP address from the request.
// It is designed to work with the fly.io reverse proxy only and is UNSAFE if not.
// It first checks for the "Fly-Client-IP" header, which is specific to fly.io.
// If the header is not present, it falls back to the remote address of the request.
// If the header is present but not a valid IP, it returns an empty string.
func GetFlyClientIP(r *http.Request) string {
	// Check for fly.io-specific header.
	if clientIP := r.Header.Get("Fly-Client-IP"); clientIP != "" {
		if ip := net.ParseIP(clientIP); ip != nil {
			return ip.String()
		}

		// If the header is present but not a valid IP, return an empty string.
		return ""
	}

	// If nonexistent, return the remote address
	return StripPort(r.RemoteAddr)
}

func StripPort(ipAddress string) string {
	host, _, err := net.SplitHostPort(ipAddress)
	if err != nil {
		return ipAddress
	}
	return host
}

func StripBearer(authHeader string) (string, error) {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}

	return "", errors.New("invalid authorization header format")
}
