package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"simplecrud/pkg/database"
	"simplecrud/pkg/vault"
	"simplecrud/pkg/web"
)

func main() {
	// Create a new client for interacting with Vault.
	vaultClient := vault.NewVaultClient()

	// Connect to MongoDB using credentials retrieved from Vault.
	mongoClient, dbName, err := database.ConnectWithRetries(vaultClient)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize the user repository.
	userRepo := database.NewUserRepository(mongoClient, dbName)

	// Start the web server in a goroutine so we can listen for shutdown signals.
	go web.StartServer(userRepo)

	// Listen for termination signals.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// Disconnect the MongoDB client.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = mongoClient.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to disconnect from database: %v", err)
	}

	log.Println("Shutdown complete.")
}
