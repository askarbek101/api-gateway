package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from panics and converts them to APIResponse
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				stack := string(debug.Stack())
				fmt.Printf("Recovery from panic: %v\nStack trace:\n%s\n", err, stack)

				// Only attempt to send error response if headers haven't been written
				if !c.Writer.Written() {
					c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
						"message": "Request could not be processed",
						"code":    http.StatusUnprocessableEntity,
					})
				} else {
					// If headers were already written, just abort
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}
