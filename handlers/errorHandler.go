package handlers

import (
	"assignment2/utils/structs"
	"log"
	"net/http"
)

type RootHandler func(http.ResponseWriter, *http.Request) structs.WrappedError

// Handles all errors in the same place.
func (fn RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := fn(w, r)         // Calls original function, then awaits errors to "bubble" back up
	if err.OrigErr == nil { // If there are no errors
		return
	}

	log.Println(err.DevMessage)                   //Logs dev message
	log.Println("\t" + err.Error())               //Logs original error
	http.Error(w, err.UsrMessage, err.StatusCode) //Returns user message and error code

}
