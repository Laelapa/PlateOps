package validate

import (
	"errors"
	"strconv"
)

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
