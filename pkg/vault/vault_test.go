package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVaultClient(t *testing.T) {

	// Call the NewVaultClient function.
	client := NewVaultClient()

	// Check that the client is not nil.
	assert.NotNil(t, client)
}

func TestGetMongoDBSecret(t *testing.T) {
	// Create a new Vault client. Make sure VAULT_ADDR and VAULT_TOKEN are properly set.
	client := NewVaultClient()

	// Call the GetMongoDBSecret function.
	secrets := GetMongoDBSecret(client)

	// Check the results.
	// Here, we are just checking that the expected keys exist. You may also wish to check the values.
	assert.NotEmpty(t, secrets["username"], "Expected username to be present")
	assert.NotEmpty(t, secrets["password"], "Expected password to be present")

}
