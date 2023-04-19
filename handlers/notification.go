package handlers

import (
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
