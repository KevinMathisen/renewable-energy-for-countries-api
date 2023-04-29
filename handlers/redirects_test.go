package handlers

import (
	"assignment2/utils/constants"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const EXPECTED_CONTENT = "This service gives information about developments related to renewable energy production for and across countries. <br> " +
	"For more information about the service, read the readme at https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2023-workspace/raphaesl/group/assignment2"

/*
Gets a response from the test URL and decodes into a string, then returns this if there are no errors
*/
func GetResponse(client http.Client, url string) (string, error) {
	log.Println("Testing URL: \"" + url + "\"...")

	//Sends Get request
	res, err := client.Get(url)
	if err != nil {
		log.Println("Get request to URL failed:")
		return "", err
	}

	//Recieves contents of the site, and decodes into string
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error during decoding:")
		return "", err
	}
	response := string(buf)

	return response, nil
}

/*
Handles opening and closing of server, alongside creating and closing client
Then calls given function for testing individual endpoints
*/
func handleRedirectLogistics(t *testing.T, f func(*testing.T, string, http.Client)) {
	//Creates instance of RenewablesCurrent handler
	handler := RootHandler(Default)

	//Runs handler instance as server
	server := httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))
	defer server.Close()

	//Creates client to speak with server
	client := http.Client{}
	defer client.CloseIdleConnections()

	log.Println("URL: ", server.URL)

	url := server.URL + constants.DEFAULT_PATH

	f(t, url, client)
}

/*
Calls all the functions for testing the default endpoints
*/
func TestHttpDefault(t *testing.T) {

	handleRedirectLogistics(t, testDefaultRedirect)
	handleRedirectLogistics(t, testServiceRedirect)
	handleRedirectLogistics(t, testServiceVersionRedirect)
	handleRedirectLogistics(t, testMangledRedirect)
}

// Tests the .../ endpoint
func testDefaultRedirect(t *testing.T, url string, client http.Client) {

	res, err := GetResponse(client, url)
	if err != nil {
		t.Fatal(err)
	}

	if res != EXPECTED_CONTENT {
		t.Fatal("Web page did not have expected content.")
	}
}

// Tests the .../energy/ endpoint
func testServiceRedirect(t *testing.T, url string, client http.Client) {
	url = url + "energy/"

	res, err := GetResponse(client, url)
	if err != nil {
		t.Fatal(err)
	}

	if res != EXPECTED_CONTENT {
		t.Fatal("Web page did not have expected content.")
	}
}

// Tests the .../energy/v1/ endpoint
func testServiceVersionRedirect(t *testing.T, url string, client http.Client) {
	url = url + "energy/" + constants.VERSION

	res, err := GetResponse(client, url)
	if err != nil {
		t.Fatal(err)
	}

	if res != EXPECTED_CONTENT {
		t.Fatal("Web page did not have expected content.")
	}
}

// Tests mangled url: .../energ/v/
func testMangledRedirect(t *testing.T, url string, client http.Client) {
	url = url + "energ/v/"

	res, err := GetResponse(client, url)
	if err != nil {
		t.Fatal(err)
	}

	if res != EXPECTED_CONTENT {
		t.Fatal("Web page did not have expected content.")
	}
}
