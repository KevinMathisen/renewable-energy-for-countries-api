package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/gateway"
	"assignment2/utils/structs"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

/*
Get renewables data for all countries given between start and end year

	w					- Responsewriter for sending error messages
	countries			- A list of countries we want to get data from
	startYear			- The first year we will get data from
	endYear				- The last year we will get data from
	createCountryOutput	- Function for creating the countryOutputs. Alternatives are creating based on years or mean.
	sortByPercentage	- If the output should be sorted by percentage. If not the output is sorted by year and IsoCode.

	return				- list of CountryOutPut structs which will be sent as json in the response, as well as error
*/
func getRenewablesForCountriesByYears(w http.ResponseWriter, countries []string, startYear int, endYear int, createCountryOutput func(http.ResponseWriter, map[string]interface{}, string, int, int) ([]structs.CountryOutput, error), sortByPercentage bool) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput
	var outputNotSorted [][]structs.CountryOutput

	// For each country
	for _, country := range countries {
		// Get the renwables data from Firestore
		renewablesCountry, err := db.GetDocumentFromFirestore(w, country, constants.RENEWABLES_COLLECTION)
		if err != nil {
			return renewablesOutput, err
		}

		// Create structs with percentage renewable value for each year specified, and save in slice
		outputCountry, err := createCountryOutput(w, renewablesCountry, country, startYear, endYear)
		if err != nil {
			return renewablesOutput, err
		}

		outputNotSorted = append(outputNotSorted, outputCountry)
	}

	// Sort output by IsoCode
	renewablesOutput = sortByIsoCode(outputNotSorted)

	// Sort output by percentage if specified
	if sortByPercentage {
		renewablesOutput = sortOutputByPercentage(renewablesOutput)
	}

	return renewablesOutput, nil
}

/*
Get renewables data for all counrties in the database between start and end year

	w					- Responsewriter for sending error messages
	startYear			- The first year we will get data from
	endYear				- The last year we will get data from
	createCountryOutput	- Function for creating the countryOutputs. Alternatives are creating based on years or mean.
	sortByPercentage	- If the output should be sorted by percentage. If not the output is sorted by year and IsoCode.

	return				- list of CountryOutPut structs which can will be sent as json in the response, as well as error
*/
func getRenewablesForAllCountriesByYears(w http.ResponseWriter, startYear int, endYear int, createCountryOutput func(http.ResponseWriter, map[string]interface{}, string, int, int) ([]structs.CountryOutput, error), sortByPercentage bool) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput
	var outputNotSorted [][]structs.CountryOutput

	// Get data from all countries from firestore
	countriesData, err := db.GetAllDocumentInCollectionFromFirestore(w, constants.RENEWABLES_COLLECTION)
	if err != nil {
		return nil, err
	}

	// For each country create structs with percentage renewable value for each year specified, and save in slice
	for key, country := range countriesData {
		outputCountry, err := createCountryOutput(w, country, key, startYear, endYear)
		if err != nil {
			return renewablesOutput, err
		}

		// If there was valid data for the year range, save to slice
		if len(outputCountry) != 0 {
			outputNotSorted = append(outputNotSorted, outputCountry)
		}

	}

	// Sort output by IsoCode
	renewablesOutput = sortByIsoCode(outputNotSorted)

	// Sort output by percentage if specified
	if sortByPercentage {
		renewablesOutput = sortOutputByPercentage(renewablesOutput)
	}

	return renewablesOutput, nil
}

/*
Should check if request is in the cache, then respond with cached response

	w	- Http responsewriter
	r 	- Http request

	return	- bool, true if there was a cache hit
*/
func checkCache(w http.ResponseWriter, r *http.Request) (bool, error) {
	var responseBody []structs.CountryOutput
	var isoCodes []string

	// Create request url path and parameters
	requestURL := strings.Replace((r.URL.Path + r.URL.RawQuery), "/", "\\", -1)

	// Check if request URL and response is in database
	if !db.DocumentInCollection(requestURL, constants.CACHE_COLLECTION) {
		return false, nil
	}

	// Get cached request
	cachedRequest, err := db.GetDocumentFromFirestore(w, requestURL, constants.CACHE_COLLECTION)
	if err != nil {
		return false, err
	}

	// Delete and dont use cached response if it is older than max cache age
	if time.Since(cachedRequest["time"].(time.Time)).Hours() > constants.MAX_CACHE_AGE_IN_HOURS {
		go db.DeleteDocument(w, requestURL, constants.CACHE_COLLECTION)
		return false, err
	}

	// Try to decode reponse saved in cache
	err = json.Unmarshal([]byte(cachedRequest["responseBody"].([]uint8)), &responseBody)
	if err != nil {
		return false, err
	}

	// Try to decode isoCodes saved in cache
	err = json.Unmarshal([]byte(cachedRequest["isoCodes"].([]uint8)), &isoCodes)
	if err != nil {
		return false, err
	}

	// Invoke webhooks
	go db.InvokeCountry(isoCodes)

	// Answer request with cached response
	err = gateway.RespondToGetRequestWithJSON(w, responseBody, http.StatusOK)
	if err != nil {
		return false, err
	}

	return true, nil
}

/*
Saves a request and its corresponding response to firestore, along with the response timestamp

	responseBody	- reponse we will save
	r				- http.request for getting the url of the request
*/
func saveToCache(responseBody []structs.CountryOutput, isoCodes []string, r *http.Request) {

	// create request id by path and parameters
	requestID := strings.Replace((r.URL.Path + r.URL.RawQuery), "/", "\\", -1)

	// Encode country strycts into json
	responseEncoded, err := json.Marshal(responseBody)
	if err != nil {
		log.Println("Error when encoding country response to json for caching")
		return
	}

	// Encode isoCodes into json
	isoCodesEncoded, err := json.Marshal(isoCodes)
	if err != nil {
		log.Println("Error when encoding isoCodes to json for caching")
		return
	}

	// Create cache map
	cachedResponse := map[string]interface{}{
		"responseBody": responseEncoded,
		"isoCodes":     isoCodesEncoded,
		"time":         firestore.ServerTimestamp,
	}

	// Save reponse with url path and parameters to firestore
	err = db.AppendDocumentToFirestore(requestID, cachedResponse, constants.CACHE_COLLECTION)
	if err != nil {
		return
	}
}

/*
Sorts a list of countryoutput by their isoCode alphabetically

	input	- Slice of slices where each subslice is countryoutputs from one country

	return 	- Slice of countryoutput sorted by isoCode
*/
func sortByIsoCode(input [][]structs.CountryOutput) []structs.CountryOutput {
	var output []structs.CountryOutput

	// Sort by isoCode alphabetically
	sort.Slice(input, func(i, j int) bool {
		return strings.Compare(input[i][0].IsoCode, input[j][0].IsoCode) == -1
	})

	// Append each subslice for each country to one slice
	for _, country := range input {
		output = append(output, country...)
	}

	return output
}

/*
Sorts a slice of countryputputs by percentage descending

	input	- Slice of countryOutput structs to be sorted

	return	- Slice of countryOutput structs sorted
*/
func sortOutputByPercentage(input []structs.CountryOutput) []structs.CountryOutput {
	sort.Slice(input, func(i, j int) bool {
		return input[i].Percentage > input[j].Percentage
	})
	return input
}
