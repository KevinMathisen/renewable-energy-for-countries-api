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

// Values for building URLs
const BEGIN_YEAR = "1990"
const END_YEAR = "2010"
const HISTORY_COUNTRY_CODE = "NOR"
const HISTORY_COUNTRY_NAME = "norway"
const HISTORY_BEGIN = "begin=" + BEGIN_YEAR
const HISTORY_END = "end=" + END_YEAR
const HISTORY_SORT_BY = "sortByValue=true"
const HISTORY_NEIGHBOURS = "neighbours=true"
const HISTORY_MEAN = "mean=true"
const HISTORY_PARAM = "?"
const HISTORY_AND = "&"

// Values for checking
var INT_BEGIN_YEAR, _ = strconv.Atoi(BEGIN_YEAR)
var INT_END_YEAR, _ = strconv.Atoi(END_YEAR)

const HISTORY_COUNTRY_OLDEST_PERCENTAGE = 67.87996  //Oldest percentage for Norway
const HISTORY_COUNTRY_LATEST_PERCENTAGE = 71.558365 //Latest percentage for Norway
const HISTORY_COUNTRY_EXPECTED_INSTANCES = 57       //Amount of instances of Norway in the dataset

const HISTORY_COUNTRY_BEGIN_PERCENTAGE = 72.44774 //Percentage for Norway in year BEGIN_YEAR
const HISTORY_COUNTRY_END_PERCENTAGE = 65.47019   //Percentage for Norway in year END_YEAR
const HISTORY_COUNTRY_BEGIN_END_INSTANCES = 21    //Amount of instances of Norway in the dataset between BEGIN_YEAR and END_YEAR

const HISTORY_COUNTRY_BEGIN_INSTANCES = 32 //Amount of instances of Norway in the dataset between BEGIN_YEAR and the end of the dataset

const HISTORY_COUNTRY_END_INSTANCES = 46 //Amount of instances of Norway in the dataset between the start of the dataset and END_YEAR

const HISTORY_COUNTRY_BEGIN_END_SORT_FIRST = 1990                //The year of the first object after sort
const HISTORY_COUNTRY_BEGIN_END_SORT_FIRST_PERCENTAGE = 72.44774 //The percentage of the first object after sort
const HISTORY_COUNTRY_BEGIN_END_SORT_LAST = 2003                 //The year of the last object after sort
const HISTORY_COUNTRY_BEGIN_END_SORT_LAST_PERCENTAGE = 63.816036 //The percentage of the last object after sort

const HISTORY_COUNTRY_MEAN = 68.01918892982457

const HISTORY_COUNTRY_BEGIN_END_MEAN = 68.63185428571428

/*
TEST COVERAGE:

/energy/v1/renewables/history/NOR
	Checks number of recieved objects
	Tests all values of the first instance
	Tests all values of the last instance

/energy/v1/renewables/history/norway
	Checks number of recieved objects
	Tests all values of the first instance
	Tests all values of the last instance

/energy/v1/renewables/history/NOR?begin=1990&end=2010
	Checks number of recieved objects
	Tests all values of the first instance
	Tests all values of the last instance

/energy/v1/renewables/history/NOR?begin=1990
	Checks number of recieved objects
	Tests all values of the first instance
	Tests all values of the last instance

/energy/v1/renewables/history/NOR?end=2010
	Checks number of recieved objects
	Tests all values of the first instance
	Tests all values of the last instance

/energy/v1/renewables/history/NOR?begin=1990&end=2010&sortByValue=true
	Checks number of recieved objects
	Tests all values of the first instance
	Tests all values of the last instance

/energy/v1/renewables/history/NOR?mean=true
	Checks number of recieved objects
	Checks whether recieved object has year value or not
	Tests percentage value of recieved object
*/

