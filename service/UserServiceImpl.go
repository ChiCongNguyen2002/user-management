package service

import (
	"User-Management/models"
	"User-Management/repository"
	"context"
	"log"
	"time"
)

type UserServiceImpl struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{UserRepo: userRepo}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	newUser, err := s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		log.Printf("CreateUser service error: %v\n", err)
		return models.User{}, err
	}
	user = newUser
	return user, nil
}

func (s *UserServiceImpl) GetUsers(ctx context.Context, page int, pageSize int) ([]models.UserRequest, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	users, total, err := s.UserRepo.GetUsers(ctx, page, pageSize)
	if err != nil {
		log.Printf("GetUsers service error: %v\n", err)
		return nil, 0, err
	}
	return users, total, nil
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, id int) (models.UserRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := s.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		log.Printf("GetUserByID service error: %v\n", err)
		return models.UserRequest{}, err
	}
	return user, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id int, user models.UserRequest) (models.UserRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	updatedUser, err := s.UserRepo.UpdateUser(ctx, id, user)
	if err != nil {
		log.Printf("UpdateUser service error: %v\n", err)
		return models.UserRequest{}, err
	}
	return updatedUser, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := s.UserRepo.DeleteUser(ctx, id)
	if err != nil {
		log.Printf("DeleteUser service error: %v\n", err)
		return err
	}
	return nil
}
