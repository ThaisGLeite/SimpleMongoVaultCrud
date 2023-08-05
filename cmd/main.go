package main

import (
	"context"
	"log"
	"time"

	"simplecrud/pkg/database"
	"simplecrud/pkg/vault"
	"simplecrud/pkg/web"

	"go.mongodb.org/mongo-driver/mongo"
)

const maxRetries = 5

func main() {
	vaultClient := vault.NewVaultClient()

	// Create a MongoDB repository
	var mongoClient *mongo.Client
	var dbName string
	var err error

	for i := 1; i <= maxRetries; i++ {
		// Connect to MongoDB
		mongoClient, dbName, err = database.ConnectDB(vaultClient)
		if err == nil {
			log.Printf("Connected to MongoDB! Database name: %s\n", dbName)
			break
		}
		log.Printf("Failed to connect to database (attempt %d of %d): %v. Retrying in %d seconds...\n", i, maxRetries, err, i*5)
		time.Sleep(time.Duration(i*5) * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts", maxRetries)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from database: %v", err)
		}
	}()

	// Create a MongoDB repository
	userRepo := database.NewUserRepository(mongoClient, dbName)

	// Start the server
	web.StartServer(userRepo) // pass the userRepo as a parameter to StartServer
}
