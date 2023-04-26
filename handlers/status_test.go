package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
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
