package main

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	logLevel := "info"
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}

	atomic := zap.NewAtomicLevel()
	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	atomic.SetLevel(level)
	logger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atomic,
	))
}

func main() {
	listenAddr := ":8080"
	if os.Getenv("LISTEN_ADDR") != "" {
		listenAddr = os.Getenv("LISTEN_ADDR")
	}

	nextcloudUrl := os.Getenv("NEXTCLOUD_OIDC_USER_API_URL")
	if nextcloudUrl == "" {
		logger.Fatal("NEXTCLOUD_OIDC_USER_API_URL environment variable is required")
	}
	nextcloudUsername := os.Getenv("NEXTCLOUD_ADMIN_USERNAME")
	if nextcloudUsername == "" {
		logger.Fatal("NEXTCLOUD_ADMIN_USERNAME environment variable is required")
	}
	nextcloudPassword := os.Getenv("NEXTCLOUD_ADMIN_PASSWORD")
	if nextcloudPassword == "" {
		logger.Fatal("NEXTCLOUD_ADMIN_PASSWORD environment variable is required")
	}
	nextcloudProviderID := os.Getenv("NEXTCLOUD_OIDC_PROVIDER_ID")
	if nextcloudProviderID == "" {
		logger.Fatal("NEXTCLOUD_OIDC_PROVIDER_ID environment variable is required")
	}

	logger.Info("Starting server", zap.String("listenAddr", listenAddr))
	nc, err := NewNextcloud(nextcloudUrl, nextcloudUsername, nextcloudPassword, nextcloudProviderID)
	if err != nil {
		logger.Fatal("Failed to create Nextcloud client", zap.Error(err))
	}

	s := NewServer(nc)
	if err := s.Start(listenAddr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
