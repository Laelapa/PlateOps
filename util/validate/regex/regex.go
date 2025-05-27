package regex

import "regexp"

var (
	Hex      = regexp.MustCompile(`^[0-9a-f]+$`)
	Username = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	Email    = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9.]{2,}$`)

	AlphanumericAndBasicSymbols = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	AlphanumericCapitalized     = regexp.MustCompile(`^[A-Z0-9]+$`)
)
