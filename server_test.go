package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	server *Server
}

func (s *ServerSuite) SetupTest() {
	gock.DisableNetworking()
	defer gock.Off()

	config := &NextcloudConfig{
		ApiUrl:     "https://example.com/ocs/v2.php/apps/user_oidc/api/v1/user",
		Username:   "admin",
		Password:   "password",
		ProviderID: "provider-id",
		Domain:     "example.com",
	}
	nc := NewNextcloud(config)

	s.server = NewServer(nc)
}

func (s *ServerSuite) TestHandleWebhook() {
	s.Run("invalid request body", func() {
		req := httptest.NewRequest("POST", "/userli", bytes.NewBuffer([]byte("invalid")))
		w := httptest.NewRecorder()
		s.server.handleUserliEvent(w, req)
		s.Equal(400, w.Code)
	})

	s.Run("unknown event type", func() {
		event := UserEvent{
			Type: "unknown",
			Data: struct {
				Email string `json:"email"`
			}{
				Email: "user@example.com",
			},
		}
		jsonData, err := json.Marshal(event)
		s.NoError(err)

		req := httptest.NewRequest("POST", "/userli", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()
		s.server.handleUserliEvent(w, req)
		s.Equal(400, w.Code)
	})
}

func (s *ServerSuite) TestHandleUserCreated() {
	s.Run("happy path", func() {
		gock.New("https://example.com").
			Post("/ocs/v2.php/apps/user_oidc/api/v1/user").
			Reply(200)

		event := UserEvent{
			Type: EventTypeUserCreated,
			Data: struct {
				Email string `json:"email"`
			}{
				Email: "user@example.com",
			},
		}
		jsonData, err := json.Marshal(event)
		s.NoError(err)

		req := httptest.NewRequest("POST", "/userli", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()
		s.server.handleUserliEvent(w, req)
		s.Equal(200, w.Code)
	})
}

func (s *ServerSuite) TestHandleUserDeleted() {
	s.Run("happy path", func() {
		gock.New("https://example.com").
			Delete("/ocs/v2.php/apps/user_oidc/api/v1/user/user").
			Reply(200)

		event := UserEvent{
			Type: EventTypeUserDeleted,
			Data: struct {
				Email string `json:"email"`
			}{
				Email: "user@example.com",
			},
		}
		jsonData, err := json.Marshal(event)
		s.NoError(err)

		req := httptest.NewRequest("POST", "/userli", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()
		s.server.handleUserliEvent(w, req)
		s.Equal(200, w.Code)
	})
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
