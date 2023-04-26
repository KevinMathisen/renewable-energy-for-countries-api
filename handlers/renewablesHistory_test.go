package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/structs"
	"encoding/json"
	"log"
	"math"
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
var INT_BEGIN_YEAR, _ = strconv.Atoi(BEGIN_YEAR) //Int value of BEGIN_YEAR
var INT_END_YEAR, _ = strconv.Atoi(END_YEAR)     //Int value of END_YEAR

const HISTORY_COUNTRY_OLDEST_PERCENTAGE = 67.87996  //Oldest percentage for Norway
const HISTORY_COUNTRY_LATEST_PERCENTAGE = 71.558365 //Latest percentage for Norway
const HISTORY_COUNTRY_EXPECTED_ENTRIES = 57         //Amount of entries Norway has in the dataset

const HISTORY_COUNTRY_BEGIN_PERCENTAGE = 72.44774 //Percentage for Norway in year BEGIN_YEAR
const HISTORY_COUNTRY_END_PERCENTAGE = 65.47019   //Percentage for Norway in year END_YEAR
const HISTORY_COUNTRY_BEGIN_END_ENTRIES = 21      //Amount of entries Norway has in the dataset between BEGIN_YEAR and END_YEAR

const HISTORY_COUNTRY_BEGIN_ENTRIES = 32 //Amount of entries Norway has in the dataset between BEGIN_YEAR and the end of the dataset

const HISTORY_COUNTRY_END_ENTRIES = 46 //Amount of entries Norway has in the dataset between the start of the dataset and END_YEAR

const HISTORY_COUNTRY_BEGIN_END_SORT_FIRST = 1990                //The year of the first object after sort
const HISTORY_COUNTRY_BEGIN_END_SORT_FIRST_PERCENTAGE = 72.44774 //The percentage of the first object after sort
const HISTORY_COUNTRY_BEGIN_END_SORT_LAST = 2003                 //The year of the last object after sort
const HISTORY_COUNTRY_BEGIN_END_SORT_LAST_PERCENTAGE = 63.816036 //The percentage of the last object after sort

const HISTORY_COUNTRY_MEAN = 68.01918892982457 //Mean percentage for Norway

const HISTORY_COUNTRY_BEGIN_END_MEAN = 68.63185428571428 //Mean percentage for Norway between BEGIN_YEAR and END_YEAR

const HISTORY_NEIGHBOUR_ENTRIES_AMOUNT = 208 //Amount of objects returned when calling for the neighbours of Norway

const HISTORY_NEIGHBOUR_BEGIN_END_AMOUNT = 84 //Amount of objects returned when calling for the nieghbours of Norway between BEGIN_YEAR and END_YEAR

const HISTORY_NEIGHBOURS_SORT_LAST_CODE = "RUS"          //The ISO code of the last country in the list when sorted by oercentage
const HISTORY_NEIGHBOURS_SORT_LAST_NAME = "russia"       //The name of the last country in the list when sorted by oercentage
const HISTORY_NEIGHBOURS_SORT_LAST_PERCENTAGE = 4.605263 //The percentage of the last country in the list when sorted by oercentage
const HISTORY_NEIGHBOURS_SORT_LAST_YEAR = 1989           //The year of the last country in the list when sorted by oercentage

const HISTORY_NEIGHBOURS_AMOUNT = 4 //The amount of neighbours Norway  has

const HISTORY_NEIGHBOURS_SORT_MEAN_LAST = 6.004957597297297 //Percentage of the last country recieved from the sorted mean of the neighbours

const HISTORY_ALL_COUNTRIES = 79 //All different countries in the dataset

const HISTORY_ALL_SORT_LAST_CODE = "SAU"                       //The ISO code of the last country in the list when sorting all countries by percentage
const HISTORY_ALL_SORT_LAST_NAME = "Saudi Arabia"              //The name of the last country in the list when sorting all countries by percentage
const HISTORY_ALL_SORT_LAST_PERCENTAGE = 0.0013665377020357142 //The percentage of the last country in the list when sorting all countries by percentage

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

/energy/v1/renewables/history/NOR?neighbours=true&begin=1990&end=2010
	Cheacks amount of returned objects

/energy/v1/renewables/history/NOR?neighbours=true&mean=true
	Cheacks amount of returned objects
	Checks whether recieved object has year value or not
	Tests the percentage of the second recieved object (Norway)

/energy/v1/renewables/history/NOR?neighbours=true&sortByValue=true
	Cheacks amount of returned objects
	Tests all values of the first instance
	Tests all values of the last instance

