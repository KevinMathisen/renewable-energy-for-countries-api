package main

import (
	h "assignment2/handlers"
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	// Set up Firestore, if a connection could not be established
	err := db.InitializeFirestore()
	if err != nil {
		println("Could not connect to Firestore:", err.Error())
		db.ReportDbState(false) //Close service on database failure, automaticly reattempt after 1 minute
	}

	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	// Save start time of service to calculate uptime
	h.Start = time.Now()

	// Handle port assignment
	port := os.Getenv(("PORT"))
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	//!!! OLD HANDLERS
	/* Set up handler enpoints
	http.HandleFunc(constants.DEFAULT_PATH, handlers.Default)
	http.HandleFunc(constants.SERVICE_PATH, handlers.Service)
	http.HandleFunc(constants.RENEWABLES_PATH, handlers.Renewables)
	http.HandleFunc(constants.RENEWABLES_CURRENT_PATH, handlers.RenewablesCurrent)
	http.HandleFunc(constants.RENEWABLES_HISTORY_PATH, handlers.RenewablesHistory)
	http.HandleFunc(constants.NOTIFICATION_PATH, handlers.Notification)
	http.HandleFunc(constants.STATUS_PATH, handlers.Status)
	*/

	// Set up handler endpoints throuh root error handler
	http.Handle(constants.DEFAULT_PATH, h.RootHandler(h.Default))
	http.Handle(constants.SERVICE_PATH, h.RootHandler(h.Service))
	http.Handle(constants.RENEWABLES_PATH, h.RootHandler(h.Renewables))
	http.Handle(constants.RENEWABLES_CURRENT_PATH, h.RootHandler(h.RenewablesCurrent))
	http.Handle(constants.RENEWABLES_HISTORY_PATH, h.RootHandler(h.RenewablesHistory))
	http.Handle(constants.NOTIFICATION_PATH, h.RootHandler(h.Notification))
	http.Handle(constants.STATUS_PATH, h.RootHandler(h.Status))

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
