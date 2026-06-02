package middleware

import (
	"cinema-backend/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	authService := services.NewAuthService(nil, jwtSecret)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CustomerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || (role != "customer" && role != "admin") {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "customer access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AdminOrCustomerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || (role != "customer" && role != "admin") {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
