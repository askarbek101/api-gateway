package services

import (
	"context"
	"fmt"
	"api-gateway/internal/models/requests"
	"api-gateway/internal/models/responses"
	"api-gateway/internal/utils/http"
	"time"
)

// UserService implements the UserService interface
type UserService struct {
	httpSender *http.HTTPSender
}

// NewUserService creates a new instance of UserService
func NewUserService(config map[string]string) *UserService {
	baseURL := config["base_url"]
	sender := http.NewHTTPSender(baseURL, 10*time.Second)

	// Enable mock mode
	sender.EnableMockMode()

	// Set up mock responses
	sender.SetMockResponse("GET", "/users", http.MockResponse{
		Data: []responses.UserResponse{
			{ID: 1, Username: "john_doe", Email: "john@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, Username: "jane_smith", Email: "jane@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	})

	sender.SetMockResponse("POST", "/users", http.MockResponse{
		Data: responses.UserResponse{
			ID:        3,
			Username:  "new_user",
			Email:     "new@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})

	// Mock for GetUserByID - will be used for paths like /users/1
	sender.SetMockResponse("GET", "/users/1", http.MockResponse{
		Data: responses.UserResponse{
			ID:        1,
			Username:  "john_doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})

	// Mock for UpdateUser
	sender.SetMockResponse("PUT", "/users/1", http.MockResponse{
		Data: responses.UserResponse{
			ID:        1,
			Username:  "updated_john_doe",
			Email:     "john.updated@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})

	// Mock for DeleteUser
	sender.SetMockResponse("DELETE", "/users/1", http.MockResponse{
		Data: nil,
	})

	return &UserService{
		httpSender: sender,
	}
}

// CreateUser handles user creation
func (s *UserService) CreateUser(ctx context.Context, req *requests.CreateUserRequest) (*responses.UserResponse, error) {
	var response responses.UserResponse
	err := s.httpSender.Post(ctx, "/users", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &response, nil
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*responses.UserResponse, error) {
	var response responses.UserResponse
	err := s.httpSender.Get(ctx, fmt.Sprintf("/users/%d", id), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &response, nil
}

// UpdateUser handles user updates
func (s *UserService) UpdateUser(ctx context.Context, id uint, req *requests.UpdateUserRequest) (*responses.UserResponse, error) {
	var response responses.UserResponse
	err := s.httpSender.Put(ctx, fmt.Sprintf("/users/%d", id), req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &response, nil
}

// DeleteUser handles user deletion
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	err := s.httpSender.Delete(ctx, fmt.Sprintf("/users/%d", id))
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) ([]responses.UserResponse, error) {
	var response []responses.UserResponse
	err := s.httpSender.Get(ctx, fmt.Sprintf("/users?page=%d&page_size=%d", page, pageSize), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return response, nil
}
