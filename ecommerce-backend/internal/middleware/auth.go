// internal/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/config"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/utils" // For JWT validation
)

// AuthMiddleware authenticates requests using JWT.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondWithError(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort() // Abort the request if unauthorized
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid Authorization header format")
			c.Abort()
			return
		}
		tokenString := parts[1]

		claims, err := utils.ValidateJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid or expired token: "+err.Error())
			c.Abort()
			return
		}

		// Store user information in Gin's context for later use by handlers
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role) // This will be used for authorization

		c.Next() // Proceed to the next handler/middleware in the chain
	}
}

// AuthorizeRole checks if the authenticated user has the required role.
func AuthorizeRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the role from the context, which was set by AuthMiddleware
		role, exists := c.Get("userRole")
		if !exists {
			utils.RespondWithError(c, http.StatusForbidden, "Role information not found (missing auth middleware?)")
			c.Abort()
			return
		}

		if role != requiredRole {
			utils.RespondWithError(c, http.StatusForbidden, "Forbidden: Insufficient permissions")
			c.Abort()
			return
		}

		c.Next() // User has the required role, proceed
	}
}
