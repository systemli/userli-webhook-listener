package main

import "os"

type Config struct {
	LogLevel   string
	ListenAddr string
	Nextcloud  *NextcloudConfig
}

func BuildConfig() *Config {
	cfg := &Config{
		LogLevel:   "info",
		ListenAddr: ":8080",
		Nextcloud:  &NextcloudConfig{},
	}

	if os.Getenv("LOG_LEVEL") != "" {
		cfg.LogLevel = os.Getenv("LOG_LEVEL")
	}
	if os.Getenv("LISTEN_ADDR") != "" {
		cfg.ListenAddr = os.Getenv("LISTEN_ADDR")
	}
	cfg.Nextcloud = &NextcloudConfig{
		ApiUrl:     getEnvOrFatal("NEXTCLOUD_OIDC_USER_API_URL"),
		Username:   getEnvOrFatal("NEXTCLOUD_ADMIN_USERNAME"),
		Password:   getEnvOrFatal("NEXTCLOUD_ADMIN_PASSWORD"),
		ProviderID: getEnvOrFatal("NEXTCLOUD_OIDC_PROVIDER_ID"),
		Domain:     getEnvOrFatal("NEXTCLOUD_USER_DOMAIN"),
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
