package user

import (
	"context"
	"errors"
	"regexp"
	"simplecrud/pkg/models"
	"simplecrud/utils"

	"golang.org/x/crypto/bcrypt"
)

// Define the Service interface for user operations.
type Service interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	UpdateUser(ctx context.Context, id string, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

// UserService implements the Service interface.
type UserService struct {
	userRepo Repository
}

// Define the Repository interface for database operations.
type Repository interface {
	FindAll(ctx context.Context) ([]models.User, error)
	FindById(ctx context.Context, id string) (models.User, error)
	Create(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, id string, user models.User) (models.User, error)
	Delete(ctx context.Context, id string) error
}

// Regular expressions to validate user name and id.
// Precompile the regular expressions for efficiency.
var (
	isAlpha         = regexp.MustCompile(`^[a-zA-Z]+$`)
	isValidObjectId = regexp.MustCompile(`^[0-9a-fA-F]{24}$`)
)

// NewService creates a new UserService with the provided user repository.
func NewService(userRepo Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetAllUsers retrieves all users from the repository.
func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.FindAll(ctx)
}

// GetUser retrieves a user by ID from the repository.
// It checks if the provided ID is a valid UUID.
func (s *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	if !isValidObjectId.MatchString(id) {
		return models.User{}, errors.New("invalid user ID")
	}
	return s.userRepo.FindById(ctx, id)
}

// CreateUser creates a new user in the repository.
// It validates the user name and password before creation, and hashes the password.
func (s *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	if user.Name == "" || len(user.Name) > 50 || !isAlpha.MatchString(user.Name) {
		return models.User{}, errors.New("invalid user name")
	}
	if user.Password == "" {
		return models.User{}, errors.New("password cannot be empty")
	}
	if !utils.IsStrongPassword(user.Password) {
		return models.User{}, errors.New("password isn't strong enough, it should have at least 8 characters, one uppercase letter, one lowercase letter, one number and one special character")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashedPassword)
	return s.userRepo.Create(ctx, user)
}

// UpdateUser updates a user by ID in the repository.
func (s *UserService) UpdateUser(ctx context.Context, id string, user models.User) (models.User, error) {
	return s.userRepo.Update(ctx, id, user)
}

// DeleteUser deletes a user by ID from the repository.
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
