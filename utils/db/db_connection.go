package db

import (
	"assignment2/utils/constants"
	"context"
	"log"

	"cloud.google.com/go/firestore" // Firestore-specific support
	firebase "firebase.google.com/go"
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
		log.Fatal("CLosing of firebase client failed. Error: ", err)
	}
}

/*
Return if a country is in the renewables collection

	isoCode	- Code of country to find

	return 	- If country given exists in database
*/
func isoCodeInDB(isoCode string) (bool, error) {
	// Check if country with ISO code exists in renewables collection
	_, err := firebaseClient.Collection(constants.RENEWABLES_COLLECTION).Doc(isoCode).Get(firestoreContext)

	// If we got error not found, return false
	if status.Code(err) == codes.NotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	// IF the country was found
	return true, nil
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
