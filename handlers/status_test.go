package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/structs"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Tests the createStatusResponse function
*/
func TestCreateStatusResponse(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllDocumentsInCollectionFromFirestore(constants.CACHE_COLLECTION)
	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	// Set up countries server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			t.Errorf("Expected 'HEAD' request, got '%s'", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Get status response
	res, err := createStatusResponse(ts.URL, time.Now())
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	// Assert response
	assert.Equal(t, "200 OK", res.CountriesApi, "Response status code should be 200.")
	assert.Equal(t, "200 OK", res.NotificationDb, "Response status code should be 200.")
	assert.NotNil(t, res.Webhooks, "Response webhooks should not be nil.")
	assert.NotNil(t, res.Version, "Response cached requests should not be nil.")
	assert.NotNil(t, res.Uptime, "Response should be nil.")
	assert.IsType(t, float64(0), res.Uptime, "Response uptime should be of type float64.")
}

/*
Tests the createStatusResponse function, with a bad response from the countries server
*/
func TestCreateStatusResponseBadUrl(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllDocumentsInCollectionFromFirestore(constants.CACHE_COLLECTION)
	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	// Set up countries server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			t.Errorf("Expected 'HEAD' request, got '%s'", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	// Get status response
	res, err := createStatusResponse(ts.URL, time.Now())
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	// Assert response
	assert.Equal(t, "500 Internal Server Error", res.CountriesApi, "Response status code should be 500.")
	assert.Equal(t, "200 OK", res.NotificationDb, "Response status code should be 200.")
	assert.NotNil(t, res.Webhooks, "Response webhooks should not be nil.")
	assert.NotNil(t, res.Version, "Response cached requests should not be nil.")
	assert.NotNil(t, res.Uptime, "Response should be nil.")
	assert.IsType(t, float64(0), res.Uptime, "Response uptime should be of type float64.")
}

/*
Tests the status handler
*/
func TestHttpStatus(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllDocumentsInCollectionFromFirestore(constants.CACHE_COLLECTION)
	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	//Creates instance of Status handler
	handler := RootHandler(Status)

	//Runs handler instance as server
	server := httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))
	defer server.Close()

	//Creates client to speak with server
	client := http.Client{}
	defer client.CloseIdleConnections()

	log.Println("URL: ", server.URL)

	url := server.URL + constants.STATUS_PATH
	log.Println("Testing URL: \"" + url + "\"...")

	//Sends Get request
	res, err := client.Get(url)
	if err != nil {
		t.Fatal("Get request to URL failed:" + err.Error())
	}

	var resObject structs.Status
	//Recieves values, and decodes into slice
	err = json.NewDecoder(res.Body).Decode(&resObject)
	if err != nil {
		t.Fatal("Error during decoding:" + err.Error())
	}

	//Asserts response
	assert.Equal(t, "200 OK", resObject.NotificationDb, "Response status code should be 200.")
	assert.NotNil(t, resObject.CountriesApi, "Response countries api should not be nil.")
	assert.NotNil(t, resObject.Webhooks, "Response webhooks should not be nil.")
	assert.NotNil(t, resObject.Version, "Response cached requests should not be nil.")
	assert.NotNil(t, resObject.Uptime, "Response should be nil.")
	assert.IsType(t, float64(0), resObject.Uptime, "Response uptime should be of type float64.")
}

/*
Tests the status handler with a bad database connection
*/
func TestHttpStatusWithBadDatabase(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllDocumentsInCollectionFromFirestore(constants.CACHE_COLLECTION)
	// Close down client before service is done running
	db.CloseFirebaseClient()

	//Creates instance of Status handler
	handler := RootHandler(Status)

	//Runs handler instance as server
	server := httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))
	defer server.Close()

	//Creates client to speak with server
	client := http.Client{}
	defer client.CloseIdleConnections()

	log.Println("URL: ", server.URL)

	url := server.URL + constants.STATUS_PATH
	log.Println("Testing URL: \"" + url + "\"...")

	//Sends Get request
	res, err := client.Get(url)
	if err != nil {
		t.Fatal("Get request to URL failed:" + err.Error())
	}

	var resObject structs.Status
	//Recieves values, and decodes into slice
	err = json.NewDecoder(res.Body).Decode(&resObject)
	if err != nil {
		t.Fatal("Error during decoding:" + err.Error())
	}

	//Asserts response
	assert.Equal(t, "503 Service Unavailable", resObject.NotificationDb, "Response status code should be 503.")
	assert.NotNil(t, resObject.CountriesApi, "Response countries api should not be nil.")
	assert.NotNil(t, resObject.Webhooks, "Response webhooks should not be nil.")
	assert.Equal(t, -1, resObject.Webhooks, "Response webhooks should be -1.")
	assert.NotNil(t, resObject.Version, "Response cached requests should not be nil.")
	assert.NotNil(t, resObject.Uptime, "Response should be nil.")
	assert.IsType(t, float64(0), resObject.Uptime, "Response uptime should be of type float64.")
}
