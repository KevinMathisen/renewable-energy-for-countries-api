package handlers

import (
	"assignment2/utils/db"
	"assignment2/utils/gateway"
	"assignment2/utils/structs"
	"net/http"
)

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
		outputCountry, err := structs.CreateCountryOutputFromData(w, renewablesCountry, country, startYear, endYear)
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
		outputCountry, err := structs.CreateMeanCountryOutputFromData(w, renewablesCountry, country, startYear, endYear)
		if err != nil {
			return renewablesOutput, err
		}

		renewablesOutput = append(renewablesOutput, outputCountry)
	}

	// TODO: sort output

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
		outputCountry, err := structs.CreateCountryOutputFromData(w, country, key, startYear, endYear)
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
		outputCountry, err := structs.CreateMeanCountryOutputFromData(w, country, key, startYear, endYear)
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
	hit, err := db.CheckCacheDBForURL(w, r.URL.Path)
	if err != nil {
		return false, err
	}

	// Cache hit
	if len(hit) != 0 {
		gateway.RespondToGetRequestWithJSON(w, hit)
		return true, nil
	}

	// No cache hit
	return false, nil
}
