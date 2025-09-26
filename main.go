package main

import (
	"fmt"
	"boilerplate-golang/internal/config"
	"boilerplate-golang/internal/manager/dbmanager"
	"boilerplate-golang/internal/manager/cachemanager"
	"boilerplate-golang/internal/manager/cronmanager"
	"boilerplate-golang/internal/server/httpserver"
	"boilerplate-golang/internal/router"
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
