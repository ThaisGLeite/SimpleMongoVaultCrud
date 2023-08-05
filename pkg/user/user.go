package user

import (
	"context"
	"errors"
	"regexp"
	"simplecrud/pkg/models"
	"simplecrud/utils"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	UpdateUser(ctx context.Context, id string, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserService struct {
	userRepo Repository
}

type Repository interface {
	FindAll(ctx context.Context) ([]models.User, error)
	FindById(ctx context.Context, id string) (models.User, error)
	Create(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, id string, user models.User) (models.User, error)
	Delete(ctx context.Context, id string) error
}

var isAlpha = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

func NewService(userRepo Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Implement the Service interface
func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.FindAll(ctx)
}

// Implement the Service interface
func (s *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	// Check if id is valid UUID
	if !isValidObjectId(id) {
		return models.User{}, errors.New("invalid user ID")
	}

	return s.userRepo.FindById(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	// Check if user name is valid
	if user.Name == "" || len(user.Name) > 50 || !isAlpha(user.Name) {
		return models.User{}, errors.New("invalid user name")
	}

	// Check if password is valid
	if user.Password == "" {
		return models.User{}, errors.New("password cannot be empty")
	}

	// Check if password is strong
	if !utils.IsStrongPassword(user.Password) {
		return models.User{}, errors.New("password isn't strong enough, it should have at least 8 characters, one uppercase letter, one lowercase letter, one number and one special character")
	}

	// Hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashedPassword)

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user models.User) (models.User, error) {
	return s.userRepo.Update(ctx, id, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

func isValidObjectId(id string) bool {
	var objectIdRegex = regexp.MustCompile(`^[0-9a-fA-F]{24}$`)
	return objectIdRegex.MatchString(id)
}
