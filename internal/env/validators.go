package env

import (
	"strconv"

	"go.uber.org/zap"

	"github.com/Laelapa/GoHome/logging"
)

const defaultPort = "8080"

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
