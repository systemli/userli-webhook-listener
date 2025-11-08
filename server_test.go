package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
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

	ncConfig := &NextcloudConfig{
		ApiUrl:     "https://example.com/ocs/v2.php/apps/user_oidc/api/v1/user",
		Username:   "admin",
		Password:   "password",
		ProviderID: "provider-id",
		Domain:     "example.com",
	}
	nc := NewNextcloud(ncConfig)

	syConfig := &SynapseConfig{
		ApiUrl:      "https://example.com/_synapse/admin/v1",
		AccessToken: "token",
		Domain:      "example.com",
	}
	sy := NewSynapse(syConfig)

	s.server = NewServer("secret", nc, sy)
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

		gock.New("https://example.com").
			Post("/_synapse/admin/v1/deactivate/user").
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

func (s *ServerSuite) TestAuthMiddleware() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}

	s.Run("valid token", func() {
		payload := []byte(`{"dummy":"data"}`)
		mac := hmac.New(sha256.New, []byte(s.server.webhookSecret))
		mac.Write(payload)
		validSignature := hex.EncodeToString(mac.Sum(nil))

		req := httptest.NewRequest("POST", "/userli", bytes.NewBuffer(payload))
		req.Header.Set("X-Webhook-Signature", validSignature)

		rr := httptest.NewRecorder()
		wrappedHandler := s.server.Authmiddleware(http.HandlerFunc(handler))
		wrappedHandler.ServeHTTP(rr, req)

		s.Equal(200, rr.Code)
	})

	s.Run("invalid token", func() {
		payload := []byte(`{"dummy":"data"}`)
		req := httptest.NewRequest("POST", "/userli", bytes.NewBuffer(payload))
		req.Header.Set("X-Webhook-Signature", "invalid-signature")

		rr := httptest.NewRecorder()
		wrappedHandler := s.server.Authmiddleware(http.HandlerFunc(handler))
		wrappedHandler.ServeHTTP(rr, req)

		s.Equal(401, rr.Code)
	})

	s.Run("missing token", func() {
		payload := []byte(`{"dummy":"data"}`)
		req := httptest.NewRequest("POST", "/userli", bytes.NewBuffer(payload))

		rr := httptest.NewRecorder()
		wrappedHandler := s.server.Authmiddleware(http.HandlerFunc(handler))
		wrappedHandler.ServeHTTP(rr, req)

		s.Equal(401, rr.Code)
	})
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
