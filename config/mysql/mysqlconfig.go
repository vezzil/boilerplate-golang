package mysql

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"boilerplate-golang/src/entity"
	"boilerplate-golang/config/toml"
)

var DB *gorm.DB

func init() {
	username := toml.GetConfig().Database.User
	password := toml.GetConfig().Config.Database.Password
	host := toml.GetConfig().Config.Database.Host
	port := toml.GetConfig().Config.Database.Port
	name := toml.GetConfig().Config.Database.Name
	charset := toml.GetConfig().Config.Database.Charset
	parseTime := toml.GetConfig().Config.Database.ParseTime
	loc := toml.GetConfig().Config.Database.Loc
	timeout := toml.GetConfig().Config.Database.Timeout

	// Build DSN from [database] section
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s&timeout=%s",
		username,
		password,
		host,
		port,
		name,
		charset,
		parseTime,
		loc,
		timeout,
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