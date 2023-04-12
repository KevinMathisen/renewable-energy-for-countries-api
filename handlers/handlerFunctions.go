package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"encoding/json"
	"errors"
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
Get country code or name, and neighbours parameter from request, then returns appropiate list of countries from these

	w	- Responsewriter
	r	- Request

	return	- Either empty list of no country specified, or one country, or country and it's neighbours
*/
func getCountriesToQuery(w http.ResponseWriter, r *http.Request) ([]string, error) {
	var countries []string

	// Get country code or name from request
	countryCodeOrName, err := getCountryCodeOrNameFromRequest(r)
	if err != nil {
		return nil, err
	}

	// Get neigbour bool from request if it is specified
	neigbours, err := getNeigboursParameterFromRequest(r)
	if err != nil {
		return nil, err
	}

	// If user didn't specify any country
	if countryCodeOrName == "" {
		return nil, nil
	}

	// If the user specified the name only
	if len(countryCodeOrName) != 3 {
		// TODO: Implement how to get the ISO code if name is given

	} else if isoCodeInDB(countryCodeOrName) {
		// Else if the user specified ISO code and it exists in the database, add the code the list of countries
		countries = append(countries, countryCodeOrName)
	}

	// If the user specified the neighbour parameter
	if neigbours {
		// TODO: Get neighbour ISO code with Restcountries API

		// TODO: Check if each isoCode is in database, if so add to list of countires
	}

	// If no countries existed in the database
	if len(countries) == 0 {
		http.Error(w, "No country with given ISO code or name exists in our service", http.StatusNotFound)
		return nil, errors.New("No country with given ISO code or name exists in our service")
	}

	return countries, nil

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

/*
Get mean renewables data for all countries given between start and end year

	w			- Responsewriter for sending error messages
	countries	- A list of countries we want to get data from
	startYear	- The first year we will get data from
	endYear		- The last year we will get data from

	return		- list of countryouput structs with no year and percentage as mean value, which can will be sent as json in the response, as well as error
*/
func getMeanRenewablesForCountriesByYears(w http.ResponseWriter, countries []string, startYear int, endYear int) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput

	// For each country
	for _, country := range countries {
		// Get the renwables data from year range
		renewablesCountry, err := getRenewablesYearsFromCountry(w, country, startYear, endYear) // TODO: Create function
		if err != nil {
			return renewablesOutput, err
		}
		// TODO: Create a CountryOutput struct with mean percentage, and correct iso code and country name
		// TODO: Create a way to calculate mean percentage given list of renewables data

		renewablesOutput = append(renewablesOutput, outputCountry)
	}
}
