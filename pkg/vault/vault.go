package vault

import (
	"encoding/json"
	"simplecrud/utils"

	vault "github.com/hashicorp/vault/api"
)

// NewVaultClient function creates and configures a new Vault client.
// It uses the VAULT_ADDR and VAULT_TOKEN environment variables for configuration.
// If these variables are not set, it uses default values.
func NewVaultClient() *vault.Client {
	// Get the Vault address from the environment variable VAULT_ADDR.
	// If VAULT_ADDR is not set, use "http://localhost:8200" as the default address.
	vaultAddr := utils.GetEnv("VAULT_ADDR", "http://localhost:8200")

	// Create a new Vault configuration with the retrieved address.
	vaultConfig := &vault.Config{
		Address: vaultAddr,
	}

	// Create a new Vault client with the configuration.
	// If there is an error, log it and terminate the program.
	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		utils.HandleError("E", "Failed to create Vault client", err)
	}

	// Get the Vault token from the environment variable VAULT_TOKEN.
	// If VAULT_TOKEN is not set, use an empty string as the default token.
	vaultToken := utils.GetEnv("VAULT_TOKEN", "")

	// Set the Vault token to the client.
	vaultClient.SetToken(vaultToken)

	return vaultClient
}

// GetMongoDBSecret function retrieves MongoDB secrets from Vault.
// It returns a map where keys are the secret names and values are the secret values.
func GetMongoDBSecret(vaultClient *vault.Client) map[string]string {
	// Read the secret from Vault at the path "secret/data/mongodb".
	// If there is an error, log it and terminate the program.
	secretValues, err := vaultClient.Logical().Read("secret/data/mongodb")
	if err != nil {
		utils.HandleError("E", "Failed to read secret", err)
	}

	// Check if the secret contains data.
	if secretValues == nil || secretValues.Data == nil || secretValues.Data["data"] == nil {
		utils.HandleError("E", "No data in secret", nil)
	}

	// Extract the data from the secret.
	data, ok := secretValues.Data["data"].(string)
	if !ok {
		utils.HandleError("E", "Data is not a string", nil)
	}

	// Unmarshal the data into a map.
	var mongodbCredentials map[string]string
	err = json.Unmarshal([]byte(data), &mongodbCredentials)
	if err != nil {
		utils.HandleError("E", "Failed to unmarshal mongodbCredentials", err)
	}

	return mongodbCredentials
}
