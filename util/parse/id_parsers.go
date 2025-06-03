package parse

import (
	"errors"
	"strconv"

	"github.com/Laelapa/PlateOps/util/validate"
)

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
