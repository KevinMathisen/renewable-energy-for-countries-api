package handlers

import (
	"assignment2/utils/db"
	"assignment2/utils/structs"
	"fmt"
	"log"
	"net/http"
	"time"
)

type RootHandler func(http.ResponseWriter, *http.Request) error

// Handles all errors in the same place.
func (fn RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	//Checks if database is online. If not, all requests are given a standard error response.
	if db.DbState {
		err = fn(w, r)  // Calls original function, then awaits errors to "bubble" back up
		if err == nil { // If there are no errors
			return
		}
	} else {
		usrMsg := fmt.Sprintf("The database is currently unavailable. Please try again later. Reattempting database connection in %v seconds.", time.Since(db.DbRestartTimerStartTime).Round(time.Second)) //Create message with time since timer was activated
		err = structs.NewError(nil, http.StatusServiceUnavailable, usrMsg, "")                                                                                                                             //Return error to user
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
