package vault

import (
	"encoding/json"
	"simplecrud/utils"

	vault "github.com/hashicorp/vault/api"
)

func NewVaultClient() *vault.Client {
	vaultAddr := utils.GetEnv("VAULT_ADDR", "http://localhost:8200")
	utils.HandleError("VAULT_ADDR not set", nil)

	vaultConfig := &vault.Config{
		Address: vaultAddr,
	}

	vaultClient, err := vault.NewClient(vaultConfig)
	utils.HandleError("Failed to create Vault client", err)

	vaultToken := utils.GetEnv("VAULT_TOKEN", "")
	utils.HandleError("VAULT_TOKEN not set", nil)

	vaultClient.SetToken(vaultToken)

	return vaultClient
}

func GetMongoDBSecret(vaultClient *vault.Client) map[string]string {
	secretValues, err := vaultClient.Logical().Read("secret/data/mongodb")
	utils.HandleError("Failed to read secret", err)

	if secretValues == nil || secretValues.Data == nil || secretValues.Data["data"] == nil {
		utils.HandleError("No data in secret", nil)
	}

	data, ok := secretValues.Data["data"].(string)
	if !ok {
		utils.HandleError("Data is not a string", nil)
	}

	var mongodbCredentials map[string]string
	err = json.Unmarshal([]byte(data), &mongodbCredentials)
	utils.HandleError("Failed to unmarshal mongodbCredentials", err)

	return mongodbCredentials
}
