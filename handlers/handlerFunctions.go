package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"encoding/json"
	"net/http"
)

/*
Responds to GET request with JSON content and body specified

	w			- Responsewriter
	jsonBody	- Any struct which will be encoded into json and sent as response body
*/
func respondToGetRequestWithJSON(w http.ResponseWriter, jsonBody interface{}) {
	// Write to content type field in response header
	w.Header().Add("content-type", constants.CONT_TYPE_JSON)

	// Encode content and write to response
	err := json.NewEncoder(w).Encode(jsonBody)
	if err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Manually set response http status code to ok
	w.WriteHeader(http.StatusOK)
}

/*
Create a request and returns response from a specified URL using specified method

	url		- URL to send request to
	method	- Method of request

	return	- http response from request or error
*/
func httpRequestFromUrl(url string, method string) (http.Response, error) {
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
Get renewables data for all countries given between start and end year

	w			- Responsewriter for sending error messages
	countries	- A list of countries we want to get data from
	startYear	- The first year we will get data from
	endYear		- The last year we will get data from

	return		- list of countryouput structs which can will be sent as json in the response, as well as error
*/
func getRenewablesForCountriesByYears(w http.ResponseWriter, countries []string, startYear int, endYear int) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput

	// For each country
	for _, country := range countries {
		// Get the renwables data from year range
		renewablesCountry, err := getRenewablesYearsFromCountry(w, country, startYear, endYear) // TODO: Create function
		if err != nil {
			return renewablesOutput, err
		}
		// TODO: Create a CountryOutput struct with renewabled and correct iso code and country name

		renewablesOutput = append(renewablesOutput, outputCountry)
	}
}

/*
Get renewables data for all counties in the database between start and end year

	w			- Responsewriter for sending error messages
	startYear	- The first year we will get data from
	endYear		- The last year we will get data from

	return		- list of countryouput structs which can will be sent as json in the response, as well as error
*/
func getRenewablesForAllCountriesByYears(w http.ResponseWriter, startYear int, endYear int) ([]structs.CountryOutput, error) {
	// var renewablesOutput []structs.CountryOutput

	// TODO: For all countries, get data, get country name, append to renewablesOutput, and return
}

/*
Should check if request is in the cache, then respond with cached response

	w	- Http responsewriter
	r 	- Http request

	return	- bool, true if there was a cache git
*/
func checkCache(w http.ResponseWriter, r *http.Request) (bool, error) {
	var hit []structs.CountryOutput

	// Check if request URL and response is in database
	hit, err := checkCacheDBForURL(w, r.URL.Path)
	if err != nil {
		return false, err
	}

	// Cache hit
	if len(hit) != 0 {
		respondToGetRequestWithJSON(w, hit)
		return true, nil
	}

	// No cache hit
	return false, nil
}
