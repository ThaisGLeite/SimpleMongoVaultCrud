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

type userServiceMock struct {
	mock.Mock
}

func (m *userServiceMock) GetAllUsers(c context.Context) ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *userServiceMock) GetUser(c context.Context, id string) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *userServiceMock) CreateUser(c context.Context, user models.User) (models.User, error) {
	args := m.Called(c, user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *userServiceMock) UpdateUser(c context.Context, id string, user models.User) (models.User, error) {
	args := m.Called(id, user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *userServiceMock) DeleteUser(c context.Context, id string) error {
	args := m.Called(c, id) // Pass both parameters here
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
