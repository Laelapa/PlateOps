package env

import (
	"strconv"
	"time"

	"github.com/Laelapa/GoHome/logging"
	"go.uber.org/zap"
)

const defaultPort = "8080"
const defaultJwtLifetime = 3 * time.Hour     // 3 hours
const defaultRtLifetime = 7 * 24 * time.Hour // 7 days

// ValidatePort checks if the provided port string is a valid port number.
// For it to be valid, it must be a number between 1 and 65535.
// If it is valid, it returns the port as a string.
// If it is invalid, it logs an error and returns the default port 8080.
func ValidatePort(port string, logger *logging.Logger) string {

	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		logger.LogAppError(
			"Invalid port number provided, using default port",
			err,
			zap.String("provided_port", port),
			zap.String("default_port", defaultPort),
		)

		return defaultPort
	}

	return port
}

func ParseLifetimeJWT(jwtLifetime string, logger *logging.Logger) time.Duration {
	lifetime, err := time.ParseDuration(jwtLifetime)
	if err != nil || lifetime <= 0 {
		logger.LogAppError(
			"Invalid JWT lifetime provided, using default lifetime",
			err,
			zap.String("provided_lifetime", jwtLifetime),
			zap.String("default_lifetime", defaultJwtLifetime.String()),
		)

		return defaultJwtLifetime
	}

	return lifetime
}

func ParseLifetimeRT(rtLifetime string, logger *logging.Logger) time.Duration {
	lifetime, err := time.ParseDuration(rtLifetime)
	if err != nil || lifetime <= 0 {
		logger.LogAppError(
			"Invalid refresh token lifetime provided, using default lifetime",
			err,
			zap.String("provided_lifetime", rtLifetime),
			zap.String("default_lifetime", defaultRtLifetime.String()),
		)

		return defaultRtLifetime
	}

	return lifetime
}
