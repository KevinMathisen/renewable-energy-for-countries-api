package db

import (
	"assignment2/utils/constants"
	"assignment2/utils/div"
	"assignment2/utils/gateway"
	"assignment2/utils/structs"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/firestore" // Firestore-specific support
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Firebase context used by Firestore functions
var firestoreContext context.Context

// Firebase client used by Firestore functions
var firebaseClient *firestore.Client

// Credentials file path
var credentials string

// Boolean variable and accompanying lock to determine the state of the database. Toggle with ReportDbState()
var (
	DbState                 bool       = true
	dbStateMutex            sync.Mutex // Mutual exclusion lock to safely toggle db error state
	dbRestartTimerMutex     sync.Mutex // Mutual exclusion lock to wait for database reconnection attempt. It's purspose is to prevent multiple timers being started.
	DbRestartTimerStartTime time.Time  // Moment of timer activation to display to user.
)

/*
Sets up Firebase client connection with credentials
Returns error
*/
func InitializeFirestore(credPath string) error {
	// Set credentials path
	credentials = credPath
	// Firebase initialisation
	firestoreContext = context.Background()

	// Load credentials from json file containing service account
	serviceAccount := option.WithCredentialsFile(credentials)
	// Create a firebase app with context and credentials
	app, err := firebase.NewApp(firestoreContext, nil, serviceAccount)
	if err != nil {
		return structs.NewError(err, http.StatusGatewayTimeout, constants.DEFAULT504, "Could not establish firebase context")
	}

	// Instantiate client and connect to Firestore
	firebaseClient, err = app.Firestore(firestoreContext)
	if err != nil {
		return structs.NewError(err, http.StatusGatewayTimeout, constants.DEFAULT504, "Could not contact firebase")
	}

	return nil
}

/*
Closes the firebase client, or logs a fatal error if it failed
*/
func CloseFirebaseClient() {
	err := firebaseClient.Close()
	if err != nil {
		log.Fatal("Closing of firebase client failed. Error: ", err)
	}
}

/*
Return if a document with ID is in the collection given

	id			- ID of document to find
	collection	- collection to search in

	return 	- If document given exists in database
*/
func DocumentInCollection(id string, collection string) bool {
	// Check if document with id exists in collection
	_, err := firebaseClient.Collection(collection).Doc(id).Get(firestoreContext)

	// If we got error not found, return false
	if status.Code(err) == codes.NotFound || err != nil {
		return false
	}
	// If the document was found
	return true
}

/*
Appends data from a map to firestore

	data			- Map of data, where each key will be the name of a document, and each element will be the document content
	collectionName	- Name of collection to add data to
*/
func AppendDataToFirestore(data map[string]map[string]interface{}, collectionName string) error {

	// For each key value pair in map, add the map to firestore
	for code, element := range data {
		err := AppendDocumentToFirestore(code, element, collectionName)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Appends a single map with specified id to firestore

	id				- Id of document we are creating
	doc				- Map of data, where each key will be one field in a document
	collectionName	- Name of collection to add data to
*/
func AppendDocumentToFirestore(id string, doc map[string]interface{}, collectionName string) error {
	_, err := firebaseClient.Collection(collectionName).Doc(id).Set(firestoreContext, doc)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

/*
Get a document from firestore

	id				- document ID to get
	collectionName	- Name of collection to get document from

	return	- Map containing data from document
*/
func GetDocumentFromFirestore(id string, collectionName string) (map[string]interface{}, error) {
	// Get reference to document
	docSnapshot, err := firebaseClient.Collection(collectionName).Doc(id).Get(firestoreContext)
	if err != nil {

		if !checkDbState() {
			ReportDbState(false)
		}

		return nil, structs.NewError(err, http.StatusBadGateway, constants.DEFAULT504, "Could not reach firestone database. Error extracting body of document "+id)
	}

	// Return the data
	return docSnapshot.Data(), nil
}

/*
Gets all documents from a collection in firestore

	collectionName	- Name of collection to get document from

	return 			- Map containing key (document id) and elements containing maps with data from each document
*/
func GetAllDocumentInCollectionFromFirestore(collectionName string) (map[string]map[string]interface{}, error) {
	// Initialize map for saving documents
	data := make(map[string]map[string]interface{})

	// Get reference to documents in collection
	iter := firebaseClient.Collection(collectionName).Documents(firestoreContext)

	for {
		// Try to go to next document in collection
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

			if !checkDbState() {
				ReportDbState(false)
			}

			return nil, structs.NewError(err, http.StatusBadGateway, constants.DEFAULT504, "Failed to iterate through documents in collection "+collectionName+" on firebase")
		}

		// Save each document with documentID as the key
		data[doc.Ref.ID] = doc.Data()
	}

	return data, nil
}

/*
Delete a document given ID if it exists

	documentID		- ID of document to delete
	collectionName	- Name of collection document is in

	return			- If deletion was succesful, if document existed, or any other errors
*/
func DeleteDocument(documentID string, collectionName string) error {

	// Get reference to document
	documentRef := firebaseClient.Collection(collectionName).Doc(documentID)

	// Get snapshot of document for testing if it exists
	documentSnap, err := documentRef.Get(firestoreContext)
	if err != nil {
		return structs.NewError(err, http.StatusNotFound, "Could not find given webhookID in database", "Could not retrieve firestone document refrence.")
	}

	// Test if any document with given ID exists
	if !documentSnap.Exists() {

		if !checkDbState() {
			ReportDbState(false)
		}

		// Error, cant delete a document that does not exist
		return structs.NewError(err, http.StatusNotFound, "Could not find given webhookID in database", "Found reference to document, but it doesn't exist in database.")

	}

	// Delete document if it exists
	documentRef.Delete(firestoreContext)
	if err != nil {
		return structs.NewError(err, http.StatusBadGateway, constants.DEFAULT504, "Could not delete document from firestone database.")
	}

	return nil
}

/*
Go through all webhooks and check if they are to be invoked

	isoCode	- Isocode of countries to be invoked, empty if all countries
*/
func InvokeCountry(isoCode []string, begin int, end int) {
	// Get reference to documents in collection
	webhooksCollection := firebaseClient.Collection(constants.WEBHOOKS_COLLECTION)
	iter := webhooksCollection.Documents(firestoreContext)

	// Go through all webhooks
	for {
		// Try to get next document in collection
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}

		// Get the webhook data
		webhook := doc.Data()

		// only want webhook if webhook country is one of the invoked countries, or we invoked all countries, or the webhook is invoked for all countries
		if len(isoCode) != 0 && webhook["country"].(string) != "ANY" && !div.Contains(isoCode, webhook["country"].(string)) {
			continue
		}

		// if year is specified, we only want to invoke if the year is between begin and end year
		if webhook["year"].(int64) != -1 && (int(webhook["year"].(int64)) < begin || int(webhook["year"].(int64)) > end) {
			continue
		}

		// Increase invocation count by one
		webhook["invocations"] = webhook["invocations"].(int64) + 1

		// Check if we have met the required invokation amount
		if webhook["invocations"].(int64)%webhook["calls"].(int64) == 0 {
			// Send post to webhook
			go gateway.PostToWebhook(webhook, doc.Ref.ID, constants.COUNTRIES_API_URL)
		}

		// Update webhook with new invocations
		webhooksCollection.Doc(doc.Ref.ID).Set(firestoreContext, webhook)

	}
}

/*
Count amount of webhooks in firestore
*/
func CountWebhooks() (int, error) {
	var amountOfWebhooks int
	// Get reference to documents in collection
	iter := firebaseClient.Collection(constants.WEBHOOKS_COLLECTION).Documents(firestoreContext)

	// Go through all webhooks, count each
	for {
		// Try to get next document in collection
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return -1, structs.NewError(err, http.StatusBadGateway, constants.DEFAULT504, "Could not retrieve doc from database while counting webhooks.")
		}

		amountOfWebhooks += 1

	}

	return amountOfWebhooks, nil
}

/*
Delete all documents in cache collection
*/
func DeleteAllDocumentsInCollectionFromFirestore(collection string) {
	// Get reference to documents in collection
	iter := firebaseClient.Collection(collection).DocumentRefs(firestoreContext)

	for {
		// Try to go to next document in collection
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}

		// Save each document with documentID as the key
		doc.Delete(firestoreContext)
	}
}

