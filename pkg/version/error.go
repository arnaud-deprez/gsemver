package version

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
	return fmt.Sprintf("%s caused by: '%v'", e.message, e.cause)
}

// NewError create an error based on a format error message
func newError(format string, args ...interface{}) Error {
	return newErrorC(nil, format, args...)
}

// NewErrorC create an error based on a cause error and a format error message
func newErrorC(cause error, format string, args ...interface{}) Error {
	return Error{message: fmt.Sprintf(format, args...), cause: cause}
}
