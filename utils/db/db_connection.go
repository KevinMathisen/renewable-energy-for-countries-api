package db

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"context"
	"errors"
	"log"
	"net/http"

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

/*
Sets up Firebase client connection with credentials
*/
func InitializeFirestore() {
	// Firebase initialisation
	firestoreContext = context.Background()

	// Load credentials from json file containing service account
	serviceAccount := option.WithCredentialsFile(constants.CREDENTIALS_FILE)
	// Create a firebase app with context and credentials
	app, err := firebase.NewApp(firestoreContext, nil, serviceAccount)
	if err != nil {
		log.Fatalln(err)
	}

	// Instantiate client and connect to Firestore
	firebaseClient, err = app.Firestore(firestoreContext)
	if err != nil {
		log.Fatalln(err)
	}

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
Return if a country is in the renewables collection

	isoCode	- Code of country to find

	return 	- If country given exists in database
*/
func IsoCodeInDB(isoCode string) bool {
	// Check if country with ISO code exists in renewables collection
	_, err := firebaseClient.Collection(constants.RENEWABLES_COLLECTION).Doc(isoCode).Get(firestoreContext)

	// If we got error not found, return false
	if status.Code(err) == codes.NotFound {
		return false
	}
	// If the country was found
	return true
}

func Firestore_test() error {
	res := firebaseClient.Collection(constants.RENEWABLES_COLLECTION).Doc("NOR")

	doc, _ := res.Get(firestoreContext)

	message := doc.Data()
	log.Println(message)

	return nil
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
Appends a single map with specified webhookID to firestore in a collection under a document in a collection

	doc				- Map of data, where each key will be one field in a document
	collectionName	- Name of root collection to add data to
	isoCode			- Document name of country the webhook is part of
	webhookID		- ID of webhook and document we will add to firestore
*/
func AppendDocumentToWebhooksFirestore(doc map[string]interface{}, collectionName string, isoCode string, webhookID string) error {
	_, err := firebaseClient.Collection(collectionName).Doc(isoCode).Collection(constants.WEBHOOK_COLLECTIONNAME).Doc(webhookID).Set(firestoreContext, doc)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

/*
Get a document from firestore

	w				- Responsewriter for error handling
	id				- document ID to get
	collectionName	- Name of collection to get document from

	return	- Map containing data from document
*/
func GetDocumentFromFirestore(w http.ResponseWriter, id string, collectionName string) (map[string]interface{}, error) {
	// Get reference to document
	docSnapshot, err := firebaseClient.Collection(collectionName).Doc(id).Get(firestoreContext)
	if err != nil {
		http.Error(w, "Error extracting body of document "+id, http.StatusInternalServerError)
		return nil, err
	}

	// Return the data
	return docSnapshot.Data(), nil
}

/*
Gets all documents from a collection in firestore

	w				- Responsewriter for error handling
	collectionName	- Name of collection to get document from

	return 			- Map containing key (document id) and elements containing maps with data from each document
*/
func GetAllDocumentInCollectionFromFirestore(w http.ResponseWriter, collectionName string) (map[string]map[string]interface{}, error) {
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
			http.Error(w, "Failed to iterate through documents in collection "+collectionName+" on firebase", http.StatusInternalServerError)
			return nil, err
		}

		// Save each document with documentID as the key
		data[doc.Ref.ID] = doc.Data()
	}

	return data, nil
}

func CheckCacheDBForURL(w http.ResponseWriter, url string) ([]structs.CountryOutput, error) {
	return nil, nil
}

/*
Delete a document given ID if it exists

	w				- HTTP responsewriter
	documentID		- ID of document to delete
	collectionName	- Name of collection document is in

	return			- If deletion was succesful, if document existed, or any other errors
*/
func DeleteDocument(w http.ResponseWriter, documentID string, collectionName string) error {

	// Get reference to document
	documentRef := firebaseClient.Collection(collectionName).Doc(documentID)

	// Get snapshot of document for testing if it exists
	documentSnap, err := documentRef.Get(firestoreContext)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Test if any document with given ID exists
	if !documentSnap.Exists() {
		// Error, cant delete a document that does not exist
		log.Println("Document in database not found")
		return errors.New("Document in database not found")
	}

	// Delete document if it exists
	documentRef.Delete(firestoreContext)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
