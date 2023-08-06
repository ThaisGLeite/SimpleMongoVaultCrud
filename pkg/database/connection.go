package database

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/vault/api"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	maxRetries      = 5
	initialInterval = 5 * time.Second
	maxInterval     = 60 * time.Second
)

func ConnectWithRetries(vaultClient *api.Client) (*mongo.Client, string, error) {
	var mongoClient *mongo.Client
	var dbName string
	var err error

	interval := initialInterval
	for i := 1; i <= maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), interval)
		defer cancel()

		mongoClient, dbName, err = ConnectDB(ctx, vaultClient)
		if err == nil {
			log.Printf("Connected to MongoDB! Database name: %s\n", dbName)
			break
		}

		log.Printf("Failed to connect to database (attempt %d of %d): %v. Retrying in %v...\n", i, maxRetries, err, interval)
		time.Sleep(interval)

		interval *= 2
		if interval > maxInterval {
			interval = maxInterval
		}
	}

	return mongoClient, dbName, err
}
