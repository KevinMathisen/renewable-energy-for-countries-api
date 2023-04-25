package gateway

import (
	"assignment2/utils/constants"
	"net/http"
	"net/http/httptest"
	"testing"
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

	country, err:= GetCountryByIso("NOR", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	// Test the response
	if country.Name != "Norway" {
		t.Errorf("Expected Norway, got %s", country.Name)
	}
}