/energy/v1/renewables/history/NOR?neighbours=true&mean=true&sortByValue=true
	Checks whether recieved object has year value or not
	Cheacks amount of returned objects
	Tests all values of the first instance
	Tests all values of the last instance

/energy/v1/renewables/history/
	Cheacks amount of returned objects

/energy/v1/renewables/history/?sortByValue=true
	Checks whether recieved object has year value or not
	Cheacks amount of returned objects
	Tests all values of the first instance
	Tests all values of the last instance
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
Tests to see if the float check is equal to float mark to 13 decimal places
This is to avoid floating point errors that seem to appear around the 13th decimal place
*/
func testPercentage(check, mark float64) string {
	checkInt := int(check * math.Pow(10, 12))
	markInt := int(mark * math.Pow(10, 12))

	if checkInt != markInt {
		return ("Mean percentage for country is not correct." +
			"\n\tExpected: " + strconv.FormatFloat(mark, 'g', -1, 64) +
			"\n\tRecieved: " + strconv.FormatFloat(check, 'g', -1, 64))
	}
	return ""
}

/*
Tests the length of given slice against the given int.
*/
func testLen(s []structs.CountryOutput, expected int) string {
	if n := len(s); n != expected {
		return ("Wrong amount of objects returned." +
			"\n\tExpected: " + strconv.Itoa(expected) +
			"\n\tRecieved: " + strconv.Itoa(n))
	}
	return ""
}

