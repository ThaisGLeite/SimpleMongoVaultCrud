package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"simplecrud/pkg/database"
	"simplecrud/pkg/handlers"
	"simplecrud/pkg/models"
	"simplecrud/pkg/user"
	"testing"
	"time"

	mim "github.com/ONSdigital/dp-mongodb-in-memory"
	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func startInMemoryMongoDB() (*mim.Server, string, error) {
	// Start in-memory MongoDB using dp-mongodb-in-memory
	mongoServer, err := mim.Start(context.Background(), "5.0.2") // Change the version if needed
	if err != nil {
		return nil, "", err
	}

	mongoURI := mongoServer.URI()

	return mongoServer, mongoURI, nil
}

func TestIntegrationCRUDOperations(t *testing.T) {
	// Start in-memory MongoDB
	mongoServer, mongoURI, err := startInMemoryMongoDB()
	if err != nil {
		t.Fatalf("Could not start in-memory MongoDB: %v", err)
	}
	defer mongoServer.Stop(context.Background()) // Make sure to stop the server at the end of the test

	// Initialize your app's database connection using mongoURI
	db, err := InitializeDatabaseConnection(mongoURI) // You'll need to implement or modify this function
	if err != nil {
		t.Fatalf("Could not connect to in-memory MongoDB: %v", err)
	}

	dbName := "testdb" // you can specify the database name here
	userRepo := database.NewUserRepository(db, dbName)

	// Set up the router specifically for the test
	router := setupTestRouter(userRepo)

	// Test Create
	userToCreate := models.User{
		Name:     "JohnDoe",
		Age:      25,
		Email:    "john.doe@example.com",
		Password: "P@ssword123",
		Address:  "123 Main St",
	}
	userJSON, _ := json.Marshal(userToCreate)
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userJSON))
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusCreated, resp.Code)
	var createdUser models.User
	_ = json.Unmarshal(resp.Body.Bytes(), &createdUser)

	// Test Get All Users
	req, err = http.NewRequest(http.MethodGet, "/users", nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var users []models.User
	_ = json.Unmarshal(resp.Body.Bytes(), &users)
	// Assuming the newly created user is the last one, or you can find it with specific logic
	createdUser = users[len(users)-1]

	// Test Get by ID
	req, err = http.NewRequest(http.MethodGet, "/users/"+createdUser.ID.Hex(), nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	// Test Update
	userToUpdate := models.User{
		Name: "JaneDoe",
	}
	userJSON, _ = json.Marshal(userToUpdate)
	req, err = http.NewRequest(http.MethodPut, "/users/"+createdUser.ID.Hex(), bytes.NewBuffer(userJSON))
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	// Test Delete
	req, err = http.NewRequest(http.MethodDelete, "/users/"+createdUser.ID.Hex(), nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusNoContent, resp.Code)

	// Clean up
	if err := db.Disconnect(context.Background()); err != nil {
		t.Fatalf("Could not close database connection: %v", err)
	}
}

func InitializeDatabaseConnection(mongoURI string) (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func setupTestRouter(userRepo user.Repository) *gin.Engine {
	// Create a new user service with the provided repository.
	userService := user.NewService(userRepo)
	// Create a new user handler with the created user service.
	userHandler := handlers.NewUserHandler(userService)

	// Create a new gin engine.
	r := gin.New()

	// Create a new rate limiter. This will limit to 1 request/second.
	limiter := tollbooth.NewLimiter(1, nil)

	// Setup the routes for the server.
	setupRoutes(r, limiter, userHandler)

	return r
}
