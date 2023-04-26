package db

import (
	"assignment2/utils/constants"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Tests the connection to the database
*/
func TestDBConnection(t *testing.T) {
	// Set up Firestore
	InitializeFirestore("../" + constants.CREDENTIALS_FILE_TESTING)
	// Close down client when service is done running
	defer CloseFirebaseClient()

	// Test connection
	data := map[string]map[string]interface{}{
		"FpLSjFbcXoEFfRsW": {
			"calls":       int64(5),
			"country":     "NOR",
			"invocations": int64(0),
			"url":         "test1.test",
			"year":        int64(-1),
		},
		"QfwLosaJKVANmUJk": {
			"calls":       int64(2),
			"country":     "ANY",
			"invocations": int64(0),
			"url":         "test2.test",
			"year":        int64(2021),
		},
	}

	
	DeleteAllDocumentsInCollectionFromFirestore(constants.WEBHOOKS_COLLECTION)

	err := AppendDataToFirestore(data, constants.WEBHOOKS_COLLECTION)
	if err != nil {
		t.Errorf("Couldn't append data to firestore: " + err.Error())
	}

	webhook1, err := GetDocumentFromFirestore("FpLSjFbcXoEFfRsW", constants.WEBHOOKS_COLLECTION)
	if err != nil {
		t.Errorf("Couldn't get document from firestore: " + err.Error())
	}
	webhook2, err := GetDocumentFromFirestore("QfwLosaJKVANmUJk", constants.WEBHOOKS_COLLECTION)
	if err != nil {
		t.Errorf("Couldn't get document from firestore: " + err.Error())
	}

	assert.Equal(t, data["FpLSjFbcXoEFfRsW"], webhook1, "Webhook 1 not equal")
	assert.Equal(t, data["QfwLosaJKVANmUJk"], webhook2, "Webhook 2 not equal")

	collection, err := GetAllDocumentInCollectionFromFirestore(constants.WEBHOOKS_COLLECTION)
	if err != nil {
		t.Errorf("Couldn't get collection from firestore: " + err.Error())
	}

	assert.Equal(t, 2, len(collection), "Collection length not equal")
	assert.Equal(t, data["FpLSjFbcXoEFfRsW"], collection["FpLSjFbcXoEFfRsW"], "Webhook 1 not equal")
	assert.Equal(t, data["QfwLosaJKVANmUJk"], collection["QfwLosaJKVANmUJk"], "Webhook 2 not equal")

	assert.True(t, DocumentInCollection("FpLSjFbcXoEFfRsW", constants.WEBHOOKS_COLLECTION), "Webhook 1 not in collection")
	assert.True(t, DocumentInCollection("QfwLosaJKVANmUJk", constants.WEBHOOKS_COLLECTION), "Webhook 2 not in collection")

	DeleteDocument("FpLSjFbcXoEFfRsW", constants.WEBHOOKS_COLLECTION)
	assert.False(t, DocumentInCollection("FpLSjFbcXoEFfRsW", constants.WEBHOOKS_COLLECTION), "Webhook 1 still in collection")
}
