package gateway

import (
	"assignment2/utils/constants"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Tests the RespondToGetRequestWithJSON function
*/
func TestRespondToGetRequestWithJSON(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Define body to be marshalled into JSON
		body := struct {
			Title string `json:"title"`
			Msg   string `json:"msg"`
		}{
			Title: "Hello",
			Msg:   "Hello World!",
		}

		RespondToGetRequestWithJSON(w, body, 200)
	}))

	defer ts.Close()

	// Make a request to the test server
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Define expected response body
	expectedBody := `{"title":"Hello","msg":"Hello World!"}`
	// Define expected response status code
	expectedStatus := 200
	// Expected content type
	expectedContentType := "application/json"

	assert.Equal(t, expectedStatus, res.StatusCode, "Response status code should be 200.")
	assert.Contains(t, string(body), expectedBody, "Response body should contain expected body.")
	assert.Equal(t, expectedContentType, res.Header.Get("content-type"), "Response content type should be equal to expected content type.")
}

/*
Tests the HttpRequestFromUrl function
*/
func TestHttpRequestFromUrl(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Define body to be marshalled into JSON
		body := struct {
			Title string `json:"title"`
			Msg   string `json:"msg"`
		}{
			Title: "Hello",
			Msg:   "Hello World!",
		}

		RespondToGetRequestWithJSON(w, body, 200)
	}))

	defer ts.Close()

	// Make a request to the test server
	res, err := HttpRequestFromUrl(ts.URL, "GET")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Define expected response body
	expectedBody := `{"title":"Hello","msg":"Hello World!"}`
	// Define expected response status code
	expectedStatus := 200
	// Expected content type
	expectedContentType := "application/json"

	assert.Equal(t, expectedStatus, res.StatusCode, "Response status code should be 200.")
	assert.Contains(t, string(body), expectedBody, "Response body should contain expected body.")
	assert.Equal(t, expectedContentType, res.Header.Get("content-type"), "Response content type should be equal to expected content type.")
}

/*
Tests the HttpRequestFromUrl function with a non-existing URL
*/
func TestHttpRequestFromUrlNonExistingUrl(t *testing.T) {
	// Make a request to the test server
	_, err := HttpRequestFromUrl("http://localhost:1234", "GET")
	if err == nil {
		t.Fatal("Expected error, but got nil.")
	}
}

/*
Tests the PostToWebhook function
*/
func TestPostToWebhook(t *testing.T) {
	count := 0
	webhookCount := 0
	apiCount := 0
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		if r.URL.Path == "/webhook" {
			webhookCount++
			if r.Method != "POST" {
				t.Fatal("Expected POST request, but got ", r.Method)
			}

			// Unmarshal request body
			var data map[string]interface{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&data)
			if err != nil {
				t.Fatal(err)
			}

			// Define expected request body
			expectedData := map[string]interface{}{
				"webhook_id": "TEST",
				"country":    "Norway",
				"calls":      float64(5),
			}

			assert.Equal(t, expectedData, data, "Request body should be equal to expected request body.")

			// Respond to request
			w.WriteHeader(http.StatusOK)
		}
		if r.URL.Path == "/api"+constants.COUNTRY_CODE_SEARCH_PATH+"NOR" {
			apiCount++
			if r.Method != "GET" {
				t.Fatal("Expected POST request, but got ", r.Method)
			}

			// Respond with JSON
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
		}

	}))

	defer ts.Close()

	data := make(map[string]interface{})

	data["url"] = ts.URL + "/webhook"
	data["country"] = "NOR"
	data["year"] = int64(-1)
	data["invocations"] = int64(5)

	PostToWebhook(data, "TEST", ts.URL+"/api")

	assert.Equal(t, 1, webhookCount, "Webhook should be called once.")
	assert.Equal(t, 1, apiCount, "API should be called once.")
	assert.Equal(t, 2, count, "Webhook and API should be called once each.")
}
