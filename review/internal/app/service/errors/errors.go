package errors

import (
	"errors"
	"fmt"
)

func IsErrResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrResourceNotFound{})
}

type ErrResourceNotFound struct {
	Inner error
}

func (err ErrResourceNotFound) Error() string {
	return fmt.Sprintf("resource not found: %s", err.Inner.Error())
}
func (err ErrResourceNotFound) Unwrap() error {
	return err.Inner
}

func (err ErrResourceNotFound) Is(target error) bool {
	if target == nil {
		return false
	}
	_, ok := target.(ErrResourceNotFound)
	return ok
}
