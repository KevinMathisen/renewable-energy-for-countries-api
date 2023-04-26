package handlers

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"fmt"
	"net/http"
)

func Default(w http.ResponseWriter, r *http.Request) error {

	// Send error message if request method is not get
	if r.Method != http.MethodGet {
		return structs.NewError(nil, http.StatusNotImplemented, "Invalid method, currently only GET is supported", "User used invalid http method")
	}

	// Set content type
	w.Header().Set("content-type", "text/html")

	// Information to display to user on root path
	outout := "This service gives information about developments related to renewable energy production for and across countries. <br> " +
		"For more information about the service, read the readme at https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2023-workspace/raphaesl/group/assignment2"

	// Write information to client
	_, err := fmt.Fprintf(w, "%v", outout)

	// Deal with potential errors
	if err != nil {
		return structs.NewError(err, http.StatusInternalServerError, constants.DEFAULT500, "Error when writing to client")
	}
	return nil
}
