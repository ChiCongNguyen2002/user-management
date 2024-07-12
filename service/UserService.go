package service

import (
	"User-Management/models"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUsers(ctx context.Context, page int, pageSize int) ([]models.UserRequest, int, error)
	GetUserByID(ctx context.Context, id int) (models.UserRequest, error)
	UpdateUser(ctx context.Context, id int, user models.UserRequest) (models.UserRequest, error)
	DeleteUser(ctx context.Context, id int) error
}
