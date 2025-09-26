package main

import (
    "github.com/gin-gonic/gin"
    "go-mvcs-boilerplate/config"
    "go-mvcs-boilerplate/mysqldb"
    "go-mvcs-boilerplate/routes"
)

func main() {
    // Load config
    cfg := config.LoadConfig()

    // Connect DB
    mysqldb.ConnectMySQL(cfg)

    // Setup Router
    r := gin.Default()
    routes.RegisterRoutes(r)

    // Start server
    r.Run(":8080")
}
