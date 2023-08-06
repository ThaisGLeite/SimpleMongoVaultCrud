package main

import (
	"context"
	"log"
	"time"

	"simplecrud/pkg/database"
	"simplecrud/pkg/vault"
	"simplecrud/pkg/web"

	"github.com/hashicorp/vault/api"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// Constants defining the retry mechanism for connecting to the database
	maxRetries      = 5
	initialInterval = 5 * time.Second
	maxInterval     = 60 * time.Second
)

func main() {
	// Create a new Vault client
	vaultClient := vault.NewVaultClient()

	// Attempt to connect to MongoDB with retries, fetching credentials from Vault
	mongoClient, dbName, err := connectWithRetries(vaultClient)
	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// Create a new user repository using the connected MongoDB client
	userRepo := database.NewUserRepository(mongoClient, dbName)

	// Start the web server
	web.StartServer(userRepo)

	// Disconnect the MongoDB client after the web server has finished running
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = mongoClient.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to disconnect from database: %v", err)
	}
}

// connectWithRetries attempts to connect to the database with exponential backoff
func connectWithRetries(vaultClient *api.Client) (*mongo.Client, string, error) {
	var mongoClient *mongo.Client
	var dbName string
	var err error

	interval := initialInterval
	for i := 1; i <= maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), interval)
		defer cancel()

		// Attempt to connect to the database using the Vault client
		mongoClient, dbName, err = database.ConnectDB(ctx, vaultClient)
		if err == nil {
			log.Printf("Connected to MongoDB! Database name: %s\n", dbName)
			break
		}

		// Log the failure and retry after waiting
		log.Printf("Failed to connect to database (attempt %d of %d): %v. Retrying in %v...\n", i, maxRetries, err, interval)
		time.Sleep(interval)

		// Double the interval for each retry, but don't exceed maxInterval
		interval *= 2
		if interval > maxInterval {
			interval = maxInterval
		}
	}

	return mongoClient, dbName, err
}
