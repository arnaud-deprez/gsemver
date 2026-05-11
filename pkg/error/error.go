package error

import "fmt"

// Error is a typical error representation that can happen during the version bump process
type Error struct {
	message string
	cause   error
}

// Error formats VersionError into a string
func (e Error) Error() string {
	if e.cause == nil {
		return e.message
	}
	return fmt.Sprintf("%s caused by: %v", e.message, e.cause)
}

// NewError create an error based on a format error message
func NewError(format string, args ...interface{}) Error {
	return NewErrorC(nil, format, args...)
}

// NewErrorC create an error based on a cause error and a format error message
func NewErrorC(cause error, format string, args ...interface{}) Error {
	return Error{message: fmt.Sprintf(format, args...), cause: cause}
}
