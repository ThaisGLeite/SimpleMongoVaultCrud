package user

import (
	"context"
	"simplecrud/pkg/models"
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

func NewService(userRepo Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Implement the Service interface
func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.FindAll(ctx)
}

func (s *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	return s.userRepo.FindById(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user models.User) (models.User, error) {
	return s.userRepo.Update(ctx, id, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
