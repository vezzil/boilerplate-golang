package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config struct maps to config.toml
type Config struct {
	App struct {
		Name string
		Port int
		Env  string
	}
	Database struct {
		Host      string
		Port      int
		User      string
		Password  string
		Name      string
		Charset   string
		ParseTime bool
		Loc       string
	}
}

// LoadConfig reads config.toml and unmarshals into Config struct
func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("❌ Error reading config file: %v", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("❌ Unable to decode into struct: %v", err)
	}

	log.Println("✅ Configuration loaded successfully")
	return cfg
}
