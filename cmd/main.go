package main

import (
	"fmt"
	"boilerplate-golang/internal/infrastructure/config"
	"boilerplate-golang/internal/infrastructure/dbmanager"
	"boilerplate-golang/internal/infrastructure/cachemanager"
	"boilerplate-golang/internal/infrastructure/cronmanager"
	"boilerplate-golang/internal/server/httpserver"
	"boilerplate-golang/internal/application/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize infrastructure managers
	dbmanager.Init()
	cachemanager.Init()
	cronmanager.Init()

	// Create HTTP server (gin.Engine)
	r := httpserver.New()

	// Register application routes
	router.Register(r)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	_ = r.Run(addr)
}
