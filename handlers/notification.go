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
		/*
			Status: Not impemented
			Message: Invalid method, currently only [methods supported] is supported
		*/
	}

	return err
}

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

	gateway.RespondToGetRequestWithJSON(w, response)

	return nil
}

/*
Saves a webhook to the correct database collection and document

	w		- Responsewriter for error handling
	webhook	- Struct which contain all relevant information about webhook to save

	return	- Type of error or nil if none
*/
func saveWebhook(w http.ResponseWriter, webhook structs.Webhook) error {
	// Create map containing data to insert into database
	webhookData := map[string]interface{}{
		"url":         webhook.Url,
		"country":     webhook.Country,
		"calls":       webhook.Calls,
		"invocations": 0,
	}

	// Save webhook to the database
	err := db.AppendDocumentToWebhooksFirestore(webhookData, constants.WEBHOOK_COLLECTIONNAME, webhook.Country, webhook.WebhookId)
	if err != nil {
		// TODO: Error handling
		return err
	}

	return nil
}
