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

func TestCreateStatusResponse(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllCachedRequestsFromFirestore()
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
	
	res, err := createStatusResponse(ts.URL, time.Now())
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}


	assert.Equal(t, "200 OK", res.CountriesApi, "Response status code should be 200.")
	assert.Equal(t, "200 OK", res.NotificationDb, "Response status code should be 200.")
	assert.NotNil(t, res.Webhooks, "Response webhooks should not be nil.")
	assert.NotNil(t, res.Version, "Response cached requests should not be nil.")
	assert.NotNil(t, res.Uptime, "Response should be nil.")
	assert.IsType(t, float64(0), res.Uptime, "Response uptime should be of type float64.")
}
