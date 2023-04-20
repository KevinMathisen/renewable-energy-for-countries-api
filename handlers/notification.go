package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"assignment2/utils/div"
	"assignment2/utils/gateway"
	"assignment2/utils/params"
	"assignment2/utils/structs"
	"net/http"
)

func Notification(w http.ResponseWriter, r *http.Request) error {
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
		return nil
		/*
			Status: Not impemented
			Message: Invalid method, currently only [methods supported] is supported
		*/
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
		//error handling
		return err
	}

	// Create and set webhookID
	webhook.WebhookId = div.CreateWebhookId()

	// Save webhook to database
	err = saveWebhook(w, webhook)
	if err != nil {
		// TODO: error handling
		return err
	}

	// Create response
	response := structs.Webhook{
		WebhookId: webhook.WebhookId,
	}

	gateway.RespondToGetRequestWithJSON(w, response, http.StatusCreated)

	return nil
}

/*
Saves a webhook to the correct database collection and document

	w		- Responsewriter for error handling
	webhook	- Struct which contain all relevant information about webhook to save

	return	- Type of error or nil if none
*/
func saveWebhook(w http.ResponseWriter, webhook structs.Webhook) error {
	var isoCode string

	// Set isoCode to ANY if no country specified, else set code provided
	if len(webhook.Country) == 0 {
		isoCode = "ANY"
	} else {
		isoCode = webhook.Country
	}

	// Create map containing data to insert into database
	webhookData := map[string]interface{}{
		"url":         webhook.Url,
		"country":     isoCode,
		"calls":       webhook.Calls,
		"invocations": 0,
	}

	// Save webhook to the database
	//err := db.AppendDocumentToWebhooksFirestore(webhookData, constants.WEBHOOK_COLLECTIONNAME, webhook.Country, webhook.WebhookId)
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

	// Try to delete webhook from database
	err = db.DeleteDocument(w, webhookID, constants.WEBHOOKS_COLLECTION)
	if err != nil {
		return err
	}

	// Create response
	response := structs.Webhook{
		WebhookId: webhookID,
	}

	// Send response to user
	gateway.RespondToGetRequestWithJSON(w, response, http.StatusOK)

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

	// Get webhooks from database
	response, err := getWebhooks(w, webhookID)
	if err != nil {
		return err
	}

	// If one webhook returned, respond with only that one struct
	if len(response) == 1 {
		gateway.RespondToGetRequestWithJSON(w, response[0], http.StatusOK)
		return nil
	}

	// Send list of webhoks as response to user
	gateway.RespondToGetRequestWithJSON(w, response, http.StatusOK)

	return nil
}
