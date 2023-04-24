package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/gateway"
	"assignment2/utils/params"
	"assignment2/utils/structs"
	"net/http"
)

func RenewablesCurrent(w http.ResponseWriter, r *http.Request) error {
	var response []structs.CountryOutput

	// Send error message if request method is not get
	if r.Method != http.MethodGet {
		return structs.NewError(nil, http.StatusNotImplemented, "Invalid method, currently only GET is supported", "")
	}

	// If cache hit, send cached response
	hit, err := checkCache(w, r)
	if hit || err != nil {
		return err
	}

	// Get the countries we are interested in finding, or empty if everyone
	countries, err := params.GetCountriesToQuery(w, r, constants.RENEWABLES_CURRENT_PATH)
	if err != nil {
		return err
	}

	sortByValue, err := params.GetBoolParameterFromRequest(w, r, "sortByValue")
	if err != nil {
		return err
	}

	// Invoke webhooks
	go db.InvokeCountry(countries)

	// Get current percentage of renewables for countries specified as a list of countryoutput structs
	response, err = getCurrentRenewablesForCountries(w, countries, sortByValue)
	if err != nil {
		return err
	}

	// Respond with list of CountryOutPut struct encoded as json to user
	err = gateway.RespondToGetRequestWithJSON(w, response, http.StatusOK)
	if err != nil {
		return err
	}

	// Save reponse to cache
	go saveToCache(response, countries, r)

	return nil
}

/*
Get renewables data for the current year from specified countires or all countries

	w			- Responsewriter for sending error messages
	countires	- Either a list of countries we want to get data from, or an empty list if we want all

	return		- Returns a list of CountryOutPut structs which can will be sent as json in the response
*/
func getCurrentRenewablesForCountries(w http.ResponseWriter, countries []string, sortByValue bool) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput
	var err error

	// Get current year
	// TODO: Get current Year
	currentYear := constants.LATEST_YEAR_DB

	// If the users specified countries, get renewables data from them in the current year
	if len(countries) != 0 {
		renewablesOutput, err = getRenewablesForCountriesByYears(w, countries, currentYear, currentYear, structs.CreateCountryOutputFromData, sortByValue)
		if err != nil {
			return renewablesOutput, err
		}
	} else {
		// If the user did not specify countires, we get renewables data from all countires in the current year
		renewablesOutput, err = getRenewablesForAllCountriesByYears(w, currentYear, currentYear, structs.CreateCountryOutputFromData, sortByValue)
		if err != nil {
			return renewablesOutput, err
		}
	}

	return renewablesOutput, nil
}
