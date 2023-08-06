package database

import (
	"context"
	"errors"
	"os"
	"simplecrud/pkg/models"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

// setup initializes the database connection by connecting to Vault and MongoDB
func setup() (*mongo.Client, string, error) {
	vaultAddress := os.Getenv("VAULT_ADDR")
	if vaultAddress == "" {
		return nil, "", errors.New("VAULT_ADDR environment variable not set") // Ensure the Vault address is set
	}
	vaultClient, err := api.NewClient(&api.Config{Address: vaultAddress})
	if err != nil {
		return nil, "", err
	}
	return ConnectDB(context.Background(), vaultClient)
}

func TestCRUDOperations(t *testing.T) {
	client, dbName, err := setup() // Setup database connection
	require.NoError(t, err)
	repo := NewUserRepository(client, dbName)

	// Test Create
	// Prepare a user object to be used in the CRUD operations test
	user := models.User{
		Name:     "JohnDoe",
		Age:      25,
		Email:    "john.doe@example.com",
		Password: "P@ssword123",
		Address:  "123 Main St",
	}

	// Test the Create operation
	_, err = repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Find all users and get the ID of the user with the same name as the created user
	allUsers, err := repo.FindAll(context.Background())
	require.NoError(t, err)
	var userIDFromDatabase string
	for _, u := range allUsers {
		if u.Name == user.Name {
			userIDFromDatabase = u.ID.Hex() // Store the user's ID for subsequent tests
			break
		}
	}

	// Test FindById
	// Fetch the user by ID and check if it matches the created user
	foundUser, err := repo.FindById(context.Background(), userIDFromDatabase)
	require.NoError(t, err)
	require.Equal(t, userIDFromDatabase, foundUser.ID.Hex())

	// Test Update
	// Prepare updated user data
	updatedUser := models.User{
		Address: "Aqui em casa",
	}

	// Test the Update operation
	_, err = repo.Update(context.Background(), userIDFromDatabase, updatedUser)
	require.NoError(t, err)

	// Test Delete
	// Test the Delete operation by removing the user
	err = repo.Delete(context.Background(), userIDFromDatabase)
	require.NoError(t, err)
}
