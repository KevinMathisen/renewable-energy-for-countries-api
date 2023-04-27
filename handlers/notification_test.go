package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/structs"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// Slice of all webhooks registered during testing
var gRegisteredWebhooks []structs.Webhook

/*
Handles opening and closing of server, alongside creating and closing client
Then calls given function for testing individual endpoints
*/
func handleNotificationLogistics(t *testing.T, f func(*testing.T, string, http.Client)) {
	handler := RootHandler(Notification)

	server := httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))
	defer server.Close()

	client := http.Client{}
	defer client.CloseIdleConnections()

	log.Println("URL: ", server.URL)

	url := server.URL + constants.NOTIFICATION_PATH

	f(t, url, client)
}

func TestHttpNotification(t *testing.T) {
	// Set up Firestore
	db.InitializeFirestore(constants.CREDENTIALS_FILE_TESTING)
	// Clears cache
	db.DeleteAllDocumentsInCollectionFromFirestore(constants.CACHE_COLLECTION)
	// Clears all webhooks
	db.DeleteAllDocumentsInCollectionFromFirestore(constants.WEBHOOKS_COLLECTION)
	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	handleNotificationLogistics(t, registerWebhook)
	handleNotificationLogistics(t, registerWebhook) //Register another webhook for increased robustness
	handleNotificationLogistics(t, getAWebhook)
	handleNotificationLogistics(t, getAllWebhooks)
	handleNotificationLogistics(t, deleteWebhook)
}

// Tests registration of new webhook
func registerWebhook(t *testing.T, url string, client http.Client) {
	req := strings.NewReader("{\"url\": \"https://webhook.site/09469c1c-abc7-4532-9175-7df9549f7d71\",\"country\": \"NOR\",\"calls\": 2}")

	//Posts a webhook
	log.Println("Testing Post to URL: \"" + url + "\"...")

	//Sends Post request
	res, err := client.Post(url, constants.CONT_TYPE_JSON, req)
	if err != nil {
		t.Fatal("Post request to URL failed:" + err.Error())
	}

	var resObject structs.Webhook
	//Recieves values, and decodes into slice
	err = json.NewDecoder(res.Body).Decode(&resObject)
	if err != nil {
		t.Fatal("Error during decoding:" + err.Error())
	}

	//If there is no WebhookId
	if resObject.WebhookId == "" {
		t.Fatal("Registered webhook does not have associated WebhookId.")
	}

	//Temporarily stores recieved webhook
	gRegisteredWebhooks = append(gRegisteredWebhooks, resObject)
}

// Tests getting a specified webhook
func getAWebhook(t *testing.T, url string, client http.Client) {
	url = url + gRegisteredWebhooks[0].WebhookId

	log.Println("Testing Get to URL: \"" + url + "\"...")

	//Sends Get request
	res, err := client.Get(url)
	if err != nil {
		t.Fatal("Get request to URL failed:" + err.Error())
	}

	//Recieves values, and decodes into struct
	var resObject structs.Webhook
	err = json.NewDecoder(res.Body).Decode(&resObject)
	if err != nil {
		t.Fatal("Error during decoding:" + err.Error())
	}

	//Checks if the webhook registered in registerWebhook() is retrievable
	if resObject.WebhookId != gRegisteredWebhooks[0].WebhookId {
		t.Fatal("Earlier registered webhook not found.")
	}
}

// Tests getting all webhooks
func getAllWebhooks(t *testing.T, url string, client http.Client) {

	log.Println("Testing Get to URL: \"" + url + "\"...")

	//Sends Get request
	res, err := client.Get(url)
	if err != nil {
		t.Fatal("Get request to URL failed:" + err.Error())
	}

	//Recieves values, and decodes into slice
	var resObject []structs.Webhook
	err = json.NewDecoder(res.Body).Decode(&resObject)
	if err != nil {
		t.Fatal("Error during decoding:" + err.Error())
	}

	//Checks if amount of recieved webhooks match the amount of registered webhooks
	if length := len(resObject); length != len(gRegisteredWebhooks) {
		t.Fatal("Amount of webhooks retrieved does not match amount of posts sent." +
			"\n\tExpected: " + strconv.Itoa(len(gRegisteredWebhooks)) +
			"\n\tRecieved: " + strconv.Itoa(length))
	}

	//Looks through resObject to see if gRegisteredWebhooks[0] is there
	found := false
	for _, v := range resObject {
		if v.WebhookId == gRegisteredWebhooks[0].WebhookId {
			found = true
		}
	}
	//If the webhook was not found
	if !found {
		t.Fatal("Could not find earlier registered webhook among recieved webhooks.")
	}
}

// Tests deleting a specified webhook
func deleteWebhook(t *testing.T, url string, client http.Client) {
	//Deletes only the first webhook
	url = url + gRegisteredWebhooks[0].WebhookId

	log.Println("Testing Delete to URL: \"" + url + "\"...")

	//Creates Delete request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal("Failed to create Delete request" + err.Error())
	}

	//Sends Delete request
	res, err := client.Do(req)
	if err != nil {
		t.Fatal("Delete request to URL failed:" + err.Error())
	}

	//If there is content in the response, it is because the deletion failed
	if res.StatusCode != http.StatusNoContent {

		var resObject string
		//Recieves values, and decodes into slice
		err = json.NewDecoder(res.Body).Decode(&resObject)
		if err != nil {
			t.Fatal("Error during decoding:" + err.Error())
		}

		t.Fatal("Error during deletion of webhook:" + resObject)
	}
}
