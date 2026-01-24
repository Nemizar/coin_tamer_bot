package errs

import (
	"errors"
	"fmt"
)

var ErrEntityAlreadyExists = errors.New("entity already exists")

type EntityAlreadyExistsError struct {
	Entity string
	Field  string
	Value  string
}

func NewEntityAlreadyExistsError(entity, field, value string) *EntityAlreadyExistsError {
	return &EntityAlreadyExistsError{
		Entity: entity,
		Field:  field,
		Value:  value,
	}
}

func (e *EntityAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s: %s with %s '%s' already exists",
		ErrEntityAlreadyExists, e.Entity, e.Field, e.Value)
}

func (e *EntityAlreadyExistsError) Unwrap() error {
	return ErrEntityAlreadyExists
}
