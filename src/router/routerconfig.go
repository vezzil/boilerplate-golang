package router

import (
    "net/http"
    "time"
    "strings"

    "github.com/gin-gonic/gin"
)

// ApplyDefaults attaches common middlewares and default routes
// such as CORS, Recovery, Health check, and a JSON 404 handler.
func ApplyDefaults(r *gin.Engine) {
    // Middlewares
    r.Use(CustomRecovery())
    r.Use(CORSMiddleware())

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "ok",
            "time":    time.Now().UTC().Format(time.RFC3339),
            "service": "boilerplate-golang",
        })
    })

    // 404 JSON response
    r.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, gin.H{
            "status":  "error",
            "message": "route not found",
            "path":    c.Request.URL.Path,
        })
    })
}

// ApplyDefaultsWithWhitelist is like ApplyDefaults, but uses a CORS whitelist.
// If the request Origin matches one of the allowed origins, it will be echoed back.
// If the allowed list contains "*", all origins are allowed (like the dev middleware).
func ApplyDefaultsWithWhitelist(r *gin.Engine, allowedOrigins []string) {
    r.Use(CustomRecovery())
    r.Use(CORSMiddlewareWithWhitelist(allowedOrigins))

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "ok",
            "time":    time.Now().UTC().Format(time.RFC3339),
            "service": "boilerplate-golang",
        })
    })

    r.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, gin.H{
            "status":  "error",
            "message": "route not found",
            "path":    c.Request.URL.Path,
        })
    })
}

// CORSMiddleware provides a permissive CORS policy suitable for development.
// Adjust AllowedOrigins/Methods/Headers as needed for production.
func CORSMiddleware() gin.HandlerFunc {
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

// CORSMiddlewareWithWhitelist allows only configured origins. If "*" is present, all are allowed.
func CORSMiddlewareWithWhitelist(allowed []string) gin.HandlerFunc {
    // Normalize allowed list to lower-case for comparison
    normalized := make([]string, 0, len(allowed))
    allowAll := false
    for _, o := range allowed {
        if o == "*" {
            allowAll = true
        }
        normalized = append(normalized, strings.ToLower(o))
    }
    return func(c *gin.Context) {
        origin := strings.ToLower(c.GetHeader("Origin"))
        allowedOrigin := ""
        if allowAll && origin != "" {
            allowedOrigin = origin
        } else {
            for _, o := range normalized {
                if origin == o {
                    allowedOrigin = origin
                    break
                }
            }
        }

        if allowedOrigin != "" {
            c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
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

// CustomRecovery recovers from panics and returns a JSON error response.
func CustomRecovery() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, err any) {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "internal server error",
        })
    })
}

