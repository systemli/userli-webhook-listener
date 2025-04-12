package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Nextcloud struct {
	client     *http.Client
	apiUrl     *url.URL
	username   string
	password   string
	providerID string
}

func NewNextcloud(apiUrl, username, password, providerID string) (*Nextcloud, error) {
	parsedURL, err := url.Parse(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Nextcloud URL: %v", err)
	}

	client := &http.Client{}

	return &Nextcloud{
		client:     client,
		apiUrl:     parsedURL,
		username:   username,
		password:   password,
		providerID: providerID,
	}, nil
}

func (n *Nextcloud) ProvisionUser(email string) error {
	// Extract the userId from the email
	userId := strings.Split(email, "@")[0]

	body := map[string]any{
		"providerId":  n.providerID,
		"userId":      userId,
		"email":       userId,
		"displayName": userId,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.apiUrl.String(), bytes.NewBuffer(jsonData))
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

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", n.apiUrl.String(), userId), nil)
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
	req.SetBasicAuth(n.username, n.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "UserliWebhookListener/1.0")
	req.Header.Set("ocs-apirequest", "true")
}
