// Package parse provides utilities for parsing and converting string values
// to their appropriate types with validation. It includes functions for
// parsing IDs and other common data types while ensuring data integrity
// through validation checks before conversion.
package parse

import (
	"errors"
	"strconv"

	"github.com/Laelapa/PlateOps/util/validate"
)


// ID converts a string ID to an int32 after validation.
// It first validates the ID, then parses the string to a 32-bit integer.
// Returns the parsed int32 ID and nil on success,
// or 0 and propagates the validation error on failure.
func ID(id string) (int32, error) {

	err := validate.ID(id)
	if err != nil {
		return 0, err
	}

	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}

	return int32(intID), nil
}
