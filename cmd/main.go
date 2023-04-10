package main

import (
	"assignment2/handlers"
	"assignment2/utils/constants"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	// Save start time of service to calculate uptime
	handlers.Start = time.Now()

	// Handle port assignment
	port := os.Getenv(("PORT"))
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Set up handler enpoints
	http.HandleFunc(constants.DEFAULT_PATH, handlers.Default)
	http.HandleFunc(constants.SERVICE_PATH, handlers.Service)
	http.HandleFunc(constants.RENEWABLES_PATH, handlers.Renewables)
	http.HandleFunc(constants.RENEWABLES_CURRENT_PATH, handlers.RenewablesCurrent)
	http.HandleFunc(constants.RENEWABLES_HISTORY_PATH, handlers.RenewablesHistory)
	http.HandleFunc(constants.NOTIFICATION_PATH, handlers.Notification)
	http.HandleFunc(constants.STATUS_PATH, handlers.Status)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
