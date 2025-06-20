package main

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type SynapseSuite struct {
	suite.Suite

	synapse *Synapse
}

func (s *SynapseSuite) SetupTest() {
	gock.DisableNetworking()
	defer gock.Off()

	config := &SynapseConfig{
		ApiUrl:      "https://example.com/_synapse/admin/v1",
		AccessToken: "token",
		Domain:      "example.com",
	}

	s.synapse = NewSynapse(config)
}

func (s *SynapseSuite) TestDeprovisionUser() {
	s.Run("happy path", func() {
		gock.New("https://example.com").
			Post("/_synapse/admin/v1/deactivate/@user:example.com").
			Reply(200)

		err := s.synapse.DeprovisionUser("user@example.com")
		s.NoError(err)
	})

	s.Run("wrong domain", func() {
		err := s.synapse.DeprovisionUser("user@example.org")
		s.Error(err)
	})

	s.Run("user not found", func() {
		gock.New("https://example.com").
			Post("/_synapse/admin/v1/deactivate/@user:example.com").
			Reply(404)

		err := s.synapse.DeprovisionUser("user@example.com")
		s.NoError(err)
	})

	s.Run("error path", func() {
		gock.New("https://example.com").
			Delete("/_synapse/admin/v1/deactivate/user").
			Reply(500)

		err := s.synapse.DeprovisionUser("user@example.com")
		s.Error(err)
	})
}

func TestSynapseSuite(t *testing.T) {
	suite.Run(t, new(SynapseSuite))
}
