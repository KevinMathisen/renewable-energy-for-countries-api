package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/structs"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

const CURRENT_COUNTRY_CODE = "NOR"                   // ISO code for current country
const CURRENT_COUNTRY_NAME = "norway"                // Name for current country
const CURRENT_COUNTRY_PERCENTAGE float64 = 71.558365 // Percentage for current country
const CURRENT_COUNTRY_CODE_NEIGHBOURS = CURRENT_COUNTRY_CODE + "?neighbours=true"
const CURRENT_SORT_BY = "?sortByValue=true"
const EXPECTED_COUNTRIES = 72 //Amount of countries expected to get from the /country/ endpoint
const EXPECTED_NEIGHBOURS = 4 //Amount of neighbours expected to get from the country with country code CURRENT_COUNTRY_CODE

var CURRENT_NEIGHBOURS_CODES = [EXPECTED_NEIGHBOURS]string{"FIN", "NOR", "RUS", "SWE"}

/*
CURRENT COVERAGE:

/energy/v1/renewables/current/
	Tests year of first country
	Tests total amount of countries

/energy/v1/renewables/current/NOR
	Tests number of recieved countries
	Tests values of object

/energy/v1/renewables/current/norway
	Tests number of recieved countries
	Tests values of object

/energy/v1/renewables/current/NOR?neighbours=true
	Tests number of recieved countries
	Tests if the countries recieved is the same as CURRENT_NEIGHBOURS_CODES

/energy/v1/renewables/current/?sortByValue=true
	Tests the order of elements based on percentage relative to eachother
*/

/*
Gets data from the test URL and decodes into a slice of type CountryOutput, then returns this if
there are no errors
*/
func getCurrentData(client http.Client, url string) ([]structs.CountryOutput, error) {
	//Sends Get request
	res, err := client.Get(url)
	if err != nil {
		log.Println("Get request to URL failed:")
		return nil, err
	}

	var resObject []structs.CountryOutput
	//Recieves values, and decodes into slice
	err = json.NewDecoder(res.Body).Decode(&resObject)
	if err != nil {
		log.Println("Error during decoding:")
		return nil, err
	}

	return resObject, nil
}

/*
Runs http tests for all the different configuration types on the renewables current endpoint
*/
func TestHttpGetRenewablesCurrent(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllCachedRequestsFromFirestore()
	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	handleCurrentLogistics(t, currentAll)
	handleCurrentLogistics(t, currentCountryCode)
	handleCurrentLogistics(t, currentCountryName)
	handleCurrentLogistics(t, currentNeighbours)
	handleCurrentLogistics(t, currentSortBy)
}

/*
Handles opening and closing of server, alongside creating and closing client
Then calls given function for testing individual endpoints
*/
func handleCurrentLogistics(t *testing.T, f func(*testing.T, string, http.Client)) {
	handler := RootHandler(RenewablesCurrent)

	server := httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))
	defer server.Close()

	client := http.Client{}
	defer client.CloseIdleConnections()

	log.Println("URL: ", server.URL)

	url := server.URL + constants.RENEWABLES_CURRENT_PATH
	f(t, url, client)
}

// Runs tests for the .../renewables/current/ endpoint
func currentAll(t *testing.T, url string, client http.Client) {

	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the .../renewables/current/ endpoint
	res, err := getCurrentData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that the year is right
	if year := res[0].Year; year != strconv.Itoa(constants.LATEST_YEAR_DB) {
		t.Fatal("Year of recieved objects is wrong." +
			"\n\tExpected: " + strconv.Itoa(constants.LATEST_YEAR_DB) +
			"\n\tRecieved: " + year)
	}

	//Checks amount of countries recieved
	if length := len(res); length != EXPECTED_COUNTRIES {
		t.Fatal("Total amount of objects recieved is wrong." +
			"\n\tExpected: " + strconv.Itoa(EXPECTED_COUNTRIES) +
			"\n\tRecieved: " + strconv.Itoa(length))
	}
}

