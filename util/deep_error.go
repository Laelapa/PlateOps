package util

import (
	"errors"
	"fmt"
)

type DeepError struct {
	BusinessErr  error
	TechnicalErr error
}

func (e *DeepError) Error() string {
	return fmt.Sprintf("%v: %v", e.BusinessErr, e.TechnicalErr)
}

func (e *DeepError) Unwrap() error {
	return e.TechnicalErr
}

func (e *DeepError) Is(target error) bool {
	return errors.Is(e.BusinessErr, target) || errors.Is(e.TechnicalErr, target)
}
