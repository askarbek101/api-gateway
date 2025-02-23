package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateRequest validates the request body against the provided struct
func ValidateRequest(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&model); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			c.Abort()
			return
		}

		if err := validate.Struct(model); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errors := make([]string, len(validationErrors))
				for i, e := range validationErrors {
					errors[i] = e.Field() + " " + e.Tag()
				}
				c.JSON(http.StatusBadRequest, gin.H{
					"errors": errors,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
