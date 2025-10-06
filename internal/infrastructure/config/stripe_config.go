package config

// StripeConfig holds the configuration for Stripe
type StripeConfig struct {
	APIKey         string `mapstructure:"api_key"`
	WebhookSecret  string `mapstructure:"webhook_secret"`
	SuccessURL     string `mapstructure:"success_url"`
	CancelURL      string `mapstructure:"cancel_url"`
	DefaultCurrency string `mapstructure:"default_currency"`
	TestMode       bool   `mapstructure:"test_mode"`
	WebhookPath    string `mapstructure:"webhook_path"`
}

// NewStripeConfig creates a new Stripe configuration from the global config
func NewStripeConfig() (*StripeConfig, error) {
	// Get the global config
	appCfg := Get()

	// Create a new StripeConfig instance
	cfg := &StripeConfig{
		APIKey:         appCfg.Stripe.APIKey,
		WebhookSecret:  appCfg.Stripe.WebhookSecret,
		SuccessURL:     appCfg.Stripe.SuccessURL,
		CancelURL:      appCfg.Stripe.CancelURL,
		DefaultCurrency: appCfg.Stripe.DefaultCurrency,
		TestMode:       appCfg.Stripe.TestMode,
		WebhookPath:    appCfg.Stripe.WebhookPath,
	}

	// Set default values if not provided
	if cfg.SuccessURL == "" {
		cfg.SuccessURL = "http://localhost:3000/success"
	}
	if cfg.CancelURL == "" {
		cfg.CancelURL = "http://localhost:3000/cancel"
	}
	if cfg.DefaultCurrency == "" {
		cfg.DefaultCurrency = "USD"
	}
	if cfg.WebhookPath == "" {
		cfg.WebhookPath = "/api/webhooks/stripe"
	}

	return cfg, nil
}
