package handlers

import (
	"net/http"
	"time"
)

// Create time of service start variable to calculate uptime
var Start time.Time

func Status(w http.ResponseWriter, r *http.Request) {

	// Send error if request is not GET:
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method, currently only GET is supported", http.StatusNotImplemented)
		return
	}

	// Generate status response
	statusRes, err := createStatusResponse(Start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle get request
	respondToGetRequestWithJSON(w, statusRes)
}
