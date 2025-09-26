package mysqldb

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"boilerplate-golang/config"
	"boilerplate-golang/src/entity"
)

var DB *gorm.DB

// ConnectMySQL initializes the database connection
func ConnectMySQL(cfg config.Config) {
	// Build DSN from [database] section
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Charset,
		cfg.Database.ParseTime,
		cfg.Database.Loc,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Auto-migrate your entities
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatal("❌ Failed to run migrations:", err)
	}

	DB = db
	log.Println("✅ MySQL connection established & migrations completed")
}