/*
* A function used for reporting on the state of the database.
* To signify a healthy database, send "true". This will disable the db-error handeling middleware.
* To signify a broken database, send "false". This will trigger the db-error handeling middleware located in the root handler within "errorHandler".
* Automaticly runs sleepAndRestartDb if action is false and goes through lock.
 */
func ReportDbState(action bool) {

	//Check if dbstate is already equal to desired state
	if DbState == action {
		return
	}

	dbStateMutex.Lock()   //Lock variable
	DbState = action      //Change dbstate
	dbStateMutex.Unlock() //Unlock variable

	//If the reported state is false, db is set to sleep and reattempt connection.
	if !action {
		sleepAndRestartDb()
	}
}

/*
* Function that sleeps for one minute before reattempting a database connection.
 */
func sleepAndRestartDb() {

	//Check if timer is already active
	if !dbRestartTimerMutex.TryLock() {
		return
	}

	//Sleep for a minute
	log.Println("Retrying for 1 minute")
	DbRestartTimerStartTime = time.Now()
	time.Sleep(1 * time.Minute)

	err := InitializeFirestore(credentials) //Reattempt database connection
	dbRestartTimerMutex.Unlock()            //Give away lock regardless of output
	if err != nil {
		sleepAndRestartDb() //On database failure, restart function
	} else {
		ReportDbState(true) //On database success, set flag to true.
	}
}

/*
* Returns the status code of the database.
 */
func GetDbResponse() (http.Response, error) {
	collectionID := constants.WEBHOOKS_COLLECTION
	docID := "non-existent-doc"

	_, err := firebaseClient.Collection(collectionID).Doc(docID).Get(firestoreContext)
	st, _ := status.FromError(err)

	// Check if the error is related to the document not being found
	if st.Code() == codes.NotFound {
		statusCode := http.StatusOK
		return http.Response{
			Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
			StatusCode: statusCode,
		}, nil
	}
	// If the error is not related to the document not being found, return a 503
	statusCode := http.StatusServiceUnavailable
	return http.Response{
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		StatusCode: statusCode,
	}, structs.NewError(err, statusCode, constants.DEFAULT503, "Could not retrieve doc from database.")
}

/*
* Function that checks if the database is healthy.
* Returns true if database is healthy, false if database is unhealthy.
 */
func checkDbState() bool {
	res, _ := GetDbResponse()
	if res.StatusCode == 200 {
		return true
	} else {
		return false
	}
}
