package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/structs"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	//"strconv"
	//"strings"
	"testing"
)

const HISTORY_COUNTRY_CODE = "NOR"
const HISTORY_COUNTRY_NAME = "norway"
const BEGIN_YEAR = "1990"
const END_YEAR = "2010"
const HISTORY_PARAM = "?"
const HISTORY_AND = "&"
const HISTORY_BEGIN = "begin=" + BEGIN_YEAR
const HISTORY_END = "end=" + END_YEAR
const HISTORY_SORT_BY = "sortByValue=true"
const HISTORY_NEIGHBOURS = "neighbours=true"
const HISTORY_MEAN = "mean=true"

/*
Gets data from the test URL and decodes into a slice of type CountryOutput, then returns this if
there are no errors
*/
func getHistoryData(client http.Client, url string) ([]structs.CountryOutput, error) {
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

func TestHttpGetRenewablesHistory(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllCachedRequestsFromFirestore()
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
	handleHistoryLogistics(t, historyCountryBeginEndMeanSort)
	//Neighbour
	handleHistoryLogistics(t, historyNeighbours)
	handleHistoryLogistics(t, historyNeighboursBeginEnd)
	handleHistoryLogistics(t, historyNeighboursMean)
	handleHistoryLogistics(t, historyNeighboursSort)
	//All
	handleHistoryLogistics(t, historyAll)
	handleHistoryLogistics(t, historyAllSort)
	handleHistoryLogistics(t, historyAllMean)
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

func historyCountryByCode(t *testing.T, url string, client http.Client) {
	fullUrl := url + HISTORY_COUNTRY_CODE
	historyCountry(t, fullUrl, client)
}

func historyCountryByName(t *testing.T, url string, client http.Client) {
	fullUrl := url + HISTORY_COUNTRY_NAME
	historyCountry(t, fullUrl, client)
}

// Runs tests for the .../renewables/history/{<NOR>/<norway>} endpoint
func historyCountry(t *testing.T, url string, client http.Client) {
	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

// Runs tests for the .../renewables/history/NOR?begin={BEGIN_YEAR}&end={END_YEAR} endpoint
func historyCountryBeginEnd(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_BEGIN + HISTORY_AND + HISTORY_END

	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

// Runs tests for the .../renewables/history/NOR?begin={BEGIN_YEAR} endpoint
func historyCountryBegin(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_BEGIN

	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

// Runs tests for the .../renewables/history/NOR?end={END_YEAR} endpoint
func historyCountryEnd(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_END

	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

// Runs tests for the .../renewables/history/NOR?begin={BEGIN_YEAR}&end={END_YEAR}&sortByValue=true endpoint
func historyCountryBeginEndSort(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_BEGIN + HISTORY_AND + HISTORY_END + HISTORY_AND + HISTORY_SORT_BY

	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

func historyCountryMean(t *testing.T, url string, client http.Client) {
	url = url + HISTORY_COUNTRY_CODE + HISTORY_PARAM + HISTORY_MEAN

	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

func historyCountryBeginEndMean(t *testing.T, url string, client http.Client) {

}

func historyCountryBeginEndMeanSort(t *testing.T, url string, client http.Client) {

}

// ------------------------------ NEIGHBOUR COUNTRY TESTS ------------------------------

func historyNeighbours(t *testing.T, url string, client http.Client) {

}

func historyNeighboursBeginEnd(t *testing.T, url string, client http.Client) {

}

func historyNeighboursSort(t *testing.T, url string, client http.Client) {

}

func historyNeighboursMean(t *testing.T, url string, client http.Client) {

}

//------------------------------ ALL COUNTRIES TESTS ------------------------------

func historyAll(t *testing.T, url string, client http.Client) {
	log.Println("Testing URL: \"" + url + "\"...")

	//Gets data from the .../renewables/history/ endpoint
	res, err := getHistoryData(client, url)
	//If there was an error during gathering or decoding of data
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(res[0].Name)
}

func historyAllSort(t *testing.T, url string, client http.Client) {

}

func historyAllMean(t *testing.T, url string, client http.Client) {

}
