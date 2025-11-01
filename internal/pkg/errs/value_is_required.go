package errs

import (
	"errors"
	"fmt"
)

var ErrValueIsRequired = errors.New("value is required")

type ValueIsRequiredError struct {
	ParamName string
}

func NewValueIsRequiredError(paramName string) *ValueIsRequiredError {
	return &ValueIsRequiredError{
		ParamName: paramName,
	}
}

func (e *ValueIsRequiredError) Error() string {
	return fmt.Sprintf("%s: %s", ErrValueIsRequired, e.ParamName)
}

func (e *ValueIsRequiredError) Unwrap() error {
	return ErrValueIsRequired
}
