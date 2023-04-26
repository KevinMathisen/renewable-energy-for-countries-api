package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/gateway"
	"assignment2/utils/params"
	"assignment2/utils/structs"
	"net/http"
)

func RenewablesHistory(w http.ResponseWriter, r *http.Request) error {
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
	countries, err := params.GetCountriesToQuery(w, r, constants.RENEWABLES_HISTORY_PATH)
	if err != nil {
		return err
	}

	// Get parameters if user specified any
	beginYear, endYear, sortByValue, getMean, err := params.GetRenewablesHistoryParameters(w, r)
	if err != nil {
		return err
	}

	// Invoke webhooks
	go db.InvokeCountry(countries, beginYear, endYear)

	// Get the historical percentage of renewables for countires specified as a list of countryoutput structs
	response, err = getHistoryRenewablesForCountries(w, countries, beginYear, endYear, sortByValue, getMean)
	if err != nil {
		return err
	}

	// Check if there was any data for the given request
	if len(response) == 0 {
		return structs.NewError(nil, http.StatusNotFound, "No data available for given request", "No data in database which satisfied the request")
	}

	// Respond with list of countryoutput struct encoded as json to user
	err = gateway.RespondToGetRequestWithJSON(w, response, http.StatusOK)
	if err != nil {
		return err
	}

	// Save reponse to cache
	go saveToCache(response, countries, beginYear, endYear, r)

	return nil
}

/*
Get renewables data for each year or mean from year range specified, from countires specified
If the user has no preference for type of data returned, the user will get historical data for each year when specifying countires, and get mean data for all countires if no counties are specified.

	w			- Responsewriter for sending error messages
	countries	- Either a list of countires we want to get data from, or an empty list if we want all
	beginYear	- The first year we will get data from. If -1 we get the default beginYear.
	endYear		- The last year we will get data from. If -1 we get the default endYear (currentyear)
	sortByvalue	- If output is to be sorted by percentage value decending
	getMean		- If the user wants to get mean values, even if countires are specified, this should be true. Does not affect output if no countries are specified, as this will always result in mean value being displayed

	return		- List of CountryOutPut structs which will be sent as json response. The struct will not have the field "year" defined if mean values are returned.
*/
func getHistoryRenewablesForCountries(w http.ResponseWriter, countries []string, beginYear int, endYear int, sortByValue bool, getMean bool) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput
	var err error

	// If countires specified and we don't want mean data, get renewables data from them in year range given
	if len(countries) != 0 && !getMean {
		renewablesOutput, err = getRenewablesForCountriesByYears(countries, beginYear, endYear, structs.CreateCountryOutputFromData, sortByValue)
		if err != nil {
			return renewablesOutput, err
		}

	} else if len(countries) != 0 && getMean {
		// If countires specified and we want mean data, get renewabled mean data from them in year range given
		renewablesOutput, err = getRenewablesForCountriesByYears(countries, beginYear, endYear, structs.CreateMeanCountryOutputFromData, sortByValue)
		if err != nil {
			return renewablesOutput, err
		}

	} else if len(countries) == 0 {
		// If no countries specified, get renewables mean data from all in year range given
		renewablesOutput, err = getRenewablesForAllCountriesByYears(beginYear, endYear, structs.CreateMeanCountryOutputFromData, sortByValue)
		if err != nil {
			return renewablesOutput, err
		}
	} // TODO: Implement functionality for getting renewables history data for all countires which are not mean values

	return renewablesOutput, nil
}
