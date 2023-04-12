package handlers

import (
	"assignment2/utils/structs"
	"errors"
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
	hit := checkCache(w, r)
	if hit {
		return
	}

	// Get the countries we are interested in finding, or empty if everyone
	countries, err := getCountriesToQuery(w, r)
	if err != nil {
		return
	}

	// Get current percentage of renewables for countries specified as a list of countryoutput structs
	response, err = getCurrentRenewablesForCountries(w, countries)
	if err != nil {
		return
	}

	// Respond with list of countryoutput struct encoded as json to user
	respondToGetRequestWithJSON(w, response)
}

/*
Get country code or name, and neighbours parameter from request, then returns appropiate list of countries from these

	w	- Responsewriter
	r	- Request

	return	- Either empty list of no country specified, or one country, or country and it's neighbours
*/
func getCountriesToQuery(w http.ResponseWriter, r *http.Request) ([]string, error) {
	var countries []string

	// Get country code or name from request
	countryCodeOrName, err := getCountryCodeOrNameFromRequest(r)
	if err != nil {
		return nil, err
	}

	// Get neigbour bool from request if it is specified
	neigbours, err := getNeigboursParameterFromRequest(r)
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

	} else if isoCodeInDB(countryCodeOrName) {
		// Else if the user specified ISO code and it exists in the database, add the code the list of countries
		countries = append(countries, countryCodeOrName)
	}

	// If the user specified the neighbour parameter
	if neigbours {
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
Get renewables data for the current year from specified countires or all countries

	w			- Responsewriter for sending error messages
	countires	- Either a list of countries we want to get data from, or an empty list if we want all

	return		- Returns a list of countryouput structs which can will be sent as json in the response
*/
func getCurrentRenewablesForCountries(w http.ResponseWriter, countries []string) ([]structs.CountryOutput, error) {
	var renewablesOutput []structs.CountryOutput

	// Get current year
	// TODO: Get current Year
	currentYear := 2021

	// If the users specified countries, get renewables data from them in the current year
	if len(countries) != 0 {
		renewablesOutput, err := getRenewablesForCountriesByYears(w, countries, currentYear, currentYear)
		if err != nil {
			return renewablesOutput, err
		}
	} else {
		// If the user did not specify countires, we get renewables data from all countires in the current year
		renewablesOutput, err := getRenewablesForAllCountriesByYears(w, currentYear, currentYear)
		if err != nil {
			return renewablesOutput, err
		}
	}

	return renewablesOutput, nil
}
