package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// MockResponse represents a mock response configuration
type MockResponse struct {
	Data       interface{}
	StatusCode int
	Error      error
}

// HTTPSender handles HTTP requests
type HTTPSender struct {
	client   *http.Client
	baseURL  string
	mockMode bool
	mockData map[string]MockResponse
}

// NewHTTPSender creates a new instance of HTTPSender
func NewHTTPSender(baseURL string, timeout time.Duration) *HTTPSender {
	return &HTTPSender{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL:  baseURL,
		mockMode: false,
		mockData: make(map[string]MockResponse),
	}
}

// EnableMockMode enables mock mode and initializes mock data
func (s *HTTPSender) EnableMockMode() {
	s.mockMode = true
}

// DisableMockMode disables mock mode
func (s *HTTPSender) DisableMockMode() {
	s.mockMode = false
}

// SetMockResponse sets a mock response for a specific path and method
func (s *HTTPSender) SetMockResponse(method, path string, response MockResponse) {
	key := fmt.Sprintf("%s:%s", method, path)
	s.mockData[key] = response
}

// SendRequest sends an HTTP request and returns the response
func (s *HTTPSender) SendRequest(ctx context.Context, method, path string, body interface{}, response interface{}) error {
	if s.mockMode {
		key := fmt.Sprintf("%s:%s", method, path)
		if mockResp, exists := s.mockData[key]; exists {
			if mockResp.Error != nil {
				return mockResp.Error
			}
			if response != nil && mockResp.Data != nil {
				mockJSON, err := json.Marshal(mockResp.Data)
				if err != nil {
					return fmt.Errorf("failed to marshal mock data: %w", err)
				}
				if err := json.Unmarshal(mockJSON, response); err != nil {
					return fmt.Errorf("failed to unmarshal mock data: %w", err)
				}
			}
			return nil
		}
		// Return default mock response if no specific mock is set
		if response != nil {
			defaultMock := map[string]interface{}{"message": "default mock response"}
			mockJSON, _ := json.Marshal(defaultMock)
			return json.Unmarshal(mockJSON, response)
		}
		return nil
	}

	url := fmt.Sprintf("%s%s", s.baseURL, path)

	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Get sends a GET request
func (s *HTTPSender) Get(ctx context.Context, path string, response interface{}) error {
	return s.SendRequest(ctx, http.MethodGet, path, nil, response)
}

// Post sends a POST request
func (s *HTTPSender) Post(ctx context.Context, path string, body interface{}, response interface{}) error {
	return s.SendRequest(ctx, http.MethodPost, path, body, response)
}

// Put sends a PUT request
func (s *HTTPSender) Put(ctx context.Context, path string, body interface{}, response interface{}) error {
	return s.SendRequest(ctx, http.MethodPut, path, body, response)
}

// Delete sends a DELETE request
func (s *HTTPSender) Delete(ctx context.Context, path string) error {
	return s.SendRequest(ctx, http.MethodDelete, path, nil, nil)
}
