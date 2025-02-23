package handlers

import (
	"fmt"
	"api-gateway/internal/services"

	"github.com/gin-gonic/gin"
)

// TestHandler handles HTTP requests related to testing
type TestHandler struct {
	testService services.ITestService
}

// NewTestHandler creates a new instance of TestHandler
func NewTestHandler(testService services.ITestService) *TestHandler {
	return &TestHandler{
		testService: testService,
	}
}

// Test handles the test endpoint request
// @Summary Test endpoint
// @Description Simple test endpoint to check if the API is working
// @Tags test
// @Accept json
// @Produce json
// @Success 200 {object} responses.MessageResponse
// @Router /test [get]
func (h *TestHandler) Test(c *gin.Context) {
	fmt.Println("TestHandler.Test")
	h.testService.Test(c)
	fmt.Println("TestHandler.Test done")
}
