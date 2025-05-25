// internal/utils/response.go
package utils

import (
	"github.com/gin-gonic/gin"
)

// RespondWithSuccess sends a successful JSON response.
func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"success": true, "data": data})
}

// RespondWithError sends an error JSON response.
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"success": false, "error": message})
}
