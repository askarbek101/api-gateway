package services

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// TestService implements the Service interface
type TestService struct{}

// NewTestService creates a new instance of TestService
func NewTestService() *TestService {
	return &TestService{}
}

// Test handles the test endpoint request
func (h *TestService) Test(c *gin.Context) {
	fmt.Println("Test")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Test2")
	c.JSON(200, gin.H{
		"time": time.Now(),
		"note": "Use Cache-Control header to control caching behavior",
	})
}
