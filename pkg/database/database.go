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
	timeout         = 10 * time.Second               // The timeout for database operations
	usersCollection = "users"                        // The MongoDB collection for users
	ErrInvalidID    = "invalid id"                   // Error message for an invalid ID
	ErrDBConnection = "failed to connect to MongoDB" // Error message for connection failure
)

// UserRepository represents the MongoDB repository for user operations
type UserRepository struct {
	client     *mongo.Client       // MongoDB client
	database   string              // MongoDB database name
	collection string              // MongoDB collection name
	validate   *validator.Validate // Validator for user struct
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(client *mongo.Client, database string) pkguser.Repository {
	return &UserRepository{
		client:     client,
		database:   database,
		collection: usersCollection,
		validate:   validator.New(), // Initialize validator for user input validation
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

// FindById finds a user by ID in the MongoDB collection
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

// Create inserts a new user into the MongoDB collection
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

// Update updates a user's details in the MongoDB collection.
func (r *UserRepository) Update(ctx context.Context, id string, user models.User) (models.User, error) {
	// Validate the user struct to ensure it meets the required constraints.
	err := r.validate.Struct(user)
	if err != nil {
		return user, err
	}

	// Convert the string ID to MongoDB's ObjectID type. Return an error if the conversion fails.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, errors.New(ErrInvalidID)
	}

	// Create a map to hold the fields that need to be updated. Only non-empty fields will be added.
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

	// Get the user collection from the MongoDB client.
	collection := r.client.Database(r.database).Collection(r.collection)

	// Perform the update operation using MongoDB's UpdateOne method.
	// The "$set" operator replaces the value of a field with the specified value.
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": updateMap,
	})
	if err != nil {
		// Return an error if the update operation fails.
		return models.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	// Return the updated user object.
	return user, nil
}

// Delete removes a user from the MongoDB collection based on the provided ID.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// Convert the string ID to MongoDB's ObjectID type. Return an error if the conversion fails.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrInvalidID, err)
	}

	// Get the user collection from the MongoDB client.
	collection := r.client.Database(r.database).Collection(r.collection)

	// Perform the delete operation using MongoDB's DeleteOne method, targeting the user with the given ObjectID.
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		// Return an error if the delete operation fails.
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Return nil to indicate successful deletion.
	return nil
}

// FindAll retrieves all users from the MongoDB collection and returns them in a slice.
func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	// Get the user collection from the MongoDB client.
	collection := r.client.Database(r.database).Collection(r.collection)

	// Initialize an empty slice to hold the retrieved users.
	var users []models.User

	// Perform the find operation to retrieve all users from the collection.
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		// Return an error if the find operation fails.
		return users, fmt.Errorf("failed to find users: %w", err)
	}

	// Ensure that the cursor is closed after the function exits.
	defer cursor.Close(ctx)

	// Iterate through the cursor, decoding each user document into a User struct.
	for cursor.Next(ctx) {
		var user models.User
		if err = cursor.Decode(&user); err != nil {
			// Return an error if the decoding fails.
			return users, fmt.Errorf("failed to decode user: %w", err)
		}
		// Append the decoded user to the users slice.
		users = append(users, user)
	}

	// Check if there were any errors during cursor iteration.
	if err = cursor.Err(); err != nil {
		return users, fmt.Errorf("failed to iterate users: %w", err)
	}

	// Return the users slice containing all retrieved users.
	return users, nil
}
