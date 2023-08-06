package database

import (
	"context"
	"simplecrud/pkg/models"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func setup() (*mongo.Client, string, error) {
	vaultClient, err := api.NewClient(&api.Config{Address: "http://vault:8200"})
	if err != nil {
		return nil, "", err
	}
	return ConnectDB(context.Background(), vaultClient)
}
func TestCRUDOperations(t *testing.T) {
	client, dbName, err := setup()
	require.NoError(t, err)
	repo := NewUserRepository(client, dbName)

	// Test Create
	user := models.User{
		Name:     "JohnDoe",
		Age:      25,
		Email:    "john.doe@example.com",
		Password: "P@ssword123",
		Address:  "123 Main St",
	}

	_, err = repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Find all users and get the ID of the user with the same name as the created user
	allUsers, err := repo.FindAll(context.Background())
	require.NoError(t, err)
	var userIDFromDatabase string
	for _, u := range allUsers {
		if u.Name == user.Name {
			userIDFromDatabase = u.ID.Hex()
			break
		}
	}

	// Test FindById
	foundUser, err := repo.FindById(context.Background(), userIDFromDatabase)
	require.NoError(t, err)
	require.Equal(t, userIDFromDatabase, foundUser.ID.Hex())

	// Test Update
	updatedUser := models.User{
		Address: "Aqui em casa",
	}

	_, err = repo.Update(context.Background(), userIDFromDatabase, updatedUser)
	require.NoError(t, err)

	// Test Delete
	err = repo.Delete(context.Background(), userIDFromDatabase)
	require.NoError(t, err)
}
