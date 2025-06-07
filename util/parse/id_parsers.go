// Package parse provides utilities for parsing and converting string values
// to their appropriate types with validation. It includes functions for
// parsing IDs and other common data types while ensuring data integrity
// through validation checks before conversion.
package parse

import (
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Laelapa/PlateOps/util/validate"
)

// ID converts a string ID to an int32 after validation.
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

// GTIN converts a string GTIN to pgtype.Text after validation.
// GTIN is expected to be a string representation of a GTIN (Global Trade Item Number).
// GTIN can be 8 to 14 characters long, representing the GTIN-8, GTIN-12, GTIN-13, or GTIN-14 formats.
// Returns the pgtype.Text representation of the GTIN and nil on success,
// or an empty pgtype.Text and propagates the validation error on failure.
func GTIN(gtin string) (pgtype.Text, error) {

	err := validate.GTIN(gtin)
	if err != nil {
		return pgtype.Text{}, err
	}

	// Convert GTIN string to pgtype.Text
	return pgtype.Text{String: gtin, Valid: true}, nil
}
