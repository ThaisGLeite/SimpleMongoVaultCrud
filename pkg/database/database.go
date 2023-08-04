package database

import (
	"context"
	"fmt"
	"log"
	"simplecrud/pkg/vault"
	"simplecrud/utils"
	"time"

	"github.com/hashicorp/vault/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(vaultClient *api.Client) (*mongo.Client, string) {
	secretValues := vault.GetMongoDBSecret(vaultClient)

	username, userOk := secretValues["username"]
	password, passOk := secretValues["password"]

	if !userOk || !passOk {
		log.Fatal("Username or password not found in secret")
	}

	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "27017")
	dbName := utils.GetEnv("DB_NAME", "devenv")

	if dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("DB_HOST, DB_PORT or DB_NAME not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, dbHost, dbPort)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	utils.HandleError("Failed to connect to MongoDB", err)

	return client, dbName
}
