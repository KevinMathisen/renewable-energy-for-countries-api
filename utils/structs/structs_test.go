package structs_test

import (
	"errors"
	"net/http"
	"testing"

	"assignment2/utils/structs"
)

func TestNewError(t *testing.T) {
	// Test case 1: Create a new error with a non-nil original error
	origErr := errors.New("original error")
	statusCode := http.StatusBadRequest
	userMsg := "Bad request"
	devMsg := "Request was not well-formed"
	err := structs.NewError(origErr, statusCode, userMsg, devMsg)
	if err == nil {
		t.Errorf("NewError() returned nil, expected non-nil error")
	}
	if err.Error() != "original error" {
		t.Errorf("NewError() returned error message '%s', expected 'original error'", err.Error())
	}
	if e, ok := err.(structs.WrappedError); ok {
		if e.StatusCode != 400 {
			t.Errorf("NewError() returned status code %d, expected %d", e.StatusCode, 400)
		}
		if e.UsrMessage != "Bad request" {
			t.Errorf("NewError() returned user message '%s', expected 'Bad request'", e.UsrMessage)
		}
		if e.DevMessage != "Request was not well-formed" {
			t.Errorf("NewError() returned dev message '%s', expected 'Request was not well-formed'", e.DevMessage)
		}
	} else {
		t.Errorf("NewError() did not return a WrappedError, expected WrappedError")
	}

	// Test case 2: Create a new error with a nil original error
	origErr = nil
	statusCode = http.StatusNotFound
	userMsg = "Not found"
	devMsg = "Requested resource could not be found"
	err = structs.NewError(origErr, statusCode, userMsg, devMsg)
	if err == nil {
		t.Errorf("NewError() returned nil, expected non-nil error")
	}
	if err.Error() != "" {
		t.Errorf("NewError() returned error message '%s', expected empty string", err.Error())
	}
	if e, ok := err.(structs.WrappedError); ok {
		if e.StatusCode != 404 {
			t.Errorf("NewError() returned status code %d, expected %d", e.StatusCode, 404)
		}
		if e.UsrMessage != "Not found" {
			t.Errorf("NewError() returned user message '%s', expected 'Not found'", e.UsrMessage)
		}
		if e.DevMessage != "Requested resource could not be found" {
			t.Errorf("NewError() returned dev message '%s', expected 'Requested resource could not be found'", e.DevMessage)
		}
	} else {
		t.Errorf("NewError() did not return a WrappedError, expected WrappedError")
	}
}

func TestMyError_Error(t *testing.T) {
	// Create a new instance of MyError with a specific error message
	errMsg := "Original Error Message"
	err := structs.NewError(errors.New(errMsg), http.StatusBadRequest, "Bad request.", "Request was not well formed")

	// Call the Error() function and compare the result to the expected error message
	if err.Error() != errMsg {
		t.Errorf("Error() returned '%s', expected '%s'", err.Error(), errMsg)
	}
}
