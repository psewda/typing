package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	// Empty is empty/blank string.
	Empty = ""
)

// TransportFunc is used to customize the http client's transport
// layer. User can inject a hook function to change the request and response params.
type TransportFunc func(*http.Request) (*http.Response, error)

// RoundTrip forwards the incoming http call to the used defined hook function.
func (f TransportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

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

// ClientWithToken creates a new http client and injects
// the access token in the authorization header.
func ClientWithToken(accessToken string) *http.Client {
	client := http.DefaultClient
	original := http.DefaultTransport
	client.Transport = TransportFunc(func(req *http.Request) (*http.Response, error) {
		req.Header.Add(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", accessToken))
		return original.RoundTrip(req)
	})
	return client
}

// ClientWithJSON creates a new http client. It has an internal transport
// function which build custom response with the specified json and
// http status code. It is useful for testing http calls.
func ClientWithJSON(j string, code int) *http.Client {
	client := http.DefaultClient
	client.Transport = TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: code,
			Body:       ioutil.NopCloser(bytes.NewBufferString(j)),
			Header:     map[string][]string{"Content-Type": {"application/json"}},
		}, nil
	})
	return client
}
