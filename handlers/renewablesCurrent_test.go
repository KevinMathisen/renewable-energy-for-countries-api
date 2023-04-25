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
	"testing"
)

const CURRENT_COUNTRY_CODE = "NOR"
const CURRENT_COUNTRY_NAME = "norway"
const CURRENT_COUNTRY_CODE_NEIGHBOURS = CURRENT_COUNTRY_CODE + "?neighbours=true"
const EXPECTED_COUNTRIES = 72 //Amount of countries expected to get from the /country/ endpoint
const EXPECTED_NEIGHBOURS = 4 //Amount of neighbours expected to get from the country with country code CURRENT_COUNTRY_CODE

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

	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	handleCurrentLogistics(t, currentAll)
	handleCurrentLogistics(t, currentCountryCode)
	handleCurrentLogistics(t, currentCountryName)
	handleCurrentLogistics(t, currentNeighbours)
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

func currentAll(t *testing.T, url string, client http.Client) {

	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the /renewables/current/... endpoint
	res, err := getCurrentData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that the year is right
	if year := res[0].Year; year != strconv.Itoa(constants.LATEST_YEAR_DB) {
		t.Fatal("Year of recieved objects is wrong. Expected: " + strconv.Itoa(constants.LATEST_YEAR_DB) + ", recieved: " + year)
	}

	//Checks amount of countries recieved
	if length := len(res); length != EXPECTED_COUNTRIES {
		t.Fatal("Total amount of objects recieved is wrong. Expected: " + strconv.Itoa(EXPECTED_COUNTRIES) + ", recieved: " + strconv.Itoa(length))
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

func currentCountry(t *testing.T, url string, client http.Client) {
	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the /renewables/current/... endpoint
	res, err := getCurrentData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that only one country was recieved
	if n := len(res); n != 1 {
		t.Fatal("Recieved more than one object. Amount recieved: " + strconv.Itoa(n))
	}

	//Checks that the year is right
	if year := res[0].Year; year != strconv.Itoa(constants.LATEST_YEAR_DB) {
		t.Fatal("Year of recieved object is wrong. Expected: " + strconv.Itoa(constants.LATEST_YEAR_DB) + ", recieved: " + year)
	}
}

func currentNeighbours(t *testing.T, url string, client http.Client) {
	fullUrl := url + CURRENT_COUNTRY_CODE_NEIGHBOURS

	log.Println("Testing URL: \"" + fullUrl + "\"...")

	//Gets data from the /renewables/current/... endpoint
	res, err := getCurrentData(client, fullUrl)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks amount of countries recieved
	if length := len(res); length != EXPECTED_NEIGHBOURS {
		t.Fatal("Total amount of objects recieved is wrong. Expected: " + strconv.Itoa(EXPECTED_NEIGHBOURS) + ", recieved: " + strconv.Itoa(length))
	}
}
