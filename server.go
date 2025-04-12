package main

import (
	"encoding/json"
	"fmt"
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
	s.router.Post("/user", s.handleUserCreated)
	s.router.Delete("/user/{email}", s.handleUserDeleted)
}

func (s *Server) handleUserCreated(w http.ResponseWriter, r *http.Request) {
	logger.Info("User created event received")

	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.nextcloud.ProvisionUser(body.Email)
	if err != nil {
		http.Error(w, "Failed to provision user in Nextcloud", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUserDeleted(w http.ResponseWriter, r *http.Request) {
	logger.Info("User deleted event received")

	email := chi.URLParam(r, "email")

	err := s.nextcloud.DeprovisionUser(email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to deprovision user in Nextcloud. Error: %s", err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
