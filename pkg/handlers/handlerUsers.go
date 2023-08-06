package handlers

import (
	"errors"
	"net/http"
	"simplecrud/pkg/models"
	user "simplecrud/pkg/user"

	"github.com/gin-gonic/gin"
)

// UserHandler struct holds a userService for user operations
type UserHandler struct {
	userService user.Service // Assuming a Service interface is defined in the user package
}

// NewUserHandler initializes a new UserHandler
func NewUserHandler(userService user.Service) *UserHandler {
	// Returns a new UserHandler with the provided userService
	return &UserHandler{
		userService: userService,
	}
}

// GetAllUsers handles the HTTP request to fetch all users
func (u *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := u.userService.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error. " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser handles the HTTP request to fetch a user by ID
func (u *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required in the request"})
		return
	}

	usuario, err := u.userService.GetUser(c, id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error. " + err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, usuario)
}

// CreateUser handles the HTTP request to create a new user
func (u *UserHandler) CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body. " + err.Error()})
		return
	}

	_, err := u.userService.CreateUser(c, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error. " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// UpdateUser handles the HTTP request to update a user
func (u *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required in the request"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body. " + err.Error()})
		return
	}

	_, err := u.userService.UpdateUser(c, id, updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error. " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser handles the HTTP request to delete a user
func (u *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required in the request"})
		return
	}

	err := u.userService.DeleteUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error. " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "User deleted successfully"})
}