/*
Gets data from the test URL and decodes into a slice of type CountryOutput, then returns this if
there are no errors
*/
func getHistoryData(client http.Client, url string) ([]structs.CountryOutput, error) {
	log.Println("Testing URL: \"" + url + "\"...")

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
Tests the length of given slice, then the values of the first object, then the values of the last
*/
func testLenLastFirst(s []structs.CountryOutput, objects int, fISO string, fName string, fYear int, fPer float64, lISO string, lName string, lYear int, lPer float64) string {
	//Checks that all instances of the country is recieved
	if n := len(s); n != objects {
		return ("Recieved more than one object." +
			"\n\tExpected: " + strconv.Itoa(objects) +
			"\n\tRecieved: " + strconv.Itoa(n))
	}

	//Checks that the data in the first recieved object is correct. Is case-insensitive on the country name.
	if s[0].IsoCode != fISO || !strings.EqualFold(s[0].Name, fName) || s[0].Year != strconv.Itoa(fYear) || s[0].Percentage != fPer {
		return ("Wrong object recieved." +
			"\n\tExpected: " + fISO + " - " + fName + " - " + strconv.Itoa(fYear) + " - " + strconv.FormatFloat(fPer, 'g', -1, 64) +
			"\n\tRecieved: " + s[0].IsoCode + " - " + s[0].Name + " - " + s[0].Year + " - " + strconv.FormatFloat(s[0].Percentage, 'g', -1, 64))
	}

	last := objects - 1
	//Checks that the data in the last recieved object is correct. Is case-insensitive on the country name.
	if s[last].IsoCode != lISO || !strings.EqualFold(s[last].Name, lName) || s[last].Year != strconv.Itoa(lYear) || s[last].Percentage != lPer {
		return ("Wrong object recieved." +
			"\n\tExpected: " + lISO + " - " + lName + " - " + strconv.Itoa(lYear) + " - " + strconv.FormatFloat(lPer, 'g', -1, 64) +
			"\n\tRecieved: " + s[last].IsoCode + " - " + s[last].Name + " - " + s[last].Year + " - " + strconv.FormatFloat(s[last].Percentage, 'g', -1, 64))
	}

	return ""
}

func TestHttpGetRenewablesHistory(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllCachedRequestsFromFirestore()
	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	//Country
	//handleHistoryLogistics(t, historyCountryByCode)
	//handleHistoryLogistics(t, historyCountryByName)
	//handleHistoryLogistics(t, historyCountryBeginEnd)
	//handleHistoryLogistics(t, historyCountryBegin)
	//handleHistoryLogistics(t, historyCountryEnd)
	//handleHistoryLogistics(t, historyCountryBeginEndSort)
	//handleHistoryLogistics(t, historyCountryMean)
	//handleHistoryLogistics(t, historyCountryBeginEndMean)
	//Neighbour
	handleHistoryLogistics(t, historyNeighbours)
	return
	handleHistoryLogistics(t, historyNeighboursBeginEnd)
	handleHistoryLogistics(t, historyNeighboursMean)
	handleHistoryLogistics(t, historyNeighboursSort)
	//All
	handleHistoryLogistics(t, historyAll)
	handleHistoryLogistics(t, historyAllSort)
	handleHistoryLogistics(t, historyAllMean)
	handleHistoryLogistics(t, historyAllSortMean)
}

/*
Handles opening and closing of server, alongside creating and closing client
Then calls given function for testing individual endpoints
*/
func handleHistoryLogistics(t *testing.T, f func(*testing.T, string, http.Client)) {
	handler := RootHandler(RenewablesHistory)

	server := httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))
	defer server.Close()

	client := http.Client{}
	defer client.CloseIdleConnections()

	log.Println("URL: ", server.URL)

	url := server.URL + constants.RENEWABLES_HISTORY_PATH

	f(t, url, client)
}

//------------------------------ SINGLE COUNTRY TESTS ------------------------------

// Calls historyCountry(...) with a country code
func historyCountryByCode(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE
	historyCountry(t, url, client)
}

// Calls historyCountry(...) with a country code
func historyCountryByName(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_NAME
	historyCountry(t, url, client)
}

// Runs tests for the .../renewables/history/{<NOR>/<norway>} endpoint
func historyCountry(t *testing.T, url string, client http.Client) {

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	err2 := testLenLastFirst(res, HISTORY_COUNTRY_EXPECTED_INSTANCES, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, constants.OLDEST_YEAR_DB, HISTORY_COUNTRY_OLDEST_PERCENTAGE, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, constants.LATEST_YEAR_DB, HISTORY_COUNTRY_LATEST_PERCENTAGE)
	if err2 != "" {
		t.Fatal(err2)
	}

}

