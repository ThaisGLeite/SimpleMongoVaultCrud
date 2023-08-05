package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"simplecrud/pkg/models"
	"simplecrud/pkg/user"
	"simplecrud/pkg/vault"
	"simplecrud/utils"
	"time"

	"github.com/go-playground/validator"
	"github.com/hashicorp/vault/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// timeout duration in seconds
	timeout = 10 * time.Second
)

type UserRepository struct {
	client     *mongo.Client
	database   string
	collection string
	validate   *validator.Validate
}

func NewUserRepository(client *mongo.Client, database string) user.Repository {
	return &UserRepository{
		client:     client,
		database:   database,
		collection: "users",
		validate:   validator.New(),
	}
}

func ConnectDB(vaultClient *api.Client) (*mongo.Client, string, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, dbHost, dbPort)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return client, dbName, nil
}

func (r *UserRepository) FindById(ctx context.Context, id string) (models.User, error) {
	// Convert string to ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid id: %w", err)
	}

	var user models.User
	collection := r.client.Database(r.database).Collection(r.collection)
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)

	return user, err
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	err := r.validate.Struct(user)
	if err != nil {
		return user, err
	}
	collection := r.client.Database(r.database).Collection(r.collection)
	_, err = collection.InsertOne(ctx, user)

	return user, err
}

func (r *UserRepository) Update(ctx context.Context, id string, user models.User) (models.User, error) {
	err := r.validate.Struct(user)
	if err != nil {
		return user, err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, errors.New("invalid ID")
	}

	// Create a map to hold the fields to update
	updateMap := make(bson.M)

	// Check each field in the User struct and add it to the update map if it is not empty
	if user.Name != "" {
		updateMap["name"] = user.Name
	}
	if user.Age != 0 {
		updateMap["age"] = user.Age
	}
	if user.Email != "" {
		updateMap["email"] = user.Email
	}
	if user.Password != "" {
		updateMap["password"] = user.Password
	}
	if user.Address != "" {
		updateMap["address"] = user.Address
	}

	collection := r.client.Database(r.database).Collection(r.collection)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": updateMap,
	})

	return user, err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// Convert string to ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	collection := r.client.Database(r.database).Collection(r.collection)
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})

	return err
}

func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	// Create an empty slice to store the decoded users
	var users []models.User

	// Use the collection's Find method to get all documents
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return users, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document one at a time
	for cursor.Next(ctx) {
		var user models.User
		if err = cursor.Decode(&user); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	// Check if the cursor encountered any errors while iterating
	if err = cursor.Err(); err != nil {
		return users, err
	}

	return users, nil
}
