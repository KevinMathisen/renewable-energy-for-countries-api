package gateway

import (
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
