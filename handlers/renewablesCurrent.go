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
