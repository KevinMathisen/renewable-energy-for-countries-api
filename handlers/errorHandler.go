package handlers

import (
	"assignment2/utils/structs"
	"log"
	"net/http"
)

/*
RootHandler is a wrapper for all handlers in the service.
It handles all errors in the same place, and logs them.
*/
type RootHandler func(http.ResponseWriter, *http.Request) error

// Handles all errors in the same place.
func (fn RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r) // Calls original function, then awaits errors to "bubble" back up
	if err == nil { // If there are no errors
		return
	}

	// If error is of type wrappederror, special logging actions will be taken.
	switch e := err.(type) {
	case structs.WrappedError:
		if e.DevMessage != "" {
			log.Println(e.DevMessage) //Logs dev message
		}
		http.Error(w, e.UsrMessage, e.StatusCode) //Returns user message and error code
	default:
		log.Println("Non-wrapped error:")
		http.Error(w, "", http.StatusInternalServerError)
	}
	log.Println("\t" + err.Error()) //Logs original error
}