// Calls currentCountry(...) with a country code
func currentCountryCode(t *testing.T, url string, client http.Client) {
	fullUrl := url + CURRENT_COUNTRY_CODE
	currentCountry(t, fullUrl, client)
}

// Calls currentCountry(...) with a country name
func currentCountryName(t *testing.T, url string, client http.Client) {
	fullUrl := url + CURRENT_COUNTRY_NAME
	currentCountry(t, fullUrl, client)
}

// Runs tests for the .../renewables/current/{countryName}/{countryCode} endpoint
func currentCountry(t *testing.T, url string, client http.Client) {
	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the endpoint
	res, err := getCurrentData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that only one country was recieved
	if n := len(res); n != 1 {
		t.Fatal("Recieved more than one object." +
			"\n\tExpected: 1" +
			"\n\tRecieved: " + strconv.Itoa(n))
	}

	//Checks that the data in retrieved object is correct. Is case-insensitive on the country name.
	if res[0].IsoCode != CURRENT_COUNTRY_CODE || !strings.EqualFold(res[0].Name, CURRENT_COUNTRY_NAME) || res[0].Year != strconv.Itoa(constants.LATEST_YEAR_DB) || res[0].Percentage != CURRENT_COUNTRY_PERCENTAGE {
		t.Fatal("Wrong object recieved." +
			"\n\tExpected: " + CURRENT_COUNTRY_CODE + " - " + CURRENT_COUNTRY_NAME + " - " + strconv.Itoa(constants.LATEST_YEAR_DB) + " - " + strconv.FormatFloat(CURRENT_COUNTRY_PERCENTAGE, 'g', -1, 64) +
			"\n\tRecieved: " + res[0].IsoCode + " - " + res[0].Name + " - " + res[0].Year + " - " + strconv.FormatFloat(res[0].Percentage, 'g', -1, 64))
	}

}

// Runs tests for the .../renewables/current/{...}?neighbours=true endpoint
func currentNeighbours(t *testing.T, url string, client http.Client) {
	fullUrl := url + CURRENT_COUNTRY_CODE_NEIGHBOURS

	log.Println("Testing URL: \"" + fullUrl + "\"...")

	//Gets data from the endpoint
	res, err := getCurrentData(client, fullUrl)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks amount of countries recieved
	if length := len(res); length != EXPECTED_NEIGHBOURS {
		t.Fatal("Total amount of objects recieved is wrong." +
			"\n\tExpected: " + strconv.Itoa(EXPECTED_NEIGHBOURS) +
			"\n\tRecieved: " + strconv.Itoa(length))
	}

	//Checks if the countries recieved are the same as the neighbours of the tested country
	//This check is dependent on the order that the countries are in
	equal := true //Assume it is correct
	for i, v := range CURRENT_NEIGHBOURS_CODES {
		if res[i].IsoCode != v {
			equal = false
		}
	}

	if !equal {
		t.Fatal("The list of neighbours is not correct." +
			"\n\tExpected: " + CURRENT_NEIGHBOURS_CODES[0] + " - " + CURRENT_NEIGHBOURS_CODES[1] + " - " + CURRENT_NEIGHBOURS_CODES[2] + " - " + CURRENT_NEIGHBOURS_CODES[3] +
			"\n\tRecieved: " + res[0].IsoCode + " - " + res[1].IsoCode + " - " + res[2].IsoCode + " - " + res[3].IsoCode)
	}
}

// Runs tests for the .../renewables/current/?sortByValue=true endpoint
func currentSortBy(t *testing.T, url string, client http.Client) {
	fullUrl := url + CURRENT_SORT_BY

	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the endpoint
	res, err := getCurrentData(client, fullUrl)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that the percentage of the first country in returned slice is larger than that of the second,
	//and that the last is smaller than the second to last
	if res[0].Percentage < res[1].Percentage || res[len(res)-2].Percentage < res[len(res)-1].Percentage {
		t.Fatal("The order of the sorted values is wrong.")
	}
}
