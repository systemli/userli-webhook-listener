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
	config := BuildConfig()
	logger.Info("Starting server", zap.String("listenAddr", config.ListenAddr))
	nc := NewNextcloud(config.Nextcloud)
	s := NewServer(nc)
	if err := s.Start(config.ListenAddr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
