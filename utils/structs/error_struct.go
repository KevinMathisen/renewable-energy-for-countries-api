package structs

/*
* Struct for wrapping errors for standardized error handling.
*
* OrigErr: Original error message
* StatusCode: Status code to show user
* UsrMessage: Error message to show user.
* DevMessage: Error message to display in logs.
 */
type WrappedError struct {
	OrigErr    error
	StatusCode int
	UsrMessage string
	DevMessage string
}

// Function for creating a new error message.
func NewError(origErr error, statusCode int, userMsg, devMsg string) error {
	return WrappedError{
		OrigErr:    origErr,
		StatusCode: statusCode,
		UsrMessage: userMsg,
		DevMessage: devMsg,
	}
}

// Returns original error in string form
func (err WrappedError) Error() string {
	if err.OrigErr != nil {
		return err.OrigErr.Error()
	}
	return ""
}

// Unwraps error
func (err WrappedError) Unwrap() error {
	return err.OrigErr
}
