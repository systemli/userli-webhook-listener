package main

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type NextcloudSuite struct {
	suite.Suite

	nextcloud *Nextcloud
}

func (s *NextcloudSuite) SetupTest() {
	gock.DisableNetworking()
	defer gock.Off()

	config := &NextcloudConfig{
		ApiUrl:     "https://example.com/ocs/v2.php/apps/user_oidc/api/v1/user",
		Username:   "admin",
		Password:   "password",
		ProviderID: "provider-id",
		Domain:     "example.com",
	}

	s.nextcloud = NewNextcloud(config)
}

func (s *NextcloudSuite) TestProvisionUser() {
	s.Run("happy path", func() {
		gock.New("https://example.com").
			Post("/ocs/v2.php/apps/user_oidc/api/v1/user").
			Reply(200)

		err := s.nextcloud.ProvisionUser("user@example.com")
		s.NoError(err)
	})

	s.Run("wrong domain", func() {
		err := s.nextcloud.ProvisionUser("user@example.org")
		s.Error(err)
	})

	s.Run("error path", func() {
		gock.New("https://example.com").
			Post("/ocs/v2.php/apps/user_oidc/api/v1/user").
			Reply(500)

		err := s.nextcloud.ProvisionUser("user@example.com")
		s.Error(err)
	})
}

func (s *NextcloudSuite) TestDeprovisionUser() {
	s.Run("happy path", func() {
		gock.New("https://example.com").
			Delete("/ocs/v2.php/apps/user_oidc/api/v1/user/user").
			Reply(200)

		err := s.nextcloud.DeprovisionUser("user@example.com")
		s.NoError(err)
	})

	s.Run("wrong domain", func() {
		err := s.nextcloud.DeprovisionUser("user@example.org")
		s.Error(err)
	})

	s.Run("user not found", func() {
		gock.New("https://example.com").
			Delete("/ocs/v2.php/apps/user_oidc/api/v1/user/user").
			Reply(404)

		err := s.nextcloud.DeprovisionUser("user@example.com")
		s.NoError(err)
	})

	s.Run("error path", func() {
		gock.New("https://example.com").
			Delete("/ocs/v2.php/apps/user_oidc/api/v1/user/user").
			Reply(500)

		err := s.nextcloud.DeprovisionUser("user@example.com")
		s.Error(err)
	})
}

func TestNextcloudSuite(t *testing.T) {
	suite.Run(t, new(NextcloudSuite))
}
