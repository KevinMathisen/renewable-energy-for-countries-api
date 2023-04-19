package handlers

import (
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
