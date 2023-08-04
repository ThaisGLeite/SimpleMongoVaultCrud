package database

import (
	"context"
	"fmt"
	"log"
	"simplecrud/pkg/models"
	"simplecrud/pkg/user"
	"simplecrud/pkg/vault"
	"simplecrud/utils"
	"time"

	"github.com/hashicorp/vault/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewUserRepository(client *mongo.Client, database string) user.Repository {
	return &UserRepository{
		client:     client,
		database:   database,
		collection: "users", // this could be any collection you want to store users in
	}
}

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

func (r *UserRepository) FindById(ctx context.Context, id string) (models.User, error) {
	// Convert string to ObjectId
	objID, _ := primitive.ObjectIDFromHex(id)

	var user models.User
	collection := r.client.Database(r.database).Collection(r.collection)
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)

	return user, err
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	_, err := collection.InsertOne(ctx, user)

	return user, err
}

func (r *UserRepository) Update(ctx context.Context, id string, user models.User) (models.User, error) {
	// Convert string to ObjectId
	objID, _ := primitive.ObjectIDFromHex(id)

	collection := r.client.Database(r.database).Collection(r.collection)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{
			"name":     user.Name,
			"age":      user.Age,
			"email":    user.Email,
			"password": user.Password,
			"address":  user.Address,
		},
	})

	return user, err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// Convert string to ObjectId
	objID, _ := primitive.ObjectIDFromHex(id)

	collection := r.client.Database(r.database).Collection(r.collection)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": objID})

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
