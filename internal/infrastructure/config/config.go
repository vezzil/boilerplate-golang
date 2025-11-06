package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AppConfig defines the full application configuration loaded from config files and env.
type AppConfig struct {
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
		Timeout   string
	}
	JWT struct {
		Secret   string        `mapstructure:"secret"`
		Issuer   string        `mapstructure:"issuer"`
		ExpireIn time.Duration `mapstructure:"expire_in"`
	} `mapstructure:"jwt"`
	Stripe struct {
		APIKey          string `mapstructure:"api_key"`
		WebhookSecret   string `mapstructure:"webhook_secret"`
		SuccessURL      string `mapstructure:"success_url"`
		CancelURL       string `mapstructure:"cancel_url"`
		DefaultCurrency string `mapstructure:"default_currency"`
		TestMode        bool   `mapstructure:"test_mode"`
		WebhookPath     string `mapstructure:"webhook_path"`
	} `mapstructure:"stripe"`
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
	CronJob struct {
		CleanupInterval string
		EmailReport     string
	}
	Log struct {
		Level string
		File  string
	}
	AWS struct {
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		Region          string `mapstructure:"region"`
		S3Bucket        string `mapstructure:"s3_bucket"`
	} `mapstructure:"aws"`
	AI struct {
		OpenAI struct {
			APIKey      string  `mapstructure:"api_key"`
			BaseURL     string  `mapstructure:"base_url"`
			Model       string  `mapstructure:"model"`
			MaxTokens   int     `mapstructure:"max_tokens"`
			Temperature float64 `mapstructure:"temperature"`
		} `mapstructure:"openai"`
		RateLimit struct {
			RequestsPerMinute int `mapstructure:"requests_per_minute"`
			RequestsPerHour   int `mapstructure:"requests_per_hour"`
		} `mapstructure:"rate_limit"`
		Timeout    time.Duration `mapstructure:"timeout"`
		RetryCount int           `mapstructure:"retry_count"`
		Enabled    bool          `mapstructure:"enabled"`
	} `mapstructure:"ai"`
}

var cfg AppConfig

// Load loads configuration from TOML file and environment variables.
// This should be called exactly once at app startup.
func Load() AppConfig {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode config into struct: %v", err)
	}

	// reasonable defaults
	if cfg.App.Port == 0 {
		cfg.App.Port = 8080
	}
	if cfg.Database.Timeout == "" {
		cfg.Database.Timeout = (10 * time.Second).String()
	}

	log.Printf("config loaded: env=%s port=%d", cfg.App.Env, cfg.App.Port)
	return cfg
}

// Get returns the loaded configuration. Call Load() first.
func Get() AppConfig { return cfg }
