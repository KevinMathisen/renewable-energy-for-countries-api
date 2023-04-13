package handlers

import (
	"log"
	"net/http"
)

type RootHandler func(http.ResponseWriter, *http.Request) WrappedError

// Handles all errors in same place.
func (fn RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := fn(w, r)         // Calls original function, then awaits error to "bubble" back up
	if err.OrigErr == nil { // If there are no errors
		return
	}

	log.Println(err.DevMessage)                   //Logs dev message
	log.Println("\t" + err.Error())               //Logs original error
	http.Error(w, err.UsrMessage, err.StatusCode) //Returns user message and error code

}

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
	return err.OrigErr.Error()
}
