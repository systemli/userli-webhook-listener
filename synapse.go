package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SynapseConfig struct {
	ApiUrl      string
	AccessToken string
	Domain      string
}

type Synapse struct {
	client *http.Client
	config *SynapseConfig
}

func NewSynapse(cfg *SynapseConfig) *Synapse {
	client := &http.Client{}

	return &Synapse{
		client: client,
		config: cfg,
	}
}

func (s *Synapse) DeprovisionUser(email string) error {
	// Extract the userId from the email
	userId := strings.Split(email, "@")[0]
	domain := strings.Split(email, "@")[1]

	if domain != s.config.Domain {
		return fmt.Errorf("domain not allowed: %s", domain)
	}

	body := map[string]any{
		"erase": true,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/deactivate/@%s:%s", s.config.ApiUrl, userId, domain), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	s.prepareRequest(req)
	res, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to deprovision user: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to deprovision user, status code: %d", res.StatusCode)
	}

	return nil
}

func (s *Synapse) prepareRequest(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", s.config.AccessToken))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "UserliWebhookListener/1.0")
	req.Header.Set("ocs-apirequest", "true")
}
