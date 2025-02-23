package services

import "github.com/gin-gonic/gin"

// Service defines the interface for HTTP Services
type ITestService interface {
	Test(c *gin.Context)
}
