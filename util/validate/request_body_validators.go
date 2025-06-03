// TODO: Note to future self: Good place to experiment with generics
package validate

import (
	"errors"

	"github.com/Laelapa/PlateOps/util/validate/regex"

	guuid "github.com/google/uuid"
)

const (
	usernameMinLength     = 5
	usernameMaxLength     = 24
	emailMaxLength        = 254
	passwordMinLength     = 8
	passwordMaxLength     = 72 // 72 mainly due to bcrypt
	foodUnitTypeMaxLength = 50
	gtinMaxLength         = 14
	realMin               = 0.0
	realMax               = 1000000.0
	stringMaxLength       = 255
	textMaxLength         = 5000 // TEXT database fields
)

func RefreshToken(token string, length int) error {

	if token == "" {
		return ErrEmptyRefreshToken
	}
	if len(token) != length*2 { // hex representation is two characters per byte
		return ErrInvalidRefreshToken
	}

	if !regex.Hex.MatchString(token) {
		return ErrInvalidRefreshToken
	}

	return nil
}

func Username(u string) error {
	if len(u) < usernameMinLength || len(u) > usernameMaxLength {
		return ErrInvalidUsername
	}
	if !regex.Username.MatchString(u) {
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
	if !regex.Email.MatchString(e) {
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

func StringRequired(s string) error {
	if s == "" {
		return ErrRequiredValueOmitted
	}

	return String(s)
}

func String(s string) error {
	if len(s) > stringMaxLength {
		return ErrStringTooLong
	}

	return nil
}

func Text(t string) error {
	if len(t) > textMaxLength {
		return ErrStringTooLong
	}

	return nil
}

func GTIN(gtin string) error {
	if len(gtin) > gtinMaxLength {
		return ErrGtinTooLong
	}

	// GTIN can be numeric or alphanumeric, but must not contain special characters.
	// Also, for consistency capitalization is enforced.
	if !regex.AlphanumericCapitalized.MatchString(gtin) {
		return ErrInvalidGtinFormat
	}

	return nil
}

func UnitType(t string) error {
	if len(t) > foodUnitTypeMaxLength {
		return ErrStringTooLong
	}

	// Unit type can be alphanumeric and underscores, but no special characters.
	if !regex.AlphanumericAndBasicSymbols.MatchString(t) {
		return ErrInvalidGtinFormat
	}

	m := map[string]bool{
		"grams":    true,
		"ml":       true,
		"items":    true,
		"portions": true,
	}

	if !m[t] {
		return ErrInvalidUnitType
	}

	return nil
}

func NonNegative[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](n T) error {
	if n < 0 {
		return ErrValueOutOfBounds
	}

	return nil
}

func Positive[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](n T) error {
	if n <= 0 {
		return ErrValueOutOfBounds
	}

	return nil
}

func Real(x float32) error {
	if x < realMin || x > realMax {
		return ErrValueOutOfBounds
	}

	return nil
}
