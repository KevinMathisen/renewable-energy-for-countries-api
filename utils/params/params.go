package params

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"errors"
	"net/http"
	"strconv"
)

/*
Get country code or name, and neighbours parameter from request, then returns appropiate list of countries from these

	w	- Responsewriter
	r	- Request

	return	- Either empty list of no country specified, or one country, or country and it's neighbours
*/
func GetCountriesToQuery(w http.ResponseWriter, r *http.Request) ([]string, error) {
	var countries []string

	// Get country code or name from request
	countryCodeOrName, err := getCountryCodeOrNameFromRequest(w, r)
	if err != nil {
		return nil, err
	}

	// Get neigbour bool from request if it is specified
	neighbours, err := getNeighboursParameterFromRequest(w, r)
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

	} else if db.IsoCodeInDB(countryCodeOrName) {
		// Else if the user specified ISO code and it exists in the database, add the code the list of countries
		countries = append(countries, countryCodeOrName)
	}

	// If the user specified the neighbour parameter
	if neighbours {
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
Get parameters from request to renewables history endpoint if any are given

	w	- Responsewriter for error messages
	r	- Request for getting parameters

	return	- Parameters from request. Ints are -1 if empty, bool values are false if empty
*/
func GetRenewablesHistoryParameters(w http.ResponseWriter, r *http.Request) (beginYear int, endYear int, sortByValue bool, getMean bool, err error) {
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

func getCountryCodeOrNameFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	return "", nil
}

func getNeighboursParameterFromRequest(w http.ResponseWriter, r *http.Request) (bool, error) {
	return false, nil
}
