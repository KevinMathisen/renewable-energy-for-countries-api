package gateway

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Tests the query for countries by iso code
*/
func TestGetCountryByIso(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if path := constants.COUNTRY_CODE_SEARCH_PATH + "NOR"; r.URL.Path != path {
			t.Errorf("Expected %s, got %s", path, r.URL.Path)
		}

		// Send response to be tested
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{
				"name": {
					"common": "Norway"
				},
				"cca3": "NOR",
				"borders": [
					"FIN",
					"SWE",
					"RUS"
				]
			}
		  ]`))
	}))
	defer ts.Close()

	country, err := GetCountryByIso("NOR", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	expected := &structs.Country{
		Name:    "Norway",
		IsoCode: "NOR",
		Borders: []string{"FIN", "SWE", "RUS"},
	}

	assert.Equal(t, country, expected, "Response body does not match expected")
}

/*
Tests the error handling of the query for countries by iso code
*/
func TestGetCountryByIsoError(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if path := constants.COUNTRY_CODE_SEARCH_PATH + "Norway"; r.URL.Path != path {
			t.Errorf("Expected %s, got %s", path, r.URL.Path)
		}

		// Send response to be tested
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(``))
	}))
	defer ts.Close()

	_, err := GetCountryByIso("Norway", ts.URL)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

/*
Tests the error handling of the query for countries by iso code
*/
func TestGetCountryByIsoError2(t *testing.T) {
	_, err := GetCountryByIso("Norway", "http://localhost:1234")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

/*
Tests the query for countries by name
*/
func TestGetCountryByName(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if path := constants.COUNTRY_NAME_SEARCH_PATH + "Norway"; r.URL.Path != path {
			t.Errorf("Expected %s, got %s", path, r.URL.Path)
		}

		// Send response to be tested
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{
				"name": {
					"common": "Norway"
				},
				"cca3": "NOR",
				"borders": [
					"FIN",
					"SWE",
					"RUS"
				]
			}
		  ]`))
	}))
	defer ts.Close()

	country, err := GetCountryByName("Norway", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	expected := &structs.Country{
		Name:    "Norway",
		IsoCode: "NOR",
		Borders: []string{"FIN", "SWE", "RUS"},
	}

	assert.Equal(t, country, expected, "Response body does not match expected")
}

/*
Tests the error handling of the query for countries by name
*/
func TestGetCountryByNameError(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if path := constants.COUNTRY_NAME_SEARCH_PATH + "Norway"; r.URL.Path != path {
			t.Errorf("Expected %s, got %s", path, r.URL.Path)
		}

		// Send response to be tested
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(``))
	}))
	defer ts.Close()

	_, err := GetCountryByName("Norway", ts.URL)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

/*
Tests the error handling of the query for countries by name
*/
func TestGetCountryByNameError2(t *testing.T) {
	_, err := GetCountryByName("Norway", "http://localhost:1234")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