/*
Tests the length of given slice, then tests the values of the first object, then the values of the last

	s:			the slice of objects to be operated on
	objects:	the amount of objects to be expected in s
	fISO:		the expected code of the fist object in s
	fName:		the expected name of the fist object in s
	fYear:		the expected year of the fist object in s
	fPer:		the expected percentage of the fist object in s
	lISO:		the expected code of the last object in s
	lName:		the expected name of the last object in s
	lYear:		the expected year of the last object in s
	lPer:		the expected percentage of the last object in s

	return:		a message explaining whats wrong if the test failed, empty string if test succeeded
*/
func testLenLastFirst(s []structs.CountryOutput, objects int, fISO string, fName string, fYear int, fPer float64, lISO string, lName string, lYear int, lPer float64) string {
	//Checks that all instances of the country is recieved
	if err := testLen(s, objects); err != "" {
		return err
	}

	//Checks that the data in the first recieved object is correct. Is case-insensitive on the country name.
	if s[0].IsoCode != fISO || !strings.EqualFold(s[0].Name, fName) || s[0].Year != strconv.Itoa(fYear) || testPercentage(s[0].Percentage, fPer) != "" {
		return ("Wrong object recieved." +
			"\n\tExpected: " + fISO + " - " + fName + " - " + strconv.Itoa(fYear) + " - " + strconv.FormatFloat(fPer, 'g', -1, 64) +
			"\n\tRecieved: " + s[0].IsoCode + " - " + s[0].Name + " - " + s[0].Year + " - " + strconv.FormatFloat(s[0].Percentage, 'g', -1, 64))
	}

	last := objects - 1
	//Checks that the data in the last recieved object is correct. Is case-insensitive on the country name.
	if s[last].IsoCode != lISO || !strings.EqualFold(s[last].Name, lName) || s[last].Year != strconv.Itoa(lYear) || testPercentage(s[last].Percentage, lPer) != "" {
		return ("Wrong object recieved." +
			"\n\tExpected: " + lISO + " - " + lName + " - " + strconv.Itoa(lYear) + " - " + strconv.FormatFloat(lPer, 'g', -1, 64) +
			"\n\tRecieved: " + s[last].IsoCode + " - " + s[last].Name + " - " + s[last].Year + " - " + strconv.FormatFloat(s[last].Percentage, 'g', -1, 64))
	}

	return ""
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

func TestHttpGetRenewablesHistory(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllDocumentsInCollectionFromFirestore(constants.CACHE_COLLECTION)
	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	//Country
	handleHistoryLogistics(t, historyCountryByCode)
	handleHistoryLogistics(t, historyCountryByName)
	handleHistoryLogistics(t, historyCountryBeginEnd)
	handleHistoryLogistics(t, historyCountryBegin)
	handleHistoryLogistics(t, historyCountryEnd)
	handleHistoryLogistics(t, historyCountryBeginEndSort)
	handleHistoryLogistics(t, historyCountryMean)
	handleHistoryLogistics(t, historyCountryBeginEndMean)
	//Neighbour
	handleHistoryLogistics(t, historyNeighbours)
	handleHistoryLogistics(t, historyNeighboursBeginEnd)
	handleHistoryLogistics(t, historyNeighboursMean)
	handleHistoryLogistics(t, historyNeighboursSort)
	handleHistoryLogistics(t, historyNeighboursMeanSort)
	//All
	handleHistoryLogistics(t, historyAll)
	handleHistoryLogistics(t, historyAllSort)
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

	if err2 := testLenLastFirst(res, HISTORY_COUNTRY_EXPECTED_ENTRIES, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, constants.OLDEST_YEAR_DB, HISTORY_COUNTRY_OLDEST_PERCENTAGE, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, constants.LATEST_YEAR_DB, HISTORY_COUNTRY_LATEST_PERCENTAGE); err2 != "" {
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

	if err2 := testLenLastFirst(res, HISTORY_COUNTRY_BEGIN_END_ENTRIES, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, INT_BEGIN_YEAR, HISTORY_COUNTRY_BEGIN_PERCENTAGE, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, INT_END_YEAR, HISTORY_COUNTRY_END_PERCENTAGE); err2 != "" {
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

	if err2 := testLenLastFirst(res, HISTORY_COUNTRY_BEGIN_ENTRIES, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, INT_BEGIN_YEAR, HISTORY_COUNTRY_BEGIN_PERCENTAGE, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, constants.LATEST_YEAR_DB, HISTORY_COUNTRY_LATEST_PERCENTAGE); err2 != "" {
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

	if err2 := testLenLastFirst(res, HISTORY_COUNTRY_END_ENTRIES, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, constants.OLDEST_YEAR_DB, HISTORY_COUNTRY_OLDEST_PERCENTAGE, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, INT_END_YEAR, HISTORY_COUNTRY_END_PERCENTAGE); err2 != "" {
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

	if err2 := testLenLastFirst(res, HISTORY_COUNTRY_BEGIN_END_ENTRIES, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, HISTORY_COUNTRY_BEGIN_END_SORT_FIRST, HISTORY_COUNTRY_BEGIN_END_SORT_FIRST_PERCENTAGE, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, HISTORY_COUNTRY_BEGIN_END_SORT_LAST, HISTORY_COUNTRY_BEGIN_END_SORT_LAST_PERCENTAGE); err2 != "" {
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
	if err2 := testLen(res, 1); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//Tests the percentage
	if err2 := testPercentage(res[0].Percentage, HISTORY_COUNTRY_MEAN); err2 != "" {
		t.Fatal(err2)
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
	if err2 := testLen(res, 1); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//Tests the percentage
	if err2 := testPercentage(res[0].Percentage, HISTORY_COUNTRY_BEGIN_END_MEAN); err2 != "" {
		t.Fatal(err2)
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

	if err2 := testLen(res, HISTORY_NEIGHBOUR_ENTRIES_AMOUNT); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?neighbours=true&begin={BEGIN_YEAR}&end={END_YEAR} endpoint
func historyNeighboursBeginEnd(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_NEIGHBOURS + HISTORY_AND + HISTORY_BEGIN + HISTORY_AND + HISTORY_END

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	if err2 := testLen(res, HISTORY_NEIGHBOUR_BEGIN_END_AMOUNT); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?neighbours=true&sortByValue=true endpoint
func historyNeighboursSort(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_NEIGHBOURS + HISTORY_AND + HISTORY_SORT_BY

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Tests the amount of recieved objects, then tests the values of the first and then the last object
	if err2 := testLenLastFirst(res, HISTORY_NEIGHBOUR_ENTRIES_AMOUNT, HISTORY_COUNTRY_CODE, HISTORY_COUNTRY_NAME, INT_BEGIN_YEAR, HISTORY_COUNTRY_BEGIN_END_SORT_FIRST_PERCENTAGE, HISTORY_NEIGHBOURS_SORT_LAST_CODE, HISTORY_NEIGHBOURS_SORT_LAST_NAME, HISTORY_NEIGHBOURS_SORT_LAST_YEAR, HISTORY_NEIGHBOURS_SORT_LAST_PERCENTAGE); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?neighbours=true&mean=true endpoint
func historyNeighboursMean(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_NEIGHBOURS + HISTORY_AND + HISTORY_MEAN

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that there are the expected amount of neighbours
	if err2 := testLen(res, HISTORY_NEIGHBOURS_AMOUNT); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//Checks that the percentage of Norway is correct. Checks that the order is correct as side-effect
	if err2 := testPercentage(res[1].Percentage, HISTORY_COUNTRY_MEAN); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?neighbours=true&mean=true&sortByValue=true endpoint
func historyNeighboursMeanSort(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_NEIGHBOURS + HISTORY_AND + HISTORY_MEAN + HISTORY_AND + HISTORY_SORT_BY

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//The code in testLenFirstLast but refactored to not include country year:
	//Tests the amount of recieved objects, then tests the values of the first and then the last object
	if err2 := testLen(res, HISTORY_NEIGHBOURS_AMOUNT); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the data in the first recieved object is correct. Is case-insensitive on the country name.
	if res[0].IsoCode != HISTORY_COUNTRY_CODE || !strings.EqualFold(res[0].Name, HISTORY_COUNTRY_NAME) || testPercentage(res[0].Percentage, HISTORY_COUNTRY_MEAN) != "" {
		t.Fatal("Wrong object recieved." +
			"\n\tExpected: " + HISTORY_COUNTRY_CODE + " - " + HISTORY_COUNTRY_NAME + " - " + strconv.FormatFloat(HISTORY_COUNTRY_MEAN, 'g', -1, 64) +
			"\n\tRecieved: " + res[0].IsoCode + " - " + res[0].Name + strconv.FormatFloat(res[0].Percentage, 'g', -1, 64))
	}

	last := len(res) - 1
	//Checks that the data in the last recieved object is correct. Is case-insensitive on the country name.
	if res[last].IsoCode != HISTORY_NEIGHBOURS_SORT_LAST_CODE || !strings.EqualFold(res[last].Name, HISTORY_NEIGHBOURS_SORT_LAST_NAME) || testPercentage(res[last].Percentage, HISTORY_NEIGHBOURS_SORT_MEAN_LAST) != "" {
		t.Fatal("Wrong object recieved." +
			"\n\tExpected: " + HISTORY_NEIGHBOURS_SORT_LAST_CODE + " - " + HISTORY_NEIGHBOURS_SORT_LAST_NAME + " - " + strconv.FormatFloat(HISTORY_NEIGHBOURS_SORT_MEAN_LAST, 'g', -1, 64) +
			"\n\tRecieved: " + res[last].IsoCode + " - " + res[last].Name + " - " + strconv.FormatFloat(res[last].Percentage, 'g', -1, 64))
	}
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

	//Checks the amount of recieved countries
	if err2 := testLen(res, HISTORY_ALL_COUNTRIES); err2 != "" {
		t.Fatal(err2)
	}
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

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//The code in testLenFirstLast but refactored to not include country year:
	//Tests the amount of recieved objects, then tests the values of the first and then the last object
	if err2 := testLen(res, HISTORY_ALL_COUNTRIES); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the data in the first recieved object is correct. Is case-insensitive on the country name.
	if res[0].IsoCode != HISTORY_COUNTRY_CODE || !strings.EqualFold(res[0].Name, HISTORY_COUNTRY_NAME) || testPercentage(res[0].Percentage, HISTORY_COUNTRY_MEAN) != "" {
		t.Fatal("Wrong object recieved." +
			"\n\tExpected: " + HISTORY_COUNTRY_CODE + " - " + HISTORY_COUNTRY_NAME + " - " + strconv.FormatFloat(HISTORY_COUNTRY_MEAN, 'g', -1, 64) +
			"\n\tRecieved: " + res[0].IsoCode + " - " + res[0].Name + strconv.FormatFloat(res[0].Percentage, 'g', -1, 64))
	}

	last := len(res) - 1
	//Checks that the data in the last recieved object is correct. Is case-insensitive on the country name.
	if res[last].IsoCode != HISTORY_ALL_SORT_LAST_CODE || !strings.EqualFold(res[last].Name, HISTORY_ALL_SORT_LAST_NAME) || testPercentage(res[last].Percentage, HISTORY_ALL_SORT_LAST_PERCENTAGE) != "" {
		t.Fatal("Wrong object recieved." +
			"\n\tExpected: " + HISTORY_ALL_SORT_LAST_CODE + " - " + HISTORY_ALL_SORT_LAST_NAME + " - " + strconv.FormatFloat(HISTORY_ALL_SORT_LAST_PERCENTAGE, 'g', -1, 64) +
			"\n\tRecieved: " + res[last].IsoCode + " - " + res[last].Name + " - " + strconv.FormatFloat(res[last].Percentage, 'g', -1, 64))
	}
}
