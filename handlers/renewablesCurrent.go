package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/gateway"
	"assignment2/utils/params"
	"assignment2/utils/structs"
	"fmt"
	"net/http"
	"time"
)

/*
Handler for current endpoint
*/
func RenewablesCurrent(w http.ResponseWriter, r *http.Request) error {
	// Check if database is online. If not, give standard error response.
	if !db.DbState {
		usrMsg := fmt.Sprintf("The database is currently unavailable. Please try again later. Reattempting database connection in %v seconds.", time.Until(db.DbRestartTimerStartTime.Add(1*time.Minute)).Round(time.Second)) //Create message with time since timer was activated
		return structs.NewError(nil, http.StatusServiceUnavailable, usrMsg, "")      
	}
	
	var response []structs.CountryOutput

	// Send error message if request method is not get
	if r.Method != http.MethodGet {
		return structs.NewError(nil, http.StatusNotImplemented, "Invalid method, currently only GET is supported", "User used invalid http method")
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
	go db.InvokeCountry(countries, constants.LATEST_YEAR_DB, constants.LATEST_YEAR_DB)

	// Get current percentage of renewables for countries specified as a list of countryoutput structs
	response, err = getCurrentRenewablesForCountries(w, countries, sortByValue)
	if err != nil {
		return err
	}

	// Check if there was any data for the given request
	if len(response) == 0 {
		return structs.NewError(nil, http.StatusNotFound, "No data available for given request", "No data in database which satisfied the request")
	}

	// Respond with list of CountryOutPut struct encoded as json to user
	err = gateway.RespondToGetRequestWithJSON(w, response, http.StatusOK)
	if err != nil {
		return err
	}

	// Save reponse to cache
	go saveToCache(response, countries, constants.LATEST_YEAR_DB, constants.LATEST_YEAR_DB, r)

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
		renewablesOutput, err = getRenewablesForCountriesByYears(countries, currentYear, currentYear, structs.CreateCountryOutputFromData, sortByValue)
		if err != nil {
			return renewablesOutput, err
		}
	} else {
		// If the user did not specify countires, we get renewables data from all countires in the current year
		renewablesOutput, err = getRenewablesForAllCountriesByYears(currentYear, currentYear, structs.CreateCountryOutputFromData, sortByValue)
		if err != nil {
			return renewablesOutput, err
		}
	}

	return renewablesOutput, nil
}
