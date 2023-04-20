package gateway

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"bytes"
	"encoding/json"
	"net/http"
)

/*
Responds to GET request with JSON content and body specified

	w			- Responsewriter
	jsonBody	- Any struct which will be encoded into json and sent as response body
*/
func RespondToGetRequestWithJSON(w http.ResponseWriter, jsonBody interface{}, status int) {
	// Write to content type field in response header
	w.Header().Add("content-type", constants.CONT_TYPE_JSON)

	// Encode content and write to response
	err := json.NewEncoder(w).Encode(jsonBody)
	if err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Manually set response http status code to ok
	w.WriteHeader(status)
}

/*
Create a request and returns response from a specified URL using specified method

	url		- URL to send request to
	method	- Method of request

	return	- http response from request or error
*/
func HttpRequestFromUrl(url string, method string) (http.Response, error) {
	// Create empty response in case of error
	var response http.Response

	// Create request
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return response, err
	}

	// Set content type to empty
	request.Header.Add("content-type", "")

	// Set up client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue http request
	res, err := client.Do(request)
	if err != nil {
		return response, err
	}

	// Return response
	return *res, nil
}

/*
Post content given to webhookURL

	data		- Map of webhook data
	webhookID	- ID of webhook
*/
func PostToWebhook(data map[string]interface{}, webhookID string) {
	var countryName string
	var err error

	// Check if country isoCode is Any, if so return no country name
	if data["country"].(string) == "ANY" {
		countryName = ""
	} else {
		// Find name from isoCode
		countryName, err = GetNameFromIsoCode(data["country"].(string))
		if err != nil {
			return
		}
	}

	// Encode struct into json
	webhookStruct := structs.Webhook{
		WebhookId: webhookID,
		Country:   countryName,
		Calls:     int(data["invocations"].(int64)),
	}

	jsonData, err := json.Marshal(webhookStruct)
	if err != nil {
		return
	}

	// Issue post request to url
	response, err := http.Post(data["url"].(string), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return
	}

	// Close reponse body at end of function
	defer response.Body.Close()
}
