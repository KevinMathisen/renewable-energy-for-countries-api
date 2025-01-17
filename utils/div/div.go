﻿package div

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
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Define possible letters in webhookID
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// Create string where each character is chosen at random
	b := make([]rune, constants.WEBHOOK_ID_LENGTH)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

/*
Returns if slice contains value
*/
func Contains(slice []string, value string) bool {
	for _, valueInSlice := range slice {
		if valueInSlice == value {
			return true
		}
	}
	return false
}

/*
Removes duplicate string values from an array before returning them
*/
func RemoveDuplicates(arr []string) []string {
	// Create map to keep track of encountered strings
	encountered := map[string]bool{}
	result := []string{}

	// Add all strings from array to encountered map
	for _, value := range arr {
		if !encountered[value] { // If string has not been encountered yet
			encountered[value] = true      // Add string to encountered map
			result = append(result, value) // Add string to result
		}
	}

	return result
}
