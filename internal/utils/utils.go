package utils

import (
	"errors"
	"fmt"
)

const (
	// Empty is empty/blank string.
	Empty = ""
)

// GetValueString returns the specified value, but if the specified
// value is empty, it returns the default value.
func GetValueString(value, def string) string {
	if len(value) == 0 {
		return def
	}
	return value
}

// Error creates a new error by injecting/formatting the inner error.
func Error(value string, inner error) error {
	msg := AppendError(value, inner)
	return errors.New(GetValueString(msg, "error"))
}

// AppendError appends the error message to the value.
func AppendError(value string, err error) string {
	if len(value) > 0 && err != nil {
		return fmt.Sprintf("%s: [%s]", value, err.Error())
	}
	return value
}
