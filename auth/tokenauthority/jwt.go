package tokenauthority

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (t *TokenAuthority) IssueJWT(userID uuid.UUID) (signedToken string, err error) {

	claims := jwt.RegisteredClaims{
		Issuer:    t.jwtIssuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(t.jwtLifetime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(t.jwtSecret))
	if err != nil {
		return "", errors.New("failed to sign JWT" + err.Error())
	}

	return signedToken, nil
}

func (t *TokenAuthority) ValidateJWT(tokenString string) (subject string, err error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {

		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.jwtSecret), nil
	})
	if err != nil {
		return "", errors.New("failed to parse JWT" + err.Error())
	}

	subject, err = token.Claims.GetSubject()
	if err != nil {
		return "", errors.New("failed to extract the subject from JWT" + err.Error())
	}

	return subject, nil
}
