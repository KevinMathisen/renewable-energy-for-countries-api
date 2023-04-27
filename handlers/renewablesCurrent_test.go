package handlers

import (
	htu "assignment2/http_test_utils"
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
TEST COVERAGE:

/energy/v1/renewables/current/NOR
	Tests number of recieved countries
	Tests all values of country

/energy/v1/renewables/current/norway
	Tests number of recieved countries
	Tests all values of country

/energy/v1/renewables/current/NOR?neighbours=true
	Tests number of recieved countries
	Tests if the countries recieved is the same as htu.NEIGHBOURS_CODES

/energy/v1/renewables/current/NOR?neighbours=true&sortByValue=true
	Tests if the countries recieved is the same as htu.SORTED_NEIGHBOURS_CODES

/energy/v1/renewables/current/
	Tests total amount of countries

/energy/v1/renewables/current/?sortByValue=true
	Tests the order of elements based on percentage relative to eachother
*/

/*
Handles opening and closing of server, alongside creating and closing client
Then calls given function for testing individual endpoints
*/
func handleCurrentLogistics(t *testing.T, f func(*testing.T, string, http.Client)) {
	//Creates instance of RenewablesCurrent handler
	handler := RootHandler(RenewablesCurrent)

	//Runs handler instance as server
	server := httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))
	defer server.Close()

	//Creates client to speak with server
	client := http.Client{}
	defer client.CloseIdleConnections()

	log.Println("URL: ", server.URL)

	url := server.URL + constants.RENEWABLES_CURRENT_PATH

	f(t, url, client)
}

/*
Runs http tests for all the different configuration types on the renewables current endpoint
*/
func TestHttpGetRenewablesCurrent(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllDocumentsInCollectionFromFirestore(constants.CACHE_COLLECTION)
	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	handleCurrentLogistics(t, currentCountryByCode)
	handleCurrentLogistics(t, currentCountryByName)
	handleCurrentLogistics(t, currentNeighbours)
	handleCurrentLogistics(t, currentNeighboursSortBy)
	handleCurrentLogistics(t, currentAll)
	handleCurrentLogistics(t, currentAllSortBy)
}

//------------------------------ SINGLE COUNTRY TESTS ------------------------------

// Calls currentCountry(...) with a country code
func currentCountryByCode(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE
	currentCountry(t, url, client)
}

// Calls currentCountry(...) with a country name
func currentCountryByName(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_NAME
	currentCountry(t, url, client)
}

// Runs tests for the .../renewables/current/{<NOR>/<norway>} endpoint
func currentCountry(t *testing.T, url string, client http.Client) {

	//Gets data from the endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that only one country was recieved
	if err2 := htu.TestLen(res, 1); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the data in recieved object is correct. Is case-insensitive on the country name.
	if err2 := htu.TestValues(res[0], htu.COUNTRY_CODE, htu.COUNTRY_NAME, constants.LATEST_YEAR_DB, htu.COUNTRY_LATEST_PERCENTAGE); err2 != "" {
		t.Fatal(err2)
	}
}

//------------------------------ NEIGHBOUR COUNTRY TESTS ------------------------------

// Runs tests for the .../renewables/current/NOR?neighbours=true endpoint
func currentNeighbours(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.NEIGHBOURS

	//Gets data from the endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks amount of countries recieved
	if err2 := htu.TestLen(res, htu.EXPECTED_NEIGHBOURS); err2 != "" {
		t.Fatal(err2)
	}

	//Checks if the countries recieved are the same as in htu.NEIGHBOURS_CODES
	if err2 := htu.TestSortedCodeList(res, htu.NEIGHBOURS_CODES); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/current/NOR?neighbours=true&sortByValue=true endpoint
func currentNeighboursSortBy(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.NEIGHBOURS + htu.AND + htu.SORT_BY

	//Gets data from the endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks if the countries recieved are sorted
	if err2 := htu.TestSortedCodeList(res, htu.SORTED_NEIGHBOURS_CODES); err2 != "" {
		t.Fatal(err2)
	}
}

//------------------------------ ALL COUNTRIES TESTS ------------------------------

// Runs tests for the .../renewables/current/ endpoint
func currentAll(t *testing.T, url string, client http.Client) {

	//Gets data from the .../renewables/current/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks amount of countries recieved
	if err2 := htu.TestLen(res, htu.CURRENT_COUNTRIES); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/current/?sortByValue=true endpoint
func currentAllSortBy(t *testing.T, url string, client http.Client) {
	url = url + htu.PARAM + htu.SORT_BY

	//Gets data from the endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that the percentage of the first country in returned slice is larger than that of the second,
	//and that the last is smaller than the second to last
	if err2 := htu.TestSortedPercentage(res[0].Percentage, res[1].Percentage, res[len(res)-2].Percentage, res[len(res)-1].Percentage); err2 != "" {
		t.Fatal(err2)
	}
}
