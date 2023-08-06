package main

import (
	"context"
	"log"
	"time"

	"simplecrud/pkg/database"
	"simplecrud/pkg/vault"
	"simplecrud/pkg/web"
)

func main() {
	// Create a new client for interacting with Vault.
	vaultClient := vault.NewVaultClient()

	// Connect to MongoDB using credentials retrieved from Vault.
	// The connection is attempted with retries, adhering to an exponential backoff strategy.
	mongoClient, dbName, err := database.ConnectWithRetries(vaultClient)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize the user repository using the connected MongoDB client.
	// This repository provides an interface to interact with user-related data in the database.
	userRepo := database.NewUserRepository(mongoClient, dbName)

	// Start the web server, providing the user repository for handling user-related HTTP requests.
	web.StartServer(userRepo)

	// Prepare for MongoDB client disconnection with a 10-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Disconnect the MongoDB client.
	// Proper disconnection ensures that all connections to the database are cleanly closed.
	if err = mongoClient.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to disconnect from database: %v", err)
	}
}
