package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router    *chi.Mux
	nextcloud *Nextcloud
}

func NewServer(nextcloud *Nextcloud) *Server {
	return &Server{
		router:    chi.NewRouter(),
		nextcloud: nextcloud,
	}
}

func (s *Server) Start(addr string) error {
	s.RegisterRoutes()

	return http.ListenAndServe(addr, s.router)
}

func (s *Server) RegisterRoutes() {
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
	logger.Info("User created event received")

	err := s.nextcloud.ProvisionUser(event.Data.Email)
	if err != nil {
		logger.Error("Failed to provision user in Nextcloud")
	}
}

func (s *Server) handleUserDeleted(event UserEvent) {
	logger.Info("User deleted event received")

	err := s.nextcloud.DeprovisionUser(event.Data.Email)
	if err != nil {
		logger.Error("Failed to deprovision user in Nextcloud")
	}
}
