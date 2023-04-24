package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/gateway"
	"assignment2/utils/structs"
	"net/http"
	"time"
)

// Create time of service start variable to calculate uptime
var Start time.Time

func Status(w http.ResponseWriter, r *http.Request) error {

	// Send error if request is not GET:
	if r.Method != http.MethodGet {
		return structs.NewError(nil, http.StatusNotImplemented, "Invalid method, currently only GET is supported", "User used invalid http method")
	}

	// Generate status response
	statusRes, err := createStatusResponse(Start)
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
func createStatusResponse(start time.Time) (structs.Status, error) {
	// Get request from countries api
	resCountry, err := gateway.HttpRequestFromUrl(constants.COUNTRIES_API_URL, http.MethodHead)
	if err != nil {
		return structs.Status{}, structs.NewError(err, http.StatusGatewayTimeout, "Error when accessing restcountries api", "")
	}

	// Get request from notification db api
	resDB, err := gateway.HttpRequestFromUrl(constants.FIRESTORE_NOTIFICATION_URL, http.MethodHead)
	if err != nil {
		return structs.Status{}, structs.NewError(err, http.StatusGatewayTimeout, "Error when accesing the database", "")
	}

	// Get amount of webhooks
	amountOfWebhooks, err := db.CountWebhooks()
	if err != nil {
		return structs.Status{}, structs.NewError(err, http.StatusGatewayTimeout, "Error when accesing the database", "")
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
