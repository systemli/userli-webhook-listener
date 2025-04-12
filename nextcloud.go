package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type NextcloudConfig struct {
	ApiUrl     string
	Username   string
	Password   string
	ProviderID string
	Domain     string
}

type Nextcloud struct {
	client *http.Client
	config *NextcloudConfig
}

func NewNextcloud(cfg *NextcloudConfig) *Nextcloud {
	client := &http.Client{}

	return &Nextcloud{
		client: client,
		config: cfg,
	}
}

func (n *Nextcloud) ProvisionUser(email string) error {
	// Extract the userId from the email
	userId := strings.Split(email, "@")[0]
	domain := strings.Split(email, "@")[1]

	if domain != n.config.Domain {
		return fmt.Errorf("domain not allowed: %s", domain)
	}

	body := map[string]any{
		"providerId":  n.config.ProviderID,
		"userId":      userId,
		"email":       userId,
		"displayName": userId,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.config.ApiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	n.prepareRequest(req)
	res, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to provision user: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to provision user, status code: %d", res.StatusCode)
	}

	return nil
}

func (n *Nextcloud) DeprovisionUser(email string) error {
	// Extract the userId from the email
	userId := strings.Split(email, "@")[0]
	domain := strings.Split(email, "@")[1]

	if domain != n.config.Domain {
		return fmt.Errorf("domain not allowed: %s", domain)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", n.config.ApiUrl, userId), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	n.prepareRequest(req)
	res, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to deprovision user: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to deprovision user, status code: %d", res.StatusCode)
	}

	return nil
}

func (n *Nextcloud) prepareRequest(req *http.Request) {
	req.SetBasicAuth(n.config.Username, n.config.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "UserliWebhookListener/1.0")
	req.Header.Set("ocs-apirequest", "true")
}
