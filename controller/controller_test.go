package controller

import (
	"User-Management/models"
	"User-Management/service"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct{}

var (
	mockUserCreate = models.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com", Password: "123456"}
	mockUser1      = models.UserRequest{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
	mockUser2      = models.UserRequest{ID: 2, Name: "Jane Doe", Email: "jane.doe@example.com"}
	mockUsers      = []models.UserRequest{mockUser1, mockUser2}
)

func (s *MockUserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	return mockUserCreate, nil
}

func (s *MockUserService) GetUserByID(ctx context.Context, id int) (models.UserRequest, error) {
	switch id {
	case 1:
		return mockUser1, nil
	case 2:
		return mockUser2, nil
	default:
		return models.UserRequest{}, errors.New("user not found")
	}
}

func (s *MockUserService) GetUsers(ctx context.Context, page, pageSize int) ([]models.UserRequest, int, error) {
	return mockUsers, len(mockUsers), nil
}

func (s *MockUserService) UpdateUser(ctx context.Context, id int, user models.UserRequest) (models.UserRequest, error) {
	if id == mockUser1.ID {
		mockUser1.Name = user.Name
		mockUser1.Email = user.Email
		return mockUser1, nil
	}
	return models.UserRequest{}, errors.New("user not found")
}

func (s *MockUserService) DeleteUser(ctx context.Context, id int) error {
	if id == 1 {
		mockUser1 = models.UserRequest{}
		return nil
	}
	return errors.New("user not found")
}

func setupRouter(userService service.UserService) *gin.Engine {
	r := gin.New()
	userController := NewUserController(userService)
	r.POST("/users", userController.CreateUser)
	r.GET("/users/:id", userController.GetUser)
	r.GET("/users", userController.GetUsers)
	r.PUT("/users/:id", userController.UpdateUser)
	r.DELETE("/users/:id", userController.DeleteUser)
	return r
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userService := &MockUserService{}
	r := setupRouter(userService)

	userJSON := `{"Name":"John Doe","Email":"john.doe@example.com","Password":"123456"}`
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users", strings.NewReader(userJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.UserResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, mockUser1.ID, response.ID)
	assert.Equal(t, mockUser1.Name, response.Name)
	assert.Equal(t, mockUser1.Email, response.Email)
}

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userService := &MockUserService{}
	r := setupRouter(userService)

	testCases := []struct {
		Name           string
		userID         int
		expectedStatus int
		expectedBody   string
	}{
		{"GetUser Success", 1, http.StatusOK, `{"id":1,"name":"John Doe","email":"john.doe@example.com"}`},
		{"GetUser Not Found", 3, http.StatusNotFound, `{"error":"user not found"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/users/"+strconv.Itoa(tc.userID), nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userService := &MockUserService{}
	r := setupRouter(userService)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/users?page=1&pageSize=10", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.GetUsersResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, len(mockUsers), response.TotalRecords)
	assert.Equal(t, mockUsers, response.Users)
}

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userService := &MockUserService{}
	r := setupRouter(userService)

	userJSON := `{"name":"Updated Name","email":"updated.email@example.com"}`
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/users/1", strings.NewReader(userJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "updated successfully", responseBody["message"])

	user := responseBody["user"].(map[string]interface{})
	assert.Equal(t, "Updated Name", user["name"])
	assert.Equal(t, "updated.email@example.com", user["email"])
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userService := &MockUserService{}
	r := setupRouter(userService)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/users/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "deleted successfully", responseBody["message"])
}
