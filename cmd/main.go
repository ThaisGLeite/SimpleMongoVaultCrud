package main

import (
	"context"
	"fmt"

	"simplecrud/pkg/database"
	"simplecrud/pkg/vault"
	"simplecrud/pkg/web"
)

func main() {
	vaultClient := vault.NewVaultClient()

	// Connect to MongoDB
	mongoClient, dbName := database.ConnectDB(vaultClient)
	defer mongoClient.Disconnect(context.Background())
	fmt.Printf("Connected to MongoDB! Database name: %s\n", dbName)

	// Create a MongoDB repository
	userRepo := database.NewUserRepository(mongoClient, dbName)

	// Start the server
	web.StartServer(userRepo) // pass the userRepo as a parameter to StartServer
}
