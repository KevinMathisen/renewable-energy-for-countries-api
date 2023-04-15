package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/gateway"
	"assignment2/utils/params"
	"assignment2/utils/structs"
	"net/http"
)

func RenewablesCurrent(w http.ResponseWriter, r *http.Request) {
	var response []structs.CountryOutput

	// Send error message if request method is not get
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method, currently only GET is supported", http.StatusNotImplemented)
		return
	}

	// If cache hit, send cached response
	hit, err := checkCache(w, r)
	if hit || err != nil {
		return
	}

	// Get the countries we are interested in finding, or empty if everyone
	countries, err := params.GetCountriesToQuery(w, r, constants.RENEWABLES_CURRENT_PATH)
	if err != nil {
		return
	}

	// Get current percentage of renewables for countries specified as a list of countryoutput structs
	response, err = getCurrentRenewablesForCountries(w, countries)
	if err != nil {
		return
	}

	// Respond with list of CountryOutPut struct encoded as json to user
	gateway.RespondToGetRequestWithJSON(w, response)
}

/*
Get renewables data for the current year from specified countires or all countries

	w			- Responsewriter for sending error messages
	countires	- Either a list of countries we want to get data from, or an empty list if we want all

	return		- Returns a list of CountryOutPut structs which can will be sent as json in the response
*/
func getCurrentRenewablesForCountries(w http.ResponseWriter, countries []string) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput
	var err error

	// Get current year
	// TODO: Get current Year
	currentYear := constants.LATEST_YEAR_DB

	// If the users specified countries, get renewables data from them in the current year
	if len(countries) != 0 {
		renewablesOutput, err = getRenewablesForCountriesByYears(w, countries, currentYear, currentYear)
		if err != nil {
			return renewablesOutput, err
		}
	} else {
		// If the user did not specify countires, we get renewables data from all countires in the current year
		renewablesOutput, err = getRenewablesForAllCountriesByYears(w, currentYear, currentYear)
		if err != nil {
			return renewablesOutput, err
		}
	}

	return renewablesOutput, nil
}
