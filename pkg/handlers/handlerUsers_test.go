package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"simplecrud/pkg/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// userServiceMock struct is used for mocking user service
type userServiceMock struct {
	mock.Mock
}

// GetAllUsers mocks the function to get all users
func (m *userServiceMock) GetAllUsers(c context.Context) ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

// GetUser mocks the function to get a single user by ID
func (m *userServiceMock) GetUser(c context.Context, id string) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

// CreateUser mocks the function to create a new user
func (m *userServiceMock) CreateUser(c context.Context, user models.User) (models.User, error) {
	args := m.Called(c, user)
	return args.Get(0).(models.User), args.Error(1)
}

// UpdateUser mocks the function to update a user by ID
func (m *userServiceMock) UpdateUser(c context.Context, id string, user models.User) (models.User, error) {
	args := m.Called(id, user)
	return args.Get(0).(models.User), args.Error(1)
}

// DeleteUser mocks the function to delete a user by ID
func (m *userServiceMock) DeleteUser(c context.Context, id string) error {
	args := m.Called(c, id)
	return args.Error(0)
}

func TestGetAllUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(userServiceMock)
	userHandler := NewUserHandler(mockUserService)

	objectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	users := []models.User{
		{ID: objectID, Name: "JohnDoe", Age: 23, Email: "teste@teste.com.br", Password: "P@ssword7", Address: "Rua Romao Batista"},
	}

	// Success case
	mockUserService.On("GetAllUsers").Return(users, nil)
	router := gin.Default()
	router.GET("/users", userHandler.GetAllUsers)
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	mockUserService.AssertExpectations(t)

	// Error case
	// Handling other cases like when an error occurs
	mockUserService = new(userServiceMock) // Create a new mock instance
	userHandler = NewUserHandler(mockUserService)
	mockUserService.On("GetAllUsers").Return([]models.User{}, errors.New("Internal Error"))
	router = gin.Default() // Create a new router instance
	router.GET("/users", userHandler.GetAllUsers)
	response = httptest.NewRecorder()
	request, _ = http.NewRequest(http.MethodGet, "/users", nil)
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	mockUserService.AssertExpectations(t)
}

// TestGetUserSuccess defines the tests for a successful GetUser request
func TestGetUserSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(userServiceMock)
	userHandler := NewUserHandler(mockUserService)
	objectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	user := models.User{ID: objectID, Name: "JohnDoe"}

	mockUserService.On("GetUser", objectID.Hex()).Return(user, nil)
	router := gin.Default()
	router.GET("/users/:id", userHandler.GetUser)
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/"+objectID.Hex(), nil)
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	mockUserService.AssertExpectations(t)
}

// TestGetUserError defines the tests for an error in a GetUser request
func TestGetUserError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(userServiceMock)
	userHandler := NewUserHandler(mockUserService)
	objectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	mockUserService.On("GetUser", objectID.Hex()).Return(models.User{}, errors.New("Not Found"))
	router := gin.Default()
	router.GET("/users/:id", userHandler.GetUser)
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/"+objectID.Hex(), nil)
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	mockUserService.AssertExpectations(t)
}

// TestCreateUserSuccess defines the tests for a successful CreateUser request
func TestCreateUserSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(userServiceMock)
	userHandler := NewUserHandler(mockUserService)

	user := models.User{
		Name:     "JohnDoe",
		Age:      23,
		Email:    "teste@teste.com.br",
		Password: "P@sswoooord7",
		Address:  "Rua Romao Batista",
	}

	payload, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	body := bytes.NewReader(payload)

	mockUserService.On("CreateUser", mock.Anything, user).Return(user, nil)
	router := gin.Default()
	router.POST("/users", userHandler.CreateUser)
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", body)
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusCreated, response.Code)
	mockUserService.AssertExpectations(t)
}

// TestCreateUserError defines the tests for an error in a CreateUser request
func TestCreateUserError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(userServiceMock)
	userHandler := NewUserHandler(mockUserService)

	user := models.User{
		Name:     "JohnDoe",
		Age:      23,
		Email:    "teste@teste.com.br",
		Password: "P@sswoooord7",
		Address:  "Rua Romao Batista",
	}

	payload, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	body := bytes.NewReader(payload)

	mockUserService.On("CreateUser", mock.Anything, user).Return(models.User{}, errors.New("Creation Failed"))
	router := gin.Default()
	router.POST("/users", userHandler.CreateUser)
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", body)
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	mockUserService.AssertExpectations(t)
}

// TestUpdateUserSuccess defines the tests for a successful UpdateUser request
func TestUpdateUserSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(userServiceMock)
	userHandler := NewUserHandler(mockUserService)
	objectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	user := models.User{ID: objectID, Name: "UpdatedUser"}

	payload, _ := json.Marshal(user)
	body := bytes.NewReader(payload)

	mockUserService.On("UpdateUser", mock.AnythingOfType("string"), mock.AnythingOfType("models.User")).Return(user, nil)
	router := gin.Default()
	router.PUT("/users/:id", userHandler.UpdateUser)
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, "/users/"+objectID.Hex(), body)
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	mockUserService.AssertExpectations(t)
}

// TestUpdateUserError defines the tests for an error in an UpdateUser request
func TestUpdateUserError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(userServiceMock)
	userHandler := NewUserHandler(mockUserService)
	objectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	user := models.User{ID: objectID, Name: "UpdatedUser"}

	payload, _ := json.Marshal(user)
	body := bytes.NewReader(payload)

	mockUserService.On("UpdateUser", mock.AnythingOfType("string"), mock.AnythingOfType("models.User")).Return(models.User{}, errors.New("Update Failed"))
	router := gin.Default()
	router.PUT("/users/:id", userHandler.UpdateUser)
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, "/users/"+objectID.Hex(), body)
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	mockUserService.AssertExpectations(t)
}

// TestDeleteUser defines the tests for deleting a user including success and error cases
func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(userServiceMock)
	userHandler := NewUserHandler(mockUserService)
	objectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	// Success case
	mockUserService.On("DeleteUser", mock.Anything, objectID.Hex()).Return(nil)
	router := gin.Default()
	router.DELETE("/users/:id", userHandler.DeleteUser)
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, "/users/"+objectID.Hex(), nil)
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusNoContent, response.Code)
	mockUserService.AssertExpectations(t)

	// Error case
	mockUserService = new(userServiceMock)
	userHandler = NewUserHandler(mockUserService)
	mockUserService.On("DeleteUser", mock.Anything, objectID.Hex()).Return(errors.New("Delete Failed"))
	router = gin.Default()
	router.DELETE("/users/:id", userHandler.DeleteUser)
	response = httptest.NewRecorder()
	request, _ = http.NewRequest(http.MethodDelete, "/users/"+objectID.Hex(), nil)
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusInternalServerError, response.Code) // Assuming you return 500 for delete failure
	mockUserService.AssertExpectations(t)
}
