package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/structs"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
		// Get the renwables data from firestore
		renewablesCountry, err := db.GetRenewablesCountryFromFirestore(w, country)
		if err != nil {
			return renewablesOutput, err
		}

		// Create structs with percentage renewable value for each year specified, and save in slice
		outputCountry, err := createCountryOutputFromData(w, renewablesCountry, country, startYear, endYear)
		if err != nil {
			return renewablesOutput, err
		}

		// TODO: Sort output

		renewablesOutput = append(renewablesOutput, outputCountry...)
	}

	// TODO: Sort output

	return renewablesOutput, nil
}

/*
Get renewables data for all counties in the database between start and end year

	w			- Responsewriter for sending error messages
	startYear	- The first year we will get data from
	endYear		- The last year we will get data from

	return		- list of countryouput structs which can will be sent as json in the response, as well as error
*/
func getRenewablesForAllCountriesByYears(w http.ResponseWriter, startYear int, endYear int) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput

	// Get data from all countries from firestore
	countriesData, err := db.GetRenewablesAllCountriesFromFirestore(w)
	if err != nil {
		return nil, err
	}

	// For each country create structs with percentage renewable value for each year specified, and save in slice
	for key, country := range countriesData {
		outputCountry, err := createCountryOutputFromData(w, country, key, startYear, endYear)
		if err != nil {
			return renewablesOutput, err
		}

		// TODO: Sort output

		renewablesOutput = append(renewablesOutput, outputCountry...)
	}

	// TODO: Sort output

	return renewablesOutput, nil
}

/*
Get mean renewables data for all countries in the database between start and end year

	w			- Responsewriter for sending error messages
	startYear	- The first year we will get data from
	endYear		- The last year we will get data from

	return		- list of countryouput structs which can will be sent as json in the response, as well as error
*/
func getMeanRenewablesForAllCountriesByYears(w http.ResponseWriter, startYear int, endYear int) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput

	// Get data from all countries from firestore
	countriesData, err := db.GetRenewablesAllCountriesFromFirestore(w)
	if err != nil {
		return nil, err
	}

	// For each country create a struct with mean value and save in slice
	for key, country := range countriesData {
		outputCountry, err := createMeanCountryOutputFromData(w, country, key, startYear, endYear)
		if err != nil {
			return renewablesOutput, err
		}

		// TODO: Sort output

		renewablesOutput = append(renewablesOutput, outputCountry)
	}

	// TODO: Sort output

	return renewablesOutput, nil
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
		renewablesCountry, err := db.GetRenewablesCountryFromFirestore(w, country)
		if err != nil {
			return renewablesOutput, err
		}

		// Create a struct with mean value for years specified
		outputCountry, err := createMeanCountryOutputFromData(w, renewablesCountry, country, startYear, endYear)
		if err != nil {
			return renewablesOutput, err
		}

		renewablesOutput = append(renewablesOutput, outputCountry)
	}

	// TODO: sort output

	return renewablesOutput, nil
}

/*
Creates a slice of countryOutput structs which can be sendt as response to requests
Goes through each year for a country, filters out the ones we want, and create a struct for each year

	w			- Responsewriter for error handling
	data		- Map which contain name of country and renewable percentages for all years of data
	isoCode		- isoCode of country we are creating structs for
	startYear	- The year in which we want to start returning data from
	endYear		- The year in which we want to stop returning data from

	return		- List of countryOutput structs which can be encoded into Json and sent as reponse to requests
*/
func createCountryOutputFromData(w http.ResponseWriter, data map[string]interface{}, isoCode string, startYear int, endYear int) ([]structs.CountryOutput, error) {
	var output []structs.CountryOutput

	// Save country name as a string
	countryName := data["name"].(string)

	// For each year of country renewables
	for year, percentage := range data {

		// Ignore name field
		if year == "name" {
			continue
		}

		// Try to convert year to an int
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			http.Error(w, "Error when creating data, could not convert year to int", http.StatusInternalServerError)
			return output, err
		}

		// Ignore years outside of scope defined by user
		if yearInt < startYear || yearInt > endYear {
			continue
		}

		// Create countryoutput with year and percentage
		countryOutput := structs.CountryOutput{
			Name:       countryName,
			IsoCode:    isoCode,
			Year:       year,
			Percentage: percentage.(float64),
		}

		// Save each countryoutput to slice
		output = append(output, countryOutput)
	}

	return output, nil
}

/*
Creates a slice of countryOutput structs with Mean value hich can be sendt as response to requests
Goes through each year for a country, filters out the ones we want, and caulcates the mean value for all years.
Then returnes a struct with the mean value.

	w			- Responsewriter for error handling
	data		- Map which contain name of country and renewable percentages for all years of data
	isoCode		- isoCode of country we are creating struct for
	startYear	- The year in which we want to start calculating mean from
	endYear		- The year in which we want to stop calculating mean from

	return		- CountryOutput struct with no year value and mean value as percentage, can be encoded into Json and sent as reponse to requests
*/
func createMeanCountryOutputFromData(w http.ResponseWriter, data map[string]interface{}, isoCode string, startYear int, endYear int) (structs.CountryOutput, error) {
	var percentages []float64

	// Save country name as a string
	countryName := data["name"].(string)

	// For each year of country renewables
	for year, percentage := range data {

		// Ignore name field
		if year == "name" {
			continue
		}

		// Try to convert year to an int
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			http.Error(w, "Error when creating data, could not convert year to int", http.StatusInternalServerError)
			return output, err
		}

		// Ignore years outside of scope defined by user
		if yearInt < startYear || yearInt > endYear {
			continue
		}

		// Add each percentage to a list of all percentages in time range
		percentages = append(percentages, percentage.(float64))
	}

	// Create a countryOutput without year and mean value as percentage
	countryOutput := structs.CountryOutput{
		Name:       countryName,
		IsoCode:    isoCode,
		Percentage: mean(percentages),
	}

	return countryOutput, nil
}

/*
Calculate mean value of list of numbers

	input	- List of float values

	return	- Average of list
*/
func mean(input []float64) float64 {
	// If there are no input
	if len(input) == 0 {
		return 0
	}

	var sum float64

	// Add all values in input to get sum
	for _, value := range input {
		sum += value
	}

	// Return mean value of input
	return sum / float64(len(input))
}
