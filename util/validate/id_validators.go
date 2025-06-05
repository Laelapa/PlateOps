package validate

import (
	"errors"
	"strconv"
)

// ID validates that the input string represents a valid positive int32 ID.
// It checks that the ID is not empty, can be parsed as a 32-bit integer,
// and has a value greater than zero.
//
// Parameters:
//   - id: The string representation of the ID to validate
//
// Returns:
//   - error: Returns an error if the ID is missing, has invalid format,
//     or is not a positive value. Returns nil if validation passes.
//
// Possible errors:
//   - "missing ID": when the input string is empty
//   - "invalid ID format": when the string cannot be parsed as an int32
//   - "invalid ID value": when the parsed integer is zero or negative
func ID(id string) error {
	if len(id) == 0 {
		return errors.New("missing ID")
	}

	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return errors.New("invalid ID format")
	}

	if intID <= 0 {
		return errors.New("invalid ID value")
	}

	return nil
}
