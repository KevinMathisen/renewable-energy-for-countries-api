package div

import (
	"assignment2/utils/constants"
	"math/rand"
	"time"
)

/*
Create random WebhookID with length 16

	return - Random webhookID
*/
func CreateWebhookId() string {
	// Set seed for rand
	rand.Seed(time.Now().UnixNano())

	// Define possible letters in webhookID
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// Create string where each character is chosen at random
	b := make([]rune, constants.WEBHOOK_ID_LENGTH)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
