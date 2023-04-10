package db

import (
	"context"

	"cloud.google.com/go/firestore" // Firestore-specific support
)

// Firebase context used by Firestore functions
var firestoreContext context.Context

// Firebase client used by Firestore functions
var firebaseClient *firestore.Client
