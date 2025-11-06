package config

import (
	"container/list"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// CORSMiddleware handles CORS configuration
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := []string{
			"http://localhost:3000", // React default port
			"http://localhost:3001", // Alternative React port
			"http://localhost:5173", // Vite default port
			"https://your-production-domain.com",
		}

		origin := c.Request.Header.Get("Origin")
		if origin != "" && slices.Contains(allowedOrigins, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// getWhitelist returns a list of paths that don't require authentication
func getWhitelist() *list.List {
	whitelist := list.New()

	// Auth routes (public - no authentication required)
	whitelist.PushBack("/api/auth/register")
	whitelist.PushBack("/api/auth/login")
	whitelist.PushBack("/api/auth/forgot-password")
	whitelist.PushBack("/api/auth/reset-password")

	// Public product routes
	whitelist.PushBack("/api/products")
	whitelist.PushBack("/api/products/")
	whitelist.PushBack("/api/products/:id")
	whitelist.PushBack("/api/categories")
	whitelist.PushBack("/api/categories/")
	whitelist.PushBack("/api/categories/:id")

	// Public file routes
	whitelist.PushBack("/api/files/:filename")

	// Health check
	whitelist.PushBack("/health")
	whitelist.PushBack("/api/health")

	// Swagger docs
	whitelist.PushBack("/swagger/*any")
	whitelist.PushBack("/docs")

	return whitelist
}

// AuthMiddleware handles JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	whitelist := getWhitelist()

	return func(c *gin.Context) {
		// Skip authentication for whitelisted routes
		for e := whitelist.Front(); e != nil; e = e.Next() {
			if strings.HasPrefix(c.Request.URL.Path, e.Value.(string)) {
				c.Next()
				return
			}
		}

		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if token is in Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		// Verify token
		claims, err := JWT.Verify(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("organization_id", claims.OrganizationID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// AdminMiddleware checks if the user has admin privileges
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Check if user has admin or super_admin role
		userRole := role.(string)
		if userRole != "admin" && userRole != "super_admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware handles panics and returns a 500 error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log the error
				// logger.Error("Recovered from panic", zap.Any("error", r))

				// Return a 500 error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
