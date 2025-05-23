package validate

import (
	"errors"
	"regexp"

	guuid "github.com/google/uuid"
)

const (
	usernameMinLength = 5
	usernameMaxLength = 24
	emailMaxLength    = 254
	passwordMinLength = 8
	passwordMaxLength = 72 // 72 mainly due to bcrypt
)

var (
	ErrEmptyRefreshToken   = errors.New("refresh token is empty")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidUsername     = errors.New("invalid username")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrEmptyEmail          = errors.New("email is empty")
	ErrInvalidEmail        = errors.New("invalid email")

	hexRegex      = regexp.MustCompile(`^[0-9a-f]+$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9.]{2,}$`)
)

func RefreshToken(token string, length int) error {

	if token == "" {
		return ErrEmptyRefreshToken
	}
	if len(token) != length*2 { // hex representation is two characters per byte
		return ErrInvalidRefreshToken
	}

	if !hexRegex.MatchString(token) {
		return ErrInvalidRefreshToken
	}

	return nil
}

func Username(u string) error {
	if len(u) < usernameMinLength || len(u) > usernameMaxLength {
		return ErrInvalidUsername
	}
	if !usernameRegex.MatchString(u) {
		return ErrInvalidUsername
	}

	return nil
}

func Email(e string) error {
	if e == "" {
		return ErrEmptyEmail
	}
	if len(e) > emailMaxLength {
		return ErrInvalidEmail
	}
	if !emailRegex.MatchString(e) {
		return ErrInvalidEmail
	}

	return nil
}

func Password(p string) error {
	if len(p) < passwordMinLength || len(p) > passwordMaxLength {
		return ErrInvalidPassword
	}

	return nil
}

func UUID(id guuid.UUID) error {
	if id == guuid.Nil {
		return errors.New("invalid or missing UUID")
	}

	return nil
}
