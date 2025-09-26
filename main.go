package main

import (
	"boilerplate-golang/config"
	"boilerplate-golang/mysql"
	"boilerplate-golang/src/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect DB
	mysql.ConnectMySQL(cfg)

	// Setup Router
	r := gin.Default()
	router.InitRouter(r)

	// Start server
	r.Run(":8080")
}
