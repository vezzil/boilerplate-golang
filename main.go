package main

import (
	"fmt"

	"boilerplate-golang/internal/application/router"
	"boilerplate-golang/internal/infrastructure/redismanager"
	"boilerplate-golang/internal/infrastructure/config"
	"boilerplate-golang/internal/infrastructure/cronmanager"
	"boilerplate-golang/internal/infrastructure/dbmanager"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize infrastructure managers
	dbmanager.Init()
	cachemanager.Init()
	cronmanager.Init()

	// Create Gin router with default middleware
	r := gin.Default()

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
