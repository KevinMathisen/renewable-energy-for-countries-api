package gateway

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
