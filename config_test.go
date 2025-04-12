package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func (s *ConfigSuite) TestBuildConfig() {
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("LISTEN_ADDR", ":8080")
	os.Setenv("NEXTCLOUD_OIDC_USER_API_URL", "https://example.com/ocs/v2.php/apps/user_oidc/api/v1/user")
	os.Setenv("NEXTCLOUD_ADMIN_USERNAME", "admin")
	os.Setenv("NEXTCLOUD_ADMIN_PASSWORD", "password")
	os.Setenv("NEXTCLOUD_OIDC_PROVIDER_ID", "1")
	os.Setenv("NEXTCLOUD_USER_DOMAIN", "example.com")

	cfg := BuildConfig()

	s.Equal("info", cfg.LogLevel)
	s.Equal(":8080", cfg.ListenAddr)

	s.NotNil(cfg.Nextcloud)
	s.Equal("https://example.com/ocs/v2.php/apps/user_oidc/api/v1/user", cfg.Nextcloud.ApiUrl)
	s.Equal("admin", cfg.Nextcloud.Username)
	s.Equal("password", cfg.Nextcloud.Password)
	s.Equal("1", cfg.Nextcloud.ProviderID)
	s.Equal("example.com", cfg.Nextcloud.Domain)
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
