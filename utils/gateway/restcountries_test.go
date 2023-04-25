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
	clearRcCache()
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
	clearRcCache()
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
	clearRcCache()
	_, err := GetCountryByIso("Norway", "http://localhost:1234")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

/*
Tests the query for countries by name
*/
func TestGetCountryByName(t *testing.T) {
	clearRcCache()
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
	clearRcCache()
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
	clearRcCache()
	_, err := GetCountryByName("Norway", "http://localhost:1234")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

/*
Tests if the cache is functioning as espected.
First query for a country by name, with no response. Then, query by iso code with response, and check if the cache is populated.
*/
func TestRcCacheIsoThenName(t *testing.T) {
	clearRcCache()
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if path := constants.COUNTRY_CODE_SEARCH_PATH + "NOR"; r.URL.Path == path {
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

		} else if path := constants.COUNTRY_NAME_SEARCH_PATH + "Norway"; r.URL.Path == path {
			// Send response to be tested
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(``))

		} else {
			t.Errorf("Incorrect path: %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	_, err := GetCountryByName("Norway", ts.URL)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	expected := &structs.Country{
		Name:    "Norway",
		IsoCode: "NOR",
		Borders: []string{"FIN", "SWE", "RUS"},
	}

	country, err := GetCountryByIso("NOR", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	assert.Equal(t, country, expected, "Response body does not match expected")

	country, err = GetCountryByName("Norway", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	assert.Equal(t, country, expected, "Response body does not match expected")
}

/*
Tests if the cache is functioning as espected.
First query for a country by iso code, with no response. Then, query by name with response, and check if the cache is populated.
*/
func TestRcCacheNameThenIso(t *testing.T) {
	clearRcCache()
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if path := constants.COUNTRY_CODE_SEARCH_PATH + "NOR"; r.URL.Path == path {
			// Send response to be tested
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(``))

		} else if path := constants.COUNTRY_NAME_SEARCH_PATH + "Norway"; r.URL.Path == path {
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

		} else {
			t.Errorf("Incorrect path: %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	_, err := GetCountryByIso("NOR", ts.URL)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	expected := &structs.Country{
		Name:    "Norway",
		IsoCode: "NOR",
		Borders: []string{"FIN", "SWE", "RUS"},
	}

	country, err := GetCountryByName("Norway", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	assert.Equal(t, country, expected, "Response body does not match expected")

	country, err = GetCountryByIso("NOR", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	assert.Equal(t, country, expected, "Response body does not match expected")
}

/*
Tests converting a country name to an iso code.
*/
func TestGetIsoCodeFromName(t *testing.T) {
	clearRcCache()
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

	expected := "NOR"
	isoCode, err := GetIsoCodeFromName("Norway", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	assert.Equal(t, isoCode, expected, "Response body does not match expected")
}

/*
Tests converting a country name to an iso code, with no response.
*/
func TestGetIsoCodeFromNameError(t *testing.T) {
	clearRcCache()
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

	_, err := GetIsoCodeFromName("Norway", ts.URL)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

/*
Tests getting the neighbours of a country.
*/
func TestGetNeighbours(t *testing.T) {
	clearRcCache()
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

	expected := []string{"FIN", "SWE", "RUS"}
	neighbours, err := GetNeighbours("NOR", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	assert.Equal(t, neighbours, expected, "Response body does not match expected")
}

/*
Tests getting the neighbours of a country, with no borders.
*/
func TestGetNeighboursNoBorders(t *testing.T) {
	clearRcCache()
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
				"borders": []
			}
		  ]`))
	}))
	defer ts.Close()

	var expected []string
	neighbours, err := GetNeighbours("NOR", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	assert.Equal(t, expected, neighbours, "Response body does not match expected")
}

/*
Tests getting the neighbours of a country, with no response.
*/
func TestGetNeighboursError(t *testing.T) {
	clearRcCache()
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
		w.Write([]byte(``))
	}))
	defer ts.Close()

	_, err := GetNeighbours("NOR", ts.URL)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
