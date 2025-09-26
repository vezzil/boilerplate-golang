package httpserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"boilerplate-golang/internal/config"
)

// New creates a configured gin.Engine with recovery and CORS middlewares.
func New() *gin.Engine {
	r := gin.New()
	r.Use(jsonRecovery())

	cfg := config.Get()
	if len(cfg.CORS.AllowedOrigins) > 0 {
		r.Use(corsWithWhitelist(cfg.CORS.AllowedOrigins))
	} else {
		r.Use(corsPermissive())
	}

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"time":    time.Now().UTC().Format(time.RFC3339),
			"service": cfg.App.Name,
		})
	})

	// 404 JSON
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "route not found",
			"path":    c.Request.URL.Path,
		})
	})

	return r
}

func jsonRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "internal server error",
		})
	})
}

func corsPermissive() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func corsWithWhitelist(allowed []string) gin.HandlerFunc {
	norm := make([]string, 0, len(allowed))
	allowAll := false
	for _, o := range allowed {
		if o == "*" {
			allowAll = true
		}
		norm = append(norm, strings.ToLower(o))
	}
	return func(c *gin.Context) {
		origin := strings.ToLower(c.GetHeader("Origin"))
		var allow string
		if allowAll && origin != "" {
			allow = origin
		} else {
			for _, o := range norm {
				if origin == o {
					allow = origin
					break
				}
			}
		}
		if allow != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allow)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
