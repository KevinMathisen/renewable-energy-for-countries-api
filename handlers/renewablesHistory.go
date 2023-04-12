package handlers

import (
	"assignment2/utils/structs"
	"net/http"
)

func RenewablesHistory(w http.ResponseWriter, r *http.Request) {
	var response []structs.CountryOutput

	// Send error message if request method is not get
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method, currently only GET is supported", http.StatusNotImplemented)
		return
	}

	// If cache hit, send cached response
	hit := checkCache(w, r)
	if hit {
		return
	}

	// Get the countries we are interested in finding, or empty if everyone
	countries, err := getCountriesToQuery(w, r)
	if err != nil {
		return
	}

	// Get parameters if user specified any
	beginYear, endYear, sortByValue, getMean, err := getRenewablesHistoryParameters(w, r)
	if err != nil {
		return
	}

	// Get the historical percentage of renewables for countires specified as a list of countryoutput structs
	response, err = getHistoryRenewablesForCountries(w, countries)
	if err != nil {
		return
	}

	// Respond with list of countryoutput struct encoded as json to user
	respondToGetRequestWithJSON(w, response)
}

/*
getHistoryRenewablesForCountries(countries, startYear, endYear, wantMean)	//History

	if startYear or endYear not specified
		set them as default values


	if countires specified and not/no mean:
		getRenewablesForCountriesByYears(isoList, startYear, endYear)					Global1

	if countires specified and mean:
		getMeanRenewablesFromCountries(isoList, startYear, endYear)						History1
			for countires
				getMeanRenewablesFromCountry(isoCode, startYear, endYEar)				History1.1


	else if no countires specified and want mean
		getMeanRenewablesFromAllCountries()												History2
			for all countires
				getMeanRenewablesFromCountry(startYear, endYEar)						History2.2

	else if no countries specified and no mean:
		getRenewablesForAllCountiresByYear(startYear, endYear)							Global2



*/
