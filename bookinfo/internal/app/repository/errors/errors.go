package errors

import (
	"errors"
	"fmt"
)

func IsErrNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrNotFound{})
}

type ErrNotFound struct {
	Inner error
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("not found: %s", err.Inner.Error())
}

func (err ErrNotFound) Unwrap() error {
	return err.Inner
}

func (err ErrNotFound) Is(target error) bool {
	if target == nil {
		return false
	}
	_, ok := target.(ErrNotFound)
	return ok
}
