package service

import (
	"User-Management/models"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockUserRepository struct{}

func (m *MockUserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	return user, nil
}

func (m *MockUserRepository) GetUsers(ctx context.Context, page, pageSize int) ([]models.UserRequest, int, error) {
	users := []models.UserRequest{
		{ID: 1, Name: "User 1", Email: "user1@example.com"},
		{ID: 2, Name: "User 2", Email: "user2@example.com"},
	}
	return users, len(users), nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (models.UserRequest, error) {
	if id == 1 {
		return models.UserRequest{ID: 1, Name: "User 1", Email: "user1@example.com"}, nil
	}
	return models.UserRequest{}, errors.New("user not found")
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, id int, user models.UserRequest) (models.UserRequest, error) {
	if id == 1 {
		user.ID = id
		return user, nil
	}
	return models.UserRequest{}, errors.New("user not found")
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) error {
	if id == 1 {
		return nil
	}
	return errors.New("user not found")
}

func TestCreateUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userService := NewUserService(mockRepo)

	user := models.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password",
	}

	createdUser, err := userService.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, createdUser.Name)
}

func TestGetUsers(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userService := NewUserService(mockRepo)

	users, total, err := userService.GetUsers(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, users, 2)
	assert.Equal(t, "User 1", users[0].Name)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userService := NewUserService(mockRepo)

	user, err := userService.GetUserByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "User 1", user.Name)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userService := NewUserService(mockRepo)

	newUserData := models.UserRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	updatedUser, err := userService.UpdateUser(context.Background(), 1, newUserData)
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedUser.ID)
	assert.Equal(t, newUserData.Name, updatedUser.Name)
	assert.Equal(t, newUserData.Email, updatedUser.Email)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userService := NewUserService(mockRepo)

	err := userService.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)
}
