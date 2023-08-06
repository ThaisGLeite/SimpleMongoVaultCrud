package user

import (
	"context"
	"errors"
	"simplecrud/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockRepository struct {
	Users []models.User
}

func (m *MockRepository) FindAll(ctx context.Context) ([]models.User, error) {
	return m.Users, nil
}

func (m *MockRepository) FindById(ctx context.Context, id string) (models.User, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	for _, user := range m.Users {
		if user.ID == objectID {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (m *MockRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	m.Users = append(m.Users, user)
	return user, nil
}

func (m *MockRepository) Update(ctx context.Context, id string, user models.User) (models.User, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	for i, u := range m.Users {
		if u.ID == objectID {
			m.Users[i] = user
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	objectID, _ := primitive.ObjectIDFromHex(id)
	for i, user := range m.Users {
		if user.ID == objectID {
			m.Users = append(m.Users[:i], m.Users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func TestGetAllUsers(t *testing.T) {
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	mockRepo := &MockRepository{
		Users: []models.User{{ID: id1, Name: "Alice"}, {ID: id2, Name: "Bob"}},
	}
	service := NewService(mockRepo)

	users, err := service.GetAllUsers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
}
func TestGetUser(t *testing.T) {
	id := primitive.NewObjectID() // Only generate the ID once
	mockRepo := &MockRepository{
		Users: []models.User{{ID: id, Name: "Alice"}}, // Use the generated ID here
	}
	service := NewService(mockRepo)

	// Testing with a valid ID
	user, err := service.GetUser(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, id.Hex(), user.ID.Hex()) // Use the generated ID here too

	// Testing with an invalid ID
	user, err = service.GetUser(context.Background(), "invalidID")
	assert.Error(t, err)
}
