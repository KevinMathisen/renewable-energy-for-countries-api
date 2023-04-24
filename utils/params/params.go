package params

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/div"
	"assignment2/utils/gateway"
	"assignment2/utils/structs"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
Get country code or name, and neighbours parameter from request, then returns appropiate list of countries from these

	w	- Responsewriter
	r	- Request

	return	- Either empty list of no country specified, or one country, or country and its neighbours
*/
func GetCountriesToQuery(w http.ResponseWriter, r *http.Request, path string) ([]string, error) {
	var countries []string
	var countriesInDB []string

	// Get country code or name from request
	countryCodeOrName, err := getCountryCodeOrNameFromRequest(w, r, path)
	if err != nil {
		return nil, err
	}

	// Get neigbour bool from request if it is specified
	neighbours, err := GetBoolParameterFromRequest(w, r, "neighbours")
	if err != nil {
		return nil, err
	}

	// If user didn't specify any country
	if countryCodeOrName == "" {
		return nil, nil
	}

	// If the user specified the name only
	if len(countryCodeOrName) != 3 {
		// Get isoCode from name
		isoCode, err := gateway.GetIsoCodeFromName(countryCodeOrName)
		if err != nil {
			return nil, err
		}

		countries = append(countries, isoCode)

	} else {
		// Else if the user specified ISO code, add the code the list of countries
		countries = append(countries, countryCodeOrName)
	}

	// If the user specified the neighbour parameter, get neighbour ISO code with Restcountries API
	if neighbours {
		country, err := gateway.GetCountryByIso(countries[0]) //Get the country object
		if err != nil {
			return nil, err
		}

		//Add borders to countries list, this will most likely include a lot of duplicates.
		countries = append(countries, country.Borders...)

		temp := countries                      //Take backup of countries
		countries = div.RemoveDuplicates(temp) //Remove duplicates from countries
	}

	// Check if each country exists in the database
	for _, isoCode := range countries {
		if db.DocumentInCollection(isoCode, constants.RENEWABLES_COLLECTION) {
			countriesInDB = append(countriesInDB, isoCode)
		}
	}

	// If no countries existed in the database
	if len(countriesInDB) == 0 {
		return nil, structs.NewError(nil, http.StatusNotFound, "No country with given ISO code or name exists in our service", "")
	}

	return countriesInDB, nil

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

	// If the parameter is not specified
	if begin == "" {
		beginYear = -1
	} else {
		// If the parameter is specified
		// Try to convert string to int
		beginYear, err = strconv.Atoi(begin)
		if err != nil && begin != "" {
			http.Error(w, "Malformed URL, invalid begin parameter set", http.StatusBadRequest)
			return -1, -1, false, false, err
		}
	}

	// Get endYear param
	end := (r.URL.Query()).Get("end")

	// If the parameter is not specified
	if end == "" {
		endYear = -1
	} else {
		// If the parameter is specified
		// Try to convert string to int
		endYear, err = strconv.Atoi(end)
		if err != nil && end != "" {
			http.Error(w, "Malformed URL, invalid end parameter set", http.StatusBadRequest)
			return -1, -1, false, false, err
		}
	}

	// If years set are outside of database scope
	if (beginYear < constants.OLDEST_YEAR_DB && beginYear != -1) || (endYear > constants.LATEST_YEAR_DB && endYear != -1) {
		http.Error(w, "Malformed URL, begin and end years have to be between "+strconv.Itoa(constants.OLDEST_YEAR_DB)+" and "+strconv.Itoa(constants.LATEST_YEAR_DB), http.StatusBadRequest)
		return -1, -1, false, false, err
	}

	// Get sortByValue param
	sortByValue, err = GetBoolParameterFromRequest(w, r, "sortByValue")
	if err != nil {
		return -1, -1, false, false, err
	}

	// Get getMean param
	getMean, err = GetBoolParameterFromRequest(w, r, "mean")
	if err != nil {
		return -1, -1, false, false, err
	}

	return beginYear, endYear, sortByValue, getMean, nil
}

/*
Get country code or name from the requests url

	w		- Responsewriter
	r		- Request
	path	- Path of endpoint used for giving correct error handling message

	return 	- Country code or name, or empty
*/
func getCountryCodeOrNameFromRequest(w http.ResponseWriter, r *http.Request, path string) (string, error) {

	// Split path into args
	args := strings.Split(r.URL.Path, "/")

	// Check if URL is correctly formatted
	if len(args) != 6 && len(args) != 7 {
		http.Error(w, "Malformed URL, Expecting format "+path+"{country?}", http.StatusBadRequest)
		return "", errors.New("malformed URL")
	}

	// Return name of country / isoCode
	return args[5], nil
}

/*
Get neighbour parameter from request

	w		- Responsewriter
	r		- Request

	return	- Bool which indicated whether user want data from neighbours of country returned
*/
func GetBoolParameterFromRequest(w http.ResponseWriter, r *http.Request, paramName string) (bool, error) {
	// Get CountryCodeOrName param
	paramString := (r.URL.Query()).Get(paramName)

	// Try to convert string to int
	paramBool, err := strconv.ParseBool(paramString)
	if err != nil && paramString != "" {
		http.Error(w, "Malformed URL, invalid "+paramName+" parameter set", http.StatusBadRequest)
		return false, err
	}

	// Return neighbours bool
	return paramBool, nil
}

/*
Get and decode webhook in json format into a webhook struct
*/
func GetWebhookFromRequest(w http.ResponseWriter, r *http.Request) (structs.Webhook, error) {
	// Decode JSON
	decoder := json.NewDecoder(r.Body)
	var webhook structs.Webhook
	if err := decoder.Decode(&webhook); err != nil {
		// Error for error in decoding
		log.Println(err.Error())
		return webhook, structs.NewError(err, http.StatusBadRequest, "Invalid request body for registration of webhook", "There was an error when decoding webhook from json.")
	}

	// Dont allow registration of webhook for country which does not exist in database
	if !db.DocumentInCollection(webhook.Country, constants.RENEWABLES_COLLECTION) {
		return webhook, structs.NewError(nil, http.StatusBadRequest, "Invalid country code for registration of webhook", "User entered a country code not in the database")
	}

	return webhook, nil
}

/*
Get webhookID from the requests url

	w		- Responsewriter
	r		- Request

	return 	- webhookID, or empty
*/
func GetWebhookIDFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	// Split path into args
	args := strings.Split(r.URL.Path, "/")

	// Check if URL is correctly formatted
	if len(args) != 5 && len(args) != 6 {
		http.Error(w, "Malformed URL, Expecting format "+constants.NOTIFICATION_PATH+"{webhookID}", http.StatusBadRequest)
		return "", errors.New("malformed URL")
	}

	// Return webhookID
	return args[4], nil

}

/*
Get webhookID from the requests url

	w		- Responsewriter
	r		- Request

	return 	- webhookID, or empty
*/
func GetWebhookIDOrNothingFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	// Split path into args
	args := strings.Split(r.URL.Path, "/")

	// Check if URL is correctly formatted
	if len(args) != 4 && len(args) != 5 && len(args) != 6 {
		http.Error(w, "Malformed URL, Expecting format "+constants.NOTIFICATION_PATH+"{webhookID?}", http.StatusBadRequest)
		return "", errors.New("malformed URL")
	}

	// If no webhookID was specified
	if len(args) == 4 {
		return "", nil
	}

	// Return webhookID
	return args[4], nil

}
