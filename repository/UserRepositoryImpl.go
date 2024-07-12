package repository

import (
	"User-Management/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
)

type UserRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (repo *UserRepositoryImpl) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	stmt, err := repo.DB.PrepareContext(ctx, "INSERT INTO users (name, email, password) VALUES "+
		"(?, ?, ?),(?, ?, ?),(?, ?, ?)")
	if err != nil {
		log.Printf("Create repository prepare error: %v\n", err)
		return models.User{}, err
	}
	defer func() {
		if cerr := stmt.Close(); cerr != nil {
			log.Printf("stmt.Close() error: %v", cerr)
		}
	}()

	result, err := stmt.ExecContext(ctx, user.Name, user.Email, user.Password, user.Name, user.Email, user.Password, user.Name, user.Email, user.Password)
	if err != nil {
		log.Printf("Create repository exec error: %v\n", err)
		return models.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Create repository last insert id error: %v\n", err)
		return models.User{}, err
	}
	user.ID = int(id)
	return user, nil
}

func (repo *UserRepositoryImpl) getCount(ctx context.Context) (int, error) {
	var totalRecords int
	err := repo.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&totalRecords)
	if err != nil {
		return 0, fmt.Errorf("getCount error: %w", err)
	}
	return totalRecords, nil
}

func (repo *UserRepositoryImpl) GetUsers(ctx context.Context, page int, pageSize int) ([]models.UserRequest, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	totalRecords, err := repo.getCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	if page > totalPages {
		return []models.UserRequest{}, totalPages, nil
	}

	query := "SELECT id, name, email FROM users LIMIT ? OFFSET ?"
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("GetUsers prepare error: %w", err)
	}
	defer func() {
		if cerr := stmt.Close(); cerr != nil {
			log.Printf("stmt.Close() error: %v", cerr)
		}
	}()

	rows, err := stmt.QueryContext(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("GetUsers query error: %w", err)
	}
	defer func() {
		if rerr := rows.Close(); rerr != nil {
			log.Printf("rows.Close() error: %v", rerr)
		}
	}()

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GetUsers rows error: %w", err)
	}

	var users []models.UserRequest
	for rows.Next() {
		var user models.UserRequest
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, 0, fmt.Errorf("GetUsers scan error: %w", err)
		}
		users = append(users, user)
	}

	return users, totalPages, nil
}

func (repo *UserRepositoryImpl) GetUserByID(ctx context.Context, id int) (models.UserRequest, error) {
	stmt, err := repo.DB.PrepareContext(ctx, "SELECT id, name, email FROM users WHERE id = ?")
	if err != nil {
		return models.UserRequest{}, fmt.Errorf("GetUserByID prepare error: %w", err)
	}
	defer func() {
		if cerr := stmt.Close(); cerr != nil {
			log.Printf("stmt.Close() error: %v", cerr)
		}
	}()

	var user models.UserRequest
	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserRequest{}, nil
		}
		return models.UserRequest{}, fmt.Errorf("GetUserByID scan error: %w", err)
	}
	return user, nil
}

func (repo *UserRepositoryImpl) UpdateUser(ctx context.Context, id int, user models.UserRequest) (models.UserRequest, error) {
	stmt, err := repo.DB.PrepareContext(ctx, "UPDATE users SET name = ?, email = ? WHERE id = ?")
	if err != nil {
		return models.UserRequest{}, fmt.Errorf("UpdateUser prepare error: %w", err)
	}
	defer func() {
		if cerr := stmt.Close(); cerr != nil {
			log.Printf("stmt.Close() error: %v", cerr)
		}
	}()

	_, err = stmt.ExecContext(ctx, user.Name, user.Email, id)
	if err != nil {
		return models.UserRequest{}, fmt.Errorf("UpdateUser exec error: %w", err)
	}

	return repo.GetUserByID(ctx, id)
}

func (repo *UserRepositoryImpl) DeleteUser(ctx context.Context, id int) error {
	stmt, err := repo.DB.PrepareContext(ctx, "DELETE FROM users WHERE id = ?")
	if err != nil {
		return fmt.Errorf("DeleteUser prepare error: %w", err)
	}
	defer func() {
		if cerr := stmt.Close(); cerr != nil {
			log.Printf("stmt.Close() error: %v", cerr)
		}
	}()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("DeleteUser exec error: %w", err)
	}
	return nil
}
