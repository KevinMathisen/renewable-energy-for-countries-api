package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/gateway"
	"assignment2/utils/structs"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Variable for time of service start, used to calculate uptime
var Start time.Time

/*
Handler for status endpoint
*/
func Status(w http.ResponseWriter, r *http.Request) error {

	// Send error if request is not GET:
	if r.Method != http.MethodGet {
		return structs.NewError(nil, http.StatusNotImplemented, "Invalid method, currently only GET is supported", "User used invalid http method")
	}

	// Generate status response
	statusRes, err := createStatusResponse(constants.COUNTRIES_API_URL, Start)
	if err != nil {
		return err
	}

	// Handle get request
	err = gateway.RespondToGetRequestWithJSON(w, statusRes, http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

/*
Creates the status json response.

	start - Start time of service
	return - Returns json of status response
*/
func createStatusResponse(countriesApiURL string, start time.Time) (structs.Status, error) {
	// Get request from countries api
	resCountry, err := gateway.HttpRequestFromUrl(countriesApiURL, http.MethodHead)
	if err != nil {
		statusCode := http.StatusServiceUnavailable
		resCountry.Status = fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
		log.Printf("Could not get response from countries api: %v", err.Error())
	}

	// Get request from notification db api
	resDB, err := db.GetDbResponse()
	if err != nil {
		log.Printf("Could not get response from notification db: %v", err.Error())
	}

	// Get amount of webhooks
	amountOfWebhooks, err := db.CountWebhooks()
	if err != nil {
		amountOfWebhooks = -1
		log.Printf("Could not count webhooks: %v", err.Error())
	}

	// Initialize the status response struct
	statusResponse := structs.Status{
		CountriesApi:   resCountry.Status,
		NotificationDb: resDB.Status,
		Webhooks:       amountOfWebhooks,
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
