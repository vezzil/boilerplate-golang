package main

import (
	"fmt"

	"boilerplate-golang/internal/application/router"
	"boilerplate-golang/internal/infrastructure/ai"
	"boilerplate-golang/internal/infrastructure/config"
	"boilerplate-golang/internal/infrastructure/cronmanager"
	"boilerplate-golang/internal/infrastructure/dbmanager"
	"boilerplate-golang/internal/infrastructure/redismanager"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize infrastructure managers
	dbmanager.Init()
	redismanager.Init()
	cronmanager.Init()
	ai.Init()

	// Create Gin router with default middleware
	r := gin.Default()

	// Configure trusted proxies for production security
	// In development: trust localhost
	// In production: specify your actual proxy IPs (load balancer, reverse proxy, etc.)
	if cfg.App.Env == "production" {
		// Set specific trusted proxy IPs for production
		// Example: r.SetTrustedProxies([]string{"10.0.0.1", "10.0.0.2"})
		r.SetTrustedProxies(nil) // Don't trust any proxies by default
	} else {
		// For local development
		r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	}

	// Apply CORS middleware from router_config
	r.Use(config.CORSMiddleware())

	// Get database connection
	db := dbmanager.GetDB()

	// Initialize JWT
	config.InitJWT()

	// Register application routes
	router.Register(r, db)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	_ = r.Run(addr)
}
