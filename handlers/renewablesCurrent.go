package handlers

import (
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
	response = getCurrentRenewablesForCountries(w, countries)

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

		// Else if the user specified ISO code, add the code the list of countries
	} else {
		countries = append(countries, countryCodeOrName)
	}

	// If the user specified the neighbour parameter
	if neigbours {
		// TODO: Get neighbour ISO code with Restcountries API, and append to countries
	}

	return countries, nil

}
