package toml

import (
	"log"

	"github.com/spf13/viper"
)

// TomlConfig struct maps to config.toml
// This struct mirrors the structure of your TOML file.
type TomlConfig struct {
	App struct {
		Name string // application name
		Port int    // application port
		Env  string // environment (e.g. development, production)
	}
	Database struct {
		Host      string // DB host, e.g. "127.0.0.1"
		Port      int    // DB port, e.g. 3306
		User      string // DB username
		Password  string // DB password
		Name      string // DB name
		Charset   string // DB charset, e.g. "utf8mb4"
		ParseTime bool   // enable parsing of time values
		Loc       string // DB timezone, e.g. "Local"
		Timeout   string // DB connection timeout
	}
	CronJob struct {
		CleanupInterval string // cleanup interval
		EmailReport     string // email report schedule
	}
	Redis struct {
		Host     string // Redis host, e.g. "127.0.0.1"
		Port     int    // Redis port, e.g. 6379
		Password string // Redis password
		DB       int    // Redis database index
	}
	Log struct {
		Level string // Log level, e.g. "debug", "info", "warn", "error"
		File  string // Log file path
	}
}

var Config TomlConfig

// LoadTomlConfig reads config.toml and unmarshals into TomlConfig struct
func LoadTomlConfig() TomlConfig {
	// Set the config file name (without extension)
	viper.SetConfigName("config")

	// Tell viper the config file type is TOML
	viper.SetConfigType("toml")

	// Look for the config file in the current directory
	viper.AddConfigPath(".")

	// Allow environment variables to override config values
	viper.AutomaticEnv()

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		// Fatal if config file can’t be read
		log.Fatalf("❌ Error reading config file: %v", err)
	}

	// Create a variable of type TomlConfig
	var cfg TomlConfig

	// Unmarshal config values into the struct
	err = viper.Unmarshal(&cfg)
	if err != nil {
		// Fatal if config values can’t map into struct
		log.Fatalf("❌ Unable to decode into struct: %v", err)
	}

	// Success message
	log.Println("✅ Configuration loaded successfully")

	// Return the config object
	return cfg
}

// GetConfig is a helper function to load and return config
func GetConfig() TomlConfig {
	return Config
}