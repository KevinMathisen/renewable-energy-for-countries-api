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
	err := db.InitializeFirestore(constants.CREDENTIALS_FILE)
	if err != nil {
		log.Println("Could not connect to Firestore:", err.Error())
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

	// Set up handler endpoints throuh root error handler
	http.Handle(constants.DEFAULT_PATH, h.RootHandler(h.Default))
	http.Handle(constants.RENEWABLES_CURRENT_PATH, h.RootHandler(h.RenewablesCurrent))
	http.Handle(constants.RENEWABLES_HISTORY_PATH, h.RootHandler(h.RenewablesHistory))
	http.Handle(constants.NOTIFICATION_PATH, h.RootHandler(h.Notification))
	http.Handle(constants.STATUS_PATH, h.RootHandler(h.Status))

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
