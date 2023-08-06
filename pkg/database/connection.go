package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"simplecrud/pkg/vault"
	"simplecrud/utils"
	"time"

	"github.com/hashicorp/vault/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Constants for controlling the retry mechanism when connecting to the database.
	maxRetries      = 5                // Maximum number of connection attempts.
	initialInterval = 5 * time.Second  // Initial delay between attempts.
	maxInterval     = 60 * time.Second // Maximum delay between attempts.
)

// ConnectWithRetries attempts to connect to MongoDB using credentials from Vault.
// If the connection attempt fails, it retries up to maxRetries times,
// using an exponential backoff strategy controlled by initialInterval and maxInterval.
func ConnectWithRetries(vaultClient *api.Client) (*mongo.Client, string, error) {
	var mongoClient *mongo.Client
	var dbName string
	var err error

	interval := initialInterval // Set the initial delay between connection attempts.

	for i := 1; i <= maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), interval)
		defer cancel()

		// Attempt to connect to the database using Vault credentials.
		mongoClient, dbName, err = ConnectDB(ctx, vaultClient)
		if err == nil {
			log.Printf("Connected to MongoDB! Database name: %s\n", dbName)
			break
		}

		// Log the failure and prepare for the next attempt.
		log.Printf("Failed to connect to database (attempt %d of %d): %v. Retrying in %v...\n", i, maxRetries, err, interval)
		time.Sleep(interval)

		// Double the delay for the next attempt, without exceeding maxInterval.
		interval *= 2
		if interval > maxInterval {
			interval = maxInterval
		}
	}

	return mongoClient, dbName, err
}

// ConnectDB establishes a connection to the MongoDB database
func ConnectDB(ctx context.Context, vaultClient *api.Client) (*mongo.Client, string, error) {
	// Get the secrets from vault
	secretValues := vault.GetMongoDBSecret(vaultClient)

	username, userOk := secretValues["username"]
	password, passOk := secretValues["password"]

	// Check if username and password are present in the secret
	if !userOk || !passOk {
		return nil, "", errors.New("username or password not found in secret")
	}

	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "27017")
	dbName := utils.GetEnv("DB_NAME", "devenv")

	// Check if DB_HOST, DB_PORT or DB_NAME are not set
	if dbHost == "" || dbPort == "" || dbName == "" {
		return nil, "", errors.New("DB_HOST, DB_PORT or DB_NAME not set")
	}

	// Connect to MongoDB
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, dbHost, dbPort)

	// Connection pooling settings
	connectionOptions := options.Client().
		ApplyURI(connectionString).
		SetMaxPoolSize(100).                 // Maximum number of connections in the connection pool
		SetMinPoolSize(10).                  // Minimum number of connections in the connection pool
		SetMaxConnIdleTime(30 * time.Minute) // Maximum idle time for a connection

	// Connect to MongoDB with the specified settings
	client, err := mongo.Connect(ctx, connectionOptions)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", ErrDBConnection, err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", ErrDBConnection, err)
	}

	return client, dbName, nil
}
