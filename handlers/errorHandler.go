package handlers

import (
	"log"
	"net/http"
)

type RootHandler func(http.ResponseWriter, *http.Request) error

// Handles all errors in same place.
func (fn RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := fn(w, r) // Calls original function, then awaits error to "bubble" back up
	if err == nil { // If there are no errors
		return
	}

	log.Println(err.DevMessage)                   //Logs dev message
	log.Println("\t" + err.Error())               //Logs original error
	http.Error(w, err.UsrMessage, err.StatusCode) //Returns user message and error code

}