// Runs tests for the .../renewables/history/NOR?begin={BEGIN_YEAR}&end={END_YEAR} endpoint
func historyCountryBeginEnd(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_BEGIN + HISTORY_AND + HISTORY_END

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	err2 := testLenLastFirst(res, HISTORY_COUNTRY_BEGIN_END_INSTANCES, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, INT_BEGIN_YEAR, HISTORY_COUNTRY_BEGIN_PERCENTAGE, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, INT_END_YEAR, HISTORY_COUNTRY_END_PERCENTAGE)
	if err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?begin={BEGIN_YEAR} endpoint
func historyCountryBegin(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_BEGIN

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	err2 := testLenLastFirst(res, HISTORY_COUNTRY_BEGIN_INSTANCES, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, INT_BEGIN_YEAR, HISTORY_COUNTRY_BEGIN_PERCENTAGE, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, constants.LATEST_YEAR_DB, HISTORY_COUNTRY_LATEST_PERCENTAGE)
	if err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?end={END_YEAR} endpoint
func historyCountryEnd(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_END

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	err2 := testLenLastFirst(res, HISTORY_COUNTRY_END_INSTANCES, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, constants.OLDEST_YEAR_DB, HISTORY_COUNTRY_OLDEST_PERCENTAGE, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, INT_END_YEAR, HISTORY_COUNTRY_END_PERCENTAGE)
	if err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?begin={BEGIN_YEAR}&end={END_YEAR}&sortByValue=true endpoint
func historyCountryBeginEndSort(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_BEGIN + HISTORY_AND + HISTORY_END + HISTORY_AND + HISTORY_SORT_BY

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	err2 := testLenLastFirst(res, HISTORY_COUNTRY_BEGIN_END_INSTANCES, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, HISTORY_COUNTRY_BEGIN_END_SORT_FIRST, HISTORY_COUNTRY_BEGIN_END_SORT_FIRST_PERCENTAGE, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, HISTORY_COUNTRY_BEGIN_END_SORT_LAST, HISTORY_COUNTRY_BEGIN_END_SORT_LAST_PERCENTAGE)
	if err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?mean=true endpoint
func historyCountryMean(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_MEAN

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks if there is only one object recieved
	if n := len(res); n != 1 {
		t.Fatal("Too many objects returned." +
			"\n\tExpected: 1" +
			"\n\tRecieved: " + strconv.Itoa(n))
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//Tests the percentage
	if res[0].Percentage != HISTORY_COUNTRY_MEAN {
		t.Fatal("Mean percentage for country is not correct." +
			"\n\tExpected: " + strconv.FormatFloat(HISTORY_COUNTRY_MEAN, 'g', -1, 64) +
			"\n\tRecieved: " + strconv.FormatFloat(res[0].Percentage, 'g', -1, 64))
	}
}

// Runs tests for the .../renewables/history/NOR?begin={BEGIN_YEAR}&end={END_YEAR}&mean=true endpoint
func historyCountryBeginEndMean(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_BEGIN + HISTORY_AND + HISTORY_END + HISTORY_AND + HISTORY_MEAN

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks if there is only one object recieved
	if n := len(res); n != 1 {
		t.Fatal("Too many objects returned." +
			"\n\tExpected: 1" +
			"\n\tRecieved: " + strconv.Itoa(n))
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//Tests the percentage
	if res[0].Percentage != HISTORY_COUNTRY_BEGIN_END_MEAN {
		t.Fatal("Mean percentage for country is not correct." +
			"\n\tExpected: " + strconv.FormatFloat(HISTORY_COUNTRY_BEGIN_END_MEAN, 'g', -1, 64) +
			"\n\tRecieved: " + strconv.FormatFloat(res[0].Percentage, 'g', -1, 64))
	}
}

// ------------------------------ NEIGHBOUR COUNTRY TESTS ------------------------------

// Runs tests for the .../renewables/history/NOR?neighbours=true endpoint
func historyNeighbours(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_NEIGHBOURS

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

// Runs tests for the .../renewables/history/NOR?begin={BEGIN_YEAR}&end={END_YEAR} endpoint
func historyNeighboursBeginEnd(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_NEIGHBOURS + HISTORY_AND + HISTORY_BEGIN + HISTORY_AND + HISTORY_END

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

// Runs tests for the .../renewables/history/NOR?sortByValue=true endpoint
func historyNeighboursSort(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_NEIGHBOURS + HISTORY_AND + HISTORY_SORT_BY

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

// Runs tests for the .../renewables/history/NOR?mean=true endpoint
func historyNeighboursMean(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_NEIGHBOURS + HISTORY_AND + HISTORY_MEAN

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

//------------------------------ ALL COUNTRIES TESTS ------------------------------

// Runs tests for the .../renewables/history/ endpoint
func historyAll(t *testing.T, url string, client http.Client) {

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

// Runs tests for the .../renewables/history/?sortByValue=true endpoint
func historyAllSort(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_PARAM + HISTORY_SORT_BY

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

// Runs tests for the .../renewables/history/?mean=true endpoint
func historyAllMean(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_PARAM + HISTORY_MEAN

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

// Runs tests for the .../renewables/history/?sortByValue=true&mean=true endpoint
func historyAllSortMean(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_PARAM + HISTORY_SORT_BY + HISTORY_AND + HISTORY_MEAN

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}
