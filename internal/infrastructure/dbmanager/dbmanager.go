package dbmanager

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"boilerplate-golang/internal/infrastructure/config"
	"boilerplate-golang/internal/application/entity"
)

var db *gorm.DB

// Init initializes the global GORM DB connection and runs migrations.
// Database connection is optional and will log a warning if it fails.
func Init() {
	cfg := config.Get()

	// Skip database initialization if host is not set
	if cfg.Database.Host == "" {
		log.Println("Database host not configured, skipping database connection")
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s&timeout=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Charset,
		cfg.Database.ParseTime,
		cfg.Database.Loc,
		cfg.Database.Timeout,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Application will continue to run without database connection")
		return
	}

	// Auto-migrate all models if database is connected
	err = db.AutoMigrate(
		&entity.User{},
		// Add other models here
	)
	if err != nil {
		log.Printf("Warning: Failed to auto-migrate database: %v", err)
	}
	log.Println("dbmanager: connected and migrated")
}

// DB returns the initialized *gorm.DB.
func DB() *gorm.DB { return db }
