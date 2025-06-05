package tokenauthority

import (
	"errors"
	"fmt"
	"time"
)

type TokenAuthority struct {
	jwtSecret     string
	jwtIssuer     string
	jwtLifetime   time.Duration
	rtSizeInBytes int
	rtLifetime    time.Duration
}

// New creates a new TokenAuthority instance with the provided configuration.
// It validates the configuration and returns a TokenAuthority instance or a joined error containing any validation errors.
//
// Parameters:
//   - jwtSecret: should be at least 16 characters long, ideally >= 32.
//   - jwtIssuer: The issuer of the JWT tokens.
//   - jwtLifetime: must be a positive value.
//   - rtSizeInBytes: must be at least 16, ideally >= 32.
//   - rtLifetime:  must be a positive value.
//
// Returns:
//   - *TokenAuthority: A pointer to the created TokenAuthority instance.
//   - error: An error containing all the errors that occurred during validation, if any.
func New(
	jwtSecret string,
	jwtIssuer string,
	jwtLifetime time.Duration,
	rtSizeInBytes int,
	rtLifetime time.Duration,
) (*TokenAuthority, error) {

	t := &TokenAuthority{
		jwtSecret:     jwtSecret,
		jwtIssuer:     jwtIssuer,
		jwtLifetime:   jwtLifetime,
		rtSizeInBytes: rtSizeInBytes,
		rtLifetime:    rtLifetime,
	}

	if err := t.validateConfig(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TokenAuthority) validateConfig() error {
	var errs []error

	if len(t.jwtSecret) < 16 {
		errstr := "jwtSecret must be at least 16 characters long, ideally >= 32"
		errs = append(errs, errors.New(errstr))
	}
	if t.jwtLifetime <= 0 {
		errs = append(errs, errors.New("jwtLifetime must be a positive value"))
	}
	if t.rtSizeInBytes < 16 {
		errs = append(errs, errors.New("rtSizeInBytes must be at least 16, ideally >= 32"))
	}
	if t.rtLifetime <= 0 {
		errs = append(errs, errors.New("rtLifetime must be a positive value"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("TokenAuthority configuration errors: %w", errors.Join(errs...))
	}

	return nil
}

func (t *TokenAuthority) GetRefreshTokenSizeInBytes() int {
	return t.rtSizeInBytes
}
