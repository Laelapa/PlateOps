package validate

import "errors"

var (
	ErrEmptyRefreshToken    = errors.New("refresh token is empty")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidUsername      = errors.New("invalid username")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrEmptyEmail           = errors.New("email is empty")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrRequiredValueOmitted = errors.New("required value omitted")
	ErrStringTooLong        = errors.New("string exceeds maximum length")
	ErrGtinTooLong          = errors.New("GTIN exceeds maximum length")
	ErrInvalidGtinFormat    = errors.New("invalid GTIN format")
	ErrInvalidUnitType      = errors.New("invalid UnitType")
	ErrValueOutOfBounds     = errors.New("a value is out of bounds")
)
