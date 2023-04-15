﻿package handlers

import (
	"assignment2/utils/db"
	"assignment2/utils/gateway"
	"assignment2/utils/structs"
	"net/http"
	"sort"
	"strings"
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
		renewablesCountry, err := db.GetRenewablesCountryFromFirestore(w, country)
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
	countriesData, err := db.GetRenewablesAllCountriesFromFirestore(w)
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
