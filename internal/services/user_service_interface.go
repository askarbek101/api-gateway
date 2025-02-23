package services

import (
	"context"
	"api-gateway/internal/models/requests"
	"api-gateway/internal/models/responses"
)

// UserService defines the interface for user-related business operations
type IUserService interface {
	// User management
	CreateUser(ctx context.Context, req *requests.CreateUserRequest) (*responses.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*responses.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, req *requests.UpdateUserRequest) (*responses.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, page, pageSize int) ([]responses.UserResponse, error)

	// Authentication operations can be added here later
	// Login(ctx context.Context, req *requests.LoginRequest) (*responses.TokenResponse, error)
	// RefreshToken(ctx context.Context, req *requests.RefreshTokenRequest) (*responses.TokenResponse, error)
}
