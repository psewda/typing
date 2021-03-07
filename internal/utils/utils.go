package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
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

// CheckLocalhostURL validates the url. The url must be a valid and localhost url.
func CheckLocalhostURL(u string) error {
	url, err := url.Parse(u)
	if err != nil {
		msg := fmt.Sprintf("the url '%s' is invalid, parsing error", u)
		return Error(msg, err)
	}

	host := strings.ToLower(url.Host)
	if !strings.Contains(host, "localhost") {
		msg := fmt.Sprintf("the url '%s' must be a localhost url", u)
		return errors.New(msg)
	}
	return nil
}
