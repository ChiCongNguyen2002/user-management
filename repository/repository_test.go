package repository

import (
	"User-Management/models"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	user := models.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password",
	}

	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs(user.Name, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	createdUser, err := userRepo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, createdUser.ID)
}

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
	mock.ExpectPrepare("SELECT id, name, email FROM users LIMIT \\? OFFSET \\?")
	mock.ExpectQuery("SELECT id, name, email FROM users LIMIT \\? OFFSET \\?").WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "John Doe", "john.doe@example.com"))

	users, _, err := repo.GetUsers(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "John Doe", users[0].Name)
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "User 1", "user1@example.com")
	mock.ExpectPrepare("SELECT id, name, email FROM users WHERE id = ?").
		ExpectQuery().
		WithArgs(1).
		WillReturnRows(rows)

	user, err := userRepo.GetUserByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "User 1", user.Name)
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "Updated User", "updated@example.com")

	mock.ExpectPrepare("UPDATE users SET name = \\?, email = \\? WHERE id = \\?").
		ExpectExec().
		WithArgs("Updated User", "updated@example.com", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare("SELECT id, name, email FROM users WHERE id = ?").
		ExpectQuery().
		WithArgs(1).
		WillReturnRows(rows)

	newUserData := models.UserRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	updatedUser, err := userRepo.UpdateUser(context.Background(), 1, newUserData)
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedUser.ID)
	assert.Equal(t, newUserData.Name, updatedUser.Name)
	assert.Equal(t, newUserData.Email, updatedUser.Email)
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	mock.ExpectPrepare("DELETE FROM users WHERE id = ?").
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = userRepo.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)
}
