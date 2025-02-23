package responses

import (
	"api-gateway/internal/models"
	"time"
)

// UserResponse represents the response body for user data
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromUser converts a User model to a UserResponse
func FromUser(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// FromUsers converts a slice of User models to UserResponses
func FromUsers(users []models.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = *FromUser(&user)
	}
	return responses
}
