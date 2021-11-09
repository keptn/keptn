package db

import "fmt"

type InvalidEventFilterError struct {
	msg string
}

func NewInvalidEventFilterError(msg string) *InvalidEventFilterError {
	invalidFilterError := &InvalidEventFilterError{
		msg: msg,
	}

	return invalidFilterError
}

func (filterError InvalidEventFilterError) Error() string {
	return fmt.Sprintf("invalid event filter: %s", filterError.msg)
}
