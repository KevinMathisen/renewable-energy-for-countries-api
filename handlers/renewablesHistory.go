package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"net/http"
	"strconv"
)

func RenewablesHistory(w http.ResponseWriter, r *http.Request) {
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
Get renewables data for each year or mean from year range specified, from countires specified
If the user has no preference for type of data returned, the user will get historical data for each year when specifying countires, and get mean data for all countires if no counties are specified.

	w			- Responsewriter for sending error messages
	countries	- Either a list of countires we want to get data from, or an empty list if we want all
	beginYear	- The first year we will get data from. If -1 we get the default beginYear.
	endYear		- The last year we will get data from. If -1 we get the default endYear (currentyear)
	sortByvalue	- If output is to be sorted by percentage value decending
	getMean		- If the user wants to get mean values, even if countires are specified, this should be true. Does not affect output if no countries are specified, as this will always result in mean value being displayed

	return		- List of countryoutput structs which will be sent as json response. The struct will not have the field "year" defined if mean values are returned.
*/
func getHistoryRenewablesForCountries(w http.ResponseWriter, countries []string, beginYear int, endYear int, sortByValue bool, getMean bool) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput

	// If beginYear not specified set to default
	if beginYear == -1 {
		// TODO: Find beginyear
		beginYear = 1965
	}

	// If endYear not specified, set to default
	if endYear == -1 {
		// TODO: Get current year, as this will be default
		endYear = 2021
	}

	// If countires specified and we don't want mean data, get renewables data from them in year range given
	if len(countries) != 0 && !getMean {
		renewablesOutput, err := getRenewablesForCountriesByYears(w, countries, beginYear, endYear)
		if err != nil {
			return renewablesOutput, err
		}

	} else if len(countries) != 0 && getMean {
		// If countires specified and we want mean data, get renewabled mean data from them in year range given
		renewablesOutput, err := getMeanRenewablesForCountriesByYears(w, countries, beginYear, endYear)
		if err != nil {
			return renewablesOutput, err
		}

	} else if len(countries) == 0 {
		// If no countries specified, get renewables mean data from all in year range given
		renewablesOutput, err := getMeanRenewablesForAllCountriesByYears(w, beginYear, endYear)
		if err != nil {
			return renewablesOutput, err
		}
	} // TODO: Implement functionality for getting renewables history data for all countires which are not mean values

	return renewablesOutput, nil
}

/*
Get parameters from request to renewables history endpoint if any are given

	w	- Responsewriter for error messages
	r	- Request for getting parameters

	return	- Parameters from request. Ints are -1 if empty, bool values are false if empty
*/
func getRenewablesHistoryParameters(w http.ResponseWriter, r *http.Request) (beginYear int, endYear int, sortByValue bool, getMean bool, err error) {
	// Get beginYear param
	begin := (r.URL.Query()).Get("begin")

	// Try to convert string to int
	beginYear, err = strconv.Atoi(begin)
	if err != nil && begin != "" {
		http.Error(w, "Malformed URL, invalid begin parameter set", http.StatusBadRequest)
		return -1, -1, false, false, err
	}

	// Get emdYear param
	end := (r.URL.Query()).Get("end")

	// Try to convert string to int
	endYear, err = strconv.Atoi(end)
	if err != nil && begin != "" {
		http.Error(w, "Malformed URL, invalid end parameter set", http.StatusBadRequest)
		return -1, -1, false, false, err
	}

	// If years set are outside of database scope
	if beginYear < constants.OLDEST_YEAR_DB || endYear > constants.LATEST_YEAR_DB {
		http.Error(w, "Malformed URL, begin and end years have to be between "+strconv.Itoa(constants.OLDEST_YEAR_DB)+" and "+strconv.Itoa(constants.LATEST_YEAR_DB), http.StatusBadRequest)
		return -1, -1, false, false, err
	}

	// Get stortByValue param
	sortBy := (r.URL.Query()).Get("sortByValue")

	// Try to convert string to int
	sortByValue, err = strconv.ParseBool(sortBy)
	if err != nil && begin != "" {
		http.Error(w, "Malformed URL, invalid sortByValue parameter set", http.StatusBadRequest)
		return -1, -1, false, false, err
	}

	// Get getMean param
	mean := (r.URL.Query()).Get("mean")

	// Try to convert string to int
	getMean, err = strconv.ParseBool(mean)
	if err != nil && begin != "" {
		http.Error(w, "Malformed URL, invalid mean parameter set", http.StatusBadRequest)
		return -1, -1, false, false, err
	}

	return beginYear, endYear, sortByValue, getMean, nil
}
