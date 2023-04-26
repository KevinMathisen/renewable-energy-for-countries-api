package handlers

import (
	htu "assignment2/http_test_utils"
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

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
	url = url + htu.COUNTRY_CODE
	historyCountry(t, url, client)
}

// Calls historyCountry(...) with a country code
func historyCountryByName(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_NAME
	historyCountry(t, url, client)
}

// Runs tests for the .../renewables/history/{<NOR>/<norway>} endpoint
func historyCountry(t *testing.T, url string, client http.Client) {

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	if err2 := htu.TestLenLastFirst(res, htu.COUNTRY_EXPECTED_ENTRIES, htu.COUNTRY_CODE, htu.COUNTRY_NAME, constants.OLDEST_YEAR_DB, htu.COUNTRY_OLDEST_PERCENTAGE, htu.COUNTRY_CODE, htu.COUNTRY_NAME, constants.LATEST_YEAR_DB, htu.COUNTRY_LATEST_PERCENTAGE); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?begin={htu.BEGIN_YEAR}&end={htu.END_YEAR} endpoint
func historyCountryBeginEnd(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.BEGIN + htu.AND + htu.END

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	if err2 := htu.TestLenLastFirst(res, htu.COUNTRY_BEGIN_END_ENTRIES, htu.COUNTRY_CODE, htu.COUNTRY_NAME, htu.INT_BEGIN_YEAR, htu.COUNTRY_BEGIN_PERCENTAGE, htu.COUNTRY_CODE, htu.COUNTRY_NAME, htu.INT_END_YEAR, htu.COUNTRY_END_PERCENTAGE); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?begin={htu.BEGIN_YEAR} endpoint
func historyCountryBegin(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.BEGIN

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	if err2 := htu.TestLenLastFirst(res, htu.COUNTRY_BEGIN_ENTRIES, htu.COUNTRY_CODE, htu.COUNTRY_NAME, htu.INT_BEGIN_YEAR, htu.COUNTRY_BEGIN_PERCENTAGE, htu.COUNTRY_CODE, htu.COUNTRY_NAME, constants.LATEST_YEAR_DB, htu.COUNTRY_LATEST_PERCENTAGE); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?end={htu.END_YEAR} endpoint
func historyCountryEnd(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.END

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	if err2 := htu.TestLenLastFirst(res, htu.COUNTRY_END_ENTRIES, htu.COUNTRY_CODE, htu.COUNTRY_NAME, constants.OLDEST_YEAR_DB, htu.COUNTRY_OLDEST_PERCENTAGE, htu.COUNTRY_CODE, htu.COUNTRY_NAME, htu.INT_END_YEAR, htu.COUNTRY_END_PERCENTAGE); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?begin={htu.BEGIN_YEAR}&end={htu.END_YEAR}&sortByValue=true endpoint
func historyCountryBeginEndSort(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.BEGIN + htu.AND + htu.END + htu.AND + htu.SORT_BY

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	if err2 := htu.TestLenLastFirst(res, htu.COUNTRY_BEGIN_END_ENTRIES, htu.COUNTRY_CODE, htu.COUNTRY_NAME, htu.COUNTRY_BEGIN_END_SORT_FIRST, htu.COUNTRY_BEGIN_END_SORT_FIRST_PERCENTAGE, htu.COUNTRY_CODE, htu.COUNTRY_NAME, htu.COUNTRY_BEGIN_END_SORT_LAST, htu.COUNTRY_BEGIN_END_SORT_LAST_PERCENTAGE); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?mean=true endpoint
func historyCountryMean(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.MEAN

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks if there is only one object recieved
	if err2 := htu.TestLen(res, 1); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//Tests the percentage
	if err2 := htu.TestPercentage(res[0].Percentage, htu.COUNTRY_MEAN); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?begin={htu.BEGIN_YEAR}&end={htu.END_YEAR}&mean=true endpoint
func historyCountryBeginEndMean(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.BEGIN + htu.AND + htu.END + htu.AND + htu.MEAN

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks if there is only one object recieved
	if err2 := htu.TestLen(res, 1); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//Tests the percentage
	if err2 := htu.TestPercentage(res[0].Percentage, htu.COUNTRY_BEGIN_END_MEAN); err2 != "" {
		t.Fatal(err2)
	}
}

// ------------------------------ NEIGHBOUR COUNTRY TESTS ------------------------------

// Runs tests for the .../renewables/history/NOR?neighbours=true endpoint
func historyNeighbours(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.NEIGHBOURS

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	if err2 := htu.TestLen(res, htu.NEIGHBOUR_ENTRIES_AMOUNT); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?neighbours=true&begin={htu.BEGIN_YEAR}&end={htu.END_YEAR} endpoint
func historyNeighboursBeginEnd(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.NEIGHBOURS + htu.AND + htu.BEGIN + htu.AND + htu.END

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	if err2 := htu.TestLen(res, htu.NEIGHBOUR_BEGIN_END_AMOUNT); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?neighbours=true&sortByValue=true endpoint
func historyNeighboursSort(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.NEIGHBOURS + htu.AND + htu.SORT_BY

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Tests the amount of recieved objects, then tests the values of the first and then the last object
	if err2 := htu.TestLenLastFirst(res, htu.NEIGHBOUR_ENTRIES_AMOUNT, htu.COUNTRY_CODE, htu.COUNTRY_NAME, htu.INT_BEGIN_YEAR, htu.COUNTRY_BEGIN_END_SORT_FIRST_PERCENTAGE, htu.NEIGHBOURS_SORT_LAST_CODE, htu.NEIGHBOURS_SORT_LAST_NAME, htu.NEIGHBOURS_SORT_LAST_YEAR, htu.NEIGHBOURS_SORT_LAST_PERCENTAGE); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?neighbours=true&mean=true endpoint
func historyNeighboursMean(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.NEIGHBOURS + htu.AND + htu.MEAN

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that there are the expected amount of neighbours
	if err2 := htu.TestLen(res, htu.EXPECTED_NEIGHBOURS); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//Checks that the percentage of Norway is correct. Checks that the order is correct as side-effect
	if err2 := htu.TestPercentage(res[1].Percentage, htu.COUNTRY_MEAN); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/NOR?neighbours=true&mean=true&sortByValue=true endpoint
func historyNeighboursMeanSort(t *testing.T, url string, client http.Client) {
	url = url + htu.COUNTRY_CODE + htu.PARAM + htu.NEIGHBOURS + htu.AND + htu.MEAN + htu.AND + htu.SORT_BY

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//The code in htu.TestLenFirstLast but refactored to not include country year:
	//Tests the amount of recieved objects, then tests the values of the first and then the last object
	if err2 := htu.TestLen(res, htu.EXPECTED_NEIGHBOURS); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the data in the first recieved object is correct. Is case-insensitive on the country name.
	if res[0].IsoCode != htu.COUNTRY_CODE || !strings.EqualFold(res[0].Name, htu.COUNTRY_NAME) || htu.TestPercentage(res[0].Percentage, htu.COUNTRY_MEAN) != "" {
		t.Fatal("Wrong object recieved." +
			"\n\tExpected: " + htu.COUNTRY_CODE + " - " + htu.COUNTRY_NAME + " - " + strconv.FormatFloat(htu.COUNTRY_MEAN, 'g', -1, 64) +
			"\n\tRecieved: " + res[0].IsoCode + " - " + res[0].Name + strconv.FormatFloat(res[0].Percentage, 'g', -1, 64))
	}

	last := len(res) - 1
	//Checks that the data in the last recieved object is correct. Is case-insensitive on the country name.
	if res[last].IsoCode != htu.NEIGHBOURS_SORT_LAST_CODE || !strings.EqualFold(res[last].Name, htu.NEIGHBOURS_SORT_LAST_NAME) || htu.TestPercentage(res[last].Percentage, htu.NEIGHBOURS_SORT_MEAN_LAST) != "" {
		t.Fatal("Wrong object recieved." +
			"\n\tExpected: " + htu.NEIGHBOURS_SORT_LAST_CODE + " - " + htu.NEIGHBOURS_SORT_LAST_NAME + " - " + strconv.FormatFloat(htu.NEIGHBOURS_SORT_MEAN_LAST, 'g', -1, 64) +
			"\n\tRecieved: " + res[last].IsoCode + " - " + res[last].Name + " - " + strconv.FormatFloat(res[last].Percentage, 'g', -1, 64))
	}
}

//------------------------------ ALL COUNTRIES TESTS ------------------------------

// Runs tests for the .../renewables/history/ endpoint
func historyAll(t *testing.T, url string, client http.Client) {

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks the amount of recieved countries
	if err2 := htu.TestLen(res, htu.ALL_COUNTRIES); err2 != "" {
		t.Fatal(err2)
	}
}

// Runs tests for the .../renewables/history/?sortByValue=true endpoint
func historyAllSort(t *testing.T, url string, client http.Client) {
	url = url + htu.PARAM + htu.SORT_BY

	//Gets data from the .../renewables/history/ endpoint
	res, err := htu.GetData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}

	//Checks that the year is not set
	if res[0].Year != "" {
		t.Fatal("Mean of an object is not supposed to have year value.")
	}

	//The code in htu.TestLenFirstLast but refactored to not include country year:
	//Tests the amount of recieved objects, then tests the values of the first and then the last object
	if err2 := htu.TestLen(res, htu.ALL_COUNTRIES); err2 != "" {
		t.Fatal(err2)
	}

	//Checks that the data in the first recieved object is correct. Is case-insensitive on the country name.
	if res[0].IsoCode != htu.COUNTRY_CODE || !strings.EqualFold(res[0].Name, htu.COUNTRY_NAME) || htu.TestPercentage(res[0].Percentage, htu.COUNTRY_MEAN) != "" {
		t.Fatal("Wrong object recieved." +
			"\n\tExpected: " + htu.COUNTRY_CODE + " - " + htu.COUNTRY_NAME + " - " + strconv.FormatFloat(htu.COUNTRY_MEAN, 'g', -1, 64) +
			"\n\tRecieved: " + res[0].IsoCode + " - " + res[0].Name + strconv.FormatFloat(res[0].Percentage, 'g', -1, 64))
	}

	last := len(res) - 1
	//Checks that the data in the last recieved object is correct. Is case-insensitive on the country name.
	if res[last].IsoCode != htu.ALL_SORT_LAST_CODE || !strings.EqualFold(res[last].Name, htu.ALL_SORT_LAST_NAME) || htu.TestPercentage(res[last].Percentage, htu.ALL_SORT_LAST_PERCENTAGE) != "" {
		t.Fatal("Wrong object recieved." +
			"\n\tExpected: " + htu.ALL_SORT_LAST_CODE + " - " + htu.ALL_SORT_LAST_NAME + " - " + strconv.FormatFloat(htu.ALL_SORT_LAST_PERCENTAGE, 'g', -1, 64) +
			"\n\tRecieved: " + res[last].IsoCode + " - " + res[last].Name + " - " + strconv.FormatFloat(res[last].Percentage, 'g', -1, 64))
	}
}
