package tokenauthority

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// IssueJWT creates and signs a new JWT for the specified user ID.
// The token includes standard registered claims (issuer, subject, issued at, expires at)
// and is signed using HMAC SHA256 algorithm with the configured secret.
// The token's lifetime and the jwtSecret are properties of the TokenAuthority instance.
//
// Parameters:
//   - userID: The UUID of the user for whom the token is being issued
//
// Returns:
//   - string: The signed JWT as a string
//   - error: An error if token signing fails, nil otherwise
func (t *TokenAuthority) IssueJWT(userID uuid.UUID) (string, error) {

	claims := jwt.RegisteredClaims{
		Issuer:    t.jwtIssuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(t.jwtLifetime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(t.jwtSecret))
	if err != nil {
		return "", errors.New("failed to sign JWT" + err.Error())
	}

	return signedToken, nil
}

// ValidateJWT parses and validates a JWT token string using HMAC SHA256 signing method.
// It verifies the token's signature against the configured JWT secret and extracts the subject claim.
// Returns the subject from the token's claims if validation succeeds, or an error if the token
// is invalid, uses an unexpected signing method, or lacks a subject claim.
//
// Parameters:
//   - tokenString: The JWT token string to validate
//
// Returns:
//   - string: The subject extracted from the token's claims, or an empty string if validation fails
//   - error: An error if the token is invalid, uses an unexpected signing method, or lacks a subject claim
func (t *TokenAuthority) ValidateJWT(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(t.jwtSecret), nil
		},
	)

	if err != nil {
		return "", errors.New("failed to parse JWT" + err.Error())
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", errors.New("failed to extract the subject from JWT" + err.Error())
	}

	return subject, nil
}
