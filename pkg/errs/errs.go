package errs

// NotFoundError is NotFound error struct.
type NotFoundError struct {
	message string
}

// Error returns error message as string.
func (e *NotFoundError) Error() string {
	return e.message
}

// UnauthorizedError is Unauthorized error struct.
type UnauthorizedError struct {
	message string
}

// Error returns error message as string.
func (e *UnauthorizedError) Error() string {
	return e.message
}

// NewNotFoundError creates new instance of NotFoundError.
func NewNotFoundError(m string) *NotFoundError {
	return &NotFoundError{
		message: m,
	}
}

// NewUnauthorizedError creates new instance of UnauthorizedError.
func NewUnauthorizedError() *UnauthorizedError {
	return &UnauthorizedError{
		message: "authorization token is invalid or expired",
	}
}
