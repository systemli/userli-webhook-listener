package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	router        *chi.Mux
	nextcloud     *Nextcloud
	synapse       *Synapse
	webhookSecret string
}

func NewServer(webhookSecret string, nextcloud *Nextcloud, synapse *Synapse) *Server {
	return &Server{
		router:        chi.NewRouter(),
		webhookSecret: webhookSecret,
		nextcloud:     nextcloud,
		synapse:       synapse,
	}
}

func (s *Server) Start(addr string) error {
	s.RegisterRoutes()

	return http.ListenAndServe(addr, s.router)
}

func (s *Server) RegisterRoutes() {
	s.router.Use(s.Authmiddleware)
	s.router.Post("/userli", s.handleUserliEvent)
}

func (s *Server) handleUserliEvent(w http.ResponseWriter, r *http.Request) {
	logger.Info("Userli event received")

	var event UserEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case EventTypeUserCreated:
		s.handleUserCreated(event)
	case EventTypeUserDeleted:
		s.handleUserDeleted(event)
	default:
		http.Error(w, "Unknown event type", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUserCreated(event UserEvent) {
	logger.With(zap.String("email", event.Data.Email)).Info("User created event received")

	err := s.nextcloud.ProvisionUser(event.Data.Email)
	if err != nil {
		logger.Error("Failed to provision user in Nextcloud")
	}
}

func (s *Server) handleUserDeleted(event UserEvent) {
	logger.With(zap.String("email", event.Data.Email)).Info("User deleted event received")

	err := s.nextcloud.DeprovisionUser(event.Data.Email)
	if err != nil {
		logger.Error("Failed to deprovision user in Nextcloud")
	}

	err = s.synapse.DeprovisionUser(event.Data.Email)
	if err != nil {
		logger.Error("Failed to deprovision user in Synapse")
	}
}

func (s *Server) Authmiddleware(next http.Handler) http.Handler {
	secret := s.webhookSecret

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get("X-Signature")
		if signature == "" {
			http.Error(w, "Missing signature header", http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)

		if !hmac.Equal([]byte(signature), []byte(hex.EncodeToString(mac.Sum(nil)))) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
