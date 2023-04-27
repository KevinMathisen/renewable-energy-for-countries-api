package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/div"
	"assignment2/utils/gateway"
	"assignment2/utils/params"
	"assignment2/utils/structs"
	"fmt"
	"net/http"
	"time"
)

/*
Handler for notification endpoint
*/
func Notification(w http.ResponseWriter, r *http.Request) error {
	// Check if database is online. If not, give standard error response.
	if !db.DbState {
		usrMsg := fmt.Sprintf("The database is currently unavailable. Please try again later. Reattempting database connection in %v seconds.", time.Until(db.DbRestartTimerStartTime.Add(1*time.Minute)).Round(time.Second)) //Create message with time since timer was activated
		return structs.NewError(nil, http.StatusServiceUnavailable, usrMsg, "")
	}

	var err error

	// Send request to different functions based on method
	switch r.Method {
	case http.MethodPost:
		err = registrationOfWebhook(w, r)
	case http.MethodDelete:
		err = deletionOfWebhook(w, r)
	case http.MethodGet:
		err = viewWebhook(w, r)
	default:
		return structs.NewError(nil, http.StatusNotImplemented, "Invalid method, currently only Post, Delete, Get supported", "User used invalid http method")
	}

	return err
}

/*
Gets webhook data from request, created ID, saves the webhook to the database, then repsponds with the id to the user
*/
func registrationOfWebhook(w http.ResponseWriter, r *http.Request) error {
	// Get request json body
	webhook, err := params.GetWebhookFromRequest(w, r)
	if err != nil {
		return err
	}

	// Create and set webhookID
	webhook.WebhookId = div.CreateWebhookId()

	// Save webhook to database
	err = saveWebhook(webhook)
	if err != nil {
		return err
	}

	// Create response
	response := structs.Webhook{
		WebhookId: webhook.WebhookId,
	}

	err = gateway.RespondToGetRequestWithJSON(w, response, http.StatusCreated)
	if err != nil {
		return err
	}

	return nil
}

/*
Saves a webhook to the correct database collection and document

	webhook	- Struct which contain all relevant information about webhook to save

	return	- Type of error or nil if none
*/
func saveWebhook(webhook structs.Webhook) error {
	var isoCode string
	var year int = -1

	// Set isoCode to ANY if no country specified, else set code provided
	if len(webhook.Country) == 0 {
		isoCode = "ANY"
	} else {
		isoCode = webhook.Country
	}

	// Set year to selected year if it is specified
	if webhook.Year > 0 {
		year = webhook.Year
	}

	// Create map containing data to insert into database
	webhookData := map[string]interface{}{
		"url":         webhook.Url,
		"country":     isoCode,
		"calls":       webhook.Calls,
		"invocations": 0,
		"year":        year,
	}

	// Save webhook to the database
	err := db.AppendDocumentToFirestore(webhook.WebhookId, webhookData, constants.WEBHOOKS_COLLECTION)
	if err != nil {
		// TODO: Error handling
		return err
	}

	return nil
}

/*
Delete webhook and respond with ID of webhook deleted to user
*/
func deletionOfWebhook(w http.ResponseWriter, r *http.Request) error {
	// Get webhookID
	webhookID, err := params.GetWebhookIDFromRequest(w, r)
	if err != nil {
		return err
	}

	// Check if the webhookID is valid
	if !checkIfValidWebhookId(webhookID) {
		return structs.NewError(nil, http.StatusNotFound, "Invalid webhookID given", "webhookID given was not found in database")
	}

	// Try to delete webhook from database
	err = db.DeleteDocument(webhookID, constants.WEBHOOKS_COLLECTION)
	if err != nil {
		return err
	}

	// Respond with status no content
	w.WriteHeader(http.StatusNoContent)

	return nil
}

/*
Get webhook or multiple if none are specified, and then respond to user
*/
func viewWebhook(w http.ResponseWriter, r *http.Request) error {
	// Get webhookID
	webhookID, err := params.GetWebhookIDOrNothingFromRequest(w, r)
	if err != nil {
		return err
	}

	// Check if the webhookID is valid
	if !checkIfValidWebhookId(webhookID) {
		return structs.NewError(nil, http.StatusNotFound, "Invalid webhookID given", "webhookID given was not found in database")
	}

	// Get webhooks from database
	response, err := getWebhooks(webhookID)
	if err != nil {
		return err
	}

	// If one webhook returned, respond with only that one struct
	if len(response) == 1 {
		err = gateway.RespondToGetRequestWithJSON(w, response[0], http.StatusOK)
		if err != nil {
			return err
		}
		return nil
	}

	// Send list of webhoks as response to user
	err = gateway.RespondToGetRequestWithJSON(w, response, http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func checkIfValidWebhookId(id string) bool {
	return id == "" || db.DocumentInCollection(id, constants.WEBHOOKS_COLLECTION)
}

/*
Get webhook from database, and create webhook structs from this data
*/
func getWebhooks(webhookID string) ([]structs.Webhook, error) {
	var webhooks []structs.Webhook
	data := make(map[string]map[string]interface{})
	var err error

	if webhookID != "" {
		// If webhookID is defined, get its data
		webhookData, err := db.GetDocumentFromFirestore(webhookID, constants.WEBHOOKS_COLLECTION)
		if err != nil {
			// Error handling
			return webhooks, err
		}
		// Save the webhooks data
		data[webhookID] = webhookData

	} else {
		// If no webhookID is given, get all webhooks data
		data, err = db.GetAllDocumentInCollectionFromFirestore(constants.WEBHOOKS_COLLECTION)
		if err != nil {
			// Error handling
			return webhooks, err
		}
	}

	for webhookID, webhookData := range data {
		// For each webhook found in database, create struct from it
		webhook := structs.CreateWebhookFromData(webhookData, webhookID)
		if err != nil {
			// Error handling
			return webhooks, err
		}

		// Save the created struct
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}
