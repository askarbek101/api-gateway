package responses

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message" example:"operation successful"`
}

// APIResponse is a generic response model for all API responses
type APIResponse[T any] struct {
	Data    T      `json:"data,omitempty"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
