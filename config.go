package main

import "os"

type Config struct {
	LogLevel      string
	ListenAddr    string
	WebhookSecret string
	Nextcloud     *NextcloudConfig
	Synapse       *SynapseConfig
}

func BuildConfig() *Config {
	cfg := &Config{
		LogLevel:   "info",
		ListenAddr: ":8080",
		Nextcloud:  &NextcloudConfig{},
		Synapse:    &SynapseConfig{},
	}

	if os.Getenv("LOG_LEVEL") != "" {
		cfg.LogLevel = os.Getenv("LOG_LEVEL")
	}
	if os.Getenv("LISTEN_ADDR") != "" {
		cfg.ListenAddr = os.Getenv("LISTEN_ADDR")
	}
	cfg.WebhookSecret = getEnvOrFatal("WEBHOOK_SECRET")
	cfg.Nextcloud = &NextcloudConfig{
		ApiUrl:     getEnvOrFatal("NEXTCLOUD_OIDC_USER_API_URL"),
		Username:   getEnvOrFatal("NEXTCLOUD_ADMIN_USERNAME"),
		Password:   getEnvOrFatal("NEXTCLOUD_ADMIN_PASSWORD"),
		ProviderID: getEnvOrFatal("NEXTCLOUD_OIDC_PROVIDER_ID"),
		Domain:     getEnvOrFatal("NEXTCLOUD_USER_DOMAIN"),
	}
	cfg.Synapse = &SynapseConfig{
		ApiUrl:      getEnvOrFatal("SYNAPSE_USER_ADMIN_API_URL"),
		AccessToken: getEnvOrFatal("SYNAPSE_ADMIN_ACCESS_TOKEN"),
		Domain:      getEnvOrFatal("SYNAPSE_USER_DOMAIN"),
	}

	return cfg
}

func getEnvOrFatal(key string) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Fatal(key + " environment variable is required")
	}
	return val
}
