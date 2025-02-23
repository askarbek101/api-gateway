package services_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"api-gateway/internal/services"
)

func TestNewTestService(t *testing.T) {
	service := services.NewTestService()
	assert.NotNil(t, service, "TestService should not be nil")
}

func TestTestService_Test(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create service and call Test method
	service := services.NewTestService()
	service.Test(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response structure
	assert.Contains(t, response, "time")
	assert.Contains(t, response, "note")
	assert.Equal(t, "Use Cache-Control header to control caching behavior", response["note"])
}
