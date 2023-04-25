package div

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Tests the creation of a random webhookID
*/
func TestCreateWebhookId(t *testing.T) {
	// Create two webhookIDs
	id1 := CreateWebhookId()
	id2 := CreateWebhookId()

	// Check if they are the same
	assert.NotEqual(t, id1, id2, "WebhookIDs are the same")
	assert.Equal(t, len(id1), 16, "WebhookID is not 16 characters long")
	assert.Equal(t, len(id2), 16, "WebhookID is not 16 characters long")
}

/*
Tests the Contains function
*/
func TestContains(t *testing.T) {
	// Create slice
	slice := []string{"a", "b", "c"}

	// Check if slice contains values
	assert.True(t, Contains(slice, "a"), "Slice does not contain value a")
	assert.True(t, Contains(slice, "b"), "Slice does not contain value b")
	assert.True(t, Contains(slice, "c"), "Slice does not contain value c")
	assert.False(t, Contains(slice, "d"), "Slice contains value d")
}
