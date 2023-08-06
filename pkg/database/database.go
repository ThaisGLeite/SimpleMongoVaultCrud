package database

import (
	"context"
	"errors"
	"fmt"
	"simplecrud/pkg/models"
	pkguser "simplecrud/pkg/user"
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
	timeout         = 10 * time.Second
	usersCollection = "users"
	ErrInvalidID    = "invalid id"
	ErrDBConnection = "failed to connect to MongoDB"
)

type UserRepository struct {
	client     *mongo.Client
	database   string
	collection string
	validate   *validator.Validate
}

func NewUserRepository(client *mongo.Client, database string) pkguser.Repository {
	return &UserRepository{
		client:     client,
		database:   database,
		collection: usersCollection,
		validate:   validator.New(),
	}
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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
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

func (r *UserRepository) FindById(ctx context.Context, id string) (models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", ErrInvalidID, err)
	}

	var user models.User
	collection := r.client.Database(r.database).Collection(r.collection)
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		// Check if the error is a "not found" error.
		if err == mongo.ErrNoDocuments {
			return models.User{}, pkguser.ErrNotFound
		}
		return models.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	err := r.validate.Struct(user)
	if err != nil {
		return user, err
	}
	collection := r.client.Database(r.database).Collection(r.collection)
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, user models.User) (models.User, error) {
	err := r.validate.Struct(user)
	if err != nil {
		return user, err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, errors.New(ErrInvalidID)
	}

	updateMap := make(bson.M)
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
	if err != nil {
		return models.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrInvalidID, err)
	}
	collection := r.client.Database(r.database).Collection(r.collection)
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	var users []models.User
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return users, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err = cursor.Decode(&user); err != nil {
			return users, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}

	if err = cursor.Err(); err != nil {
		return users, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, nil
}
