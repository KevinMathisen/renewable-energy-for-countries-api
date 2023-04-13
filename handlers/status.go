package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/gateway"
	"assignment2/utils/structs"
	"net/http"
	"strconv"
	"time"
)

// Create time of service start variable to calculate uptime
var Start time.Time

func Status(w http.ResponseWriter, r *http.Request) {

	// Send error if request is not GET:
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method, currently only GET is supported", http.StatusNotImplemented)
		return
	}

	// Generate status response
	statusRes, err := createStatusResponse(Start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle get request
	gateway.RespondToGetRequestWithJSON(w, statusRes)
}

/*
Creates the status json response.

	start - Start time of service
	return - Returns json of status response
*/
func createStatusResponse(start time.Time) (structs.Status, error) {
	// Get request from countries api
	resCountry, err := gateway.HttpRequestFromUrl(constants.COUNTRIES_API_URL, http.MethodHead)
	if err != nil {
		return structs.Status{}, err
	}

	// TODO: Get request from Notification Database in Firebase
	// TODO: Get amount of webhooks from Firebase

	// Initialize the status response struct
	statusResponse := structs.Status{
		CountriesApi:   resCountry.Status,
		NotificationDb: strconv.Itoa(http.StatusNotImplemented),
		Webhooks:       0,
		Version:        constants.VERSION,
		Uptime:         calculateUptimeInSeconds(),
	}

	return statusResponse, nil
}

/*
Calculates time elapsed since start of service in seconds

	return	- Returns in seconds uptime of service
*/
func calculateUptimeInSeconds() float64 {
	return float64(time.Since(Start).Seconds())
}
