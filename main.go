package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/user", handleUserCreated).Methods(http.MethodPost)
	r.HandleFunc("/user/{email}", handleUserDeleted).Methods(http.MethodDelete)

	http.Handle("/", r)

	fmt.Println("Listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handleUserCreated(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract the email parameter
	email := r.PostFormValue("email")
	if email == "" {
		http.Error(w, "Missing email parameter", http.StatusBadRequest)
		return
	}

	fmt.Println("User created event received")

	// Send create request to the Nextcloud OIDC user API
	resp, err := provisionNextcloudUser(email)
	if err != nil {
		http.Error(w, "Failed to provision user in Nextcloud", http.StatusInternalServerError)
	} else {
		w.WriteHeader(resp.StatusCode)
	}
}

func handleUserDeleted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	if email == "" {
		http.Error(w, "Missing email parameter", http.StatusBadRequest)
		return
	}

	fmt.Println("User deleted event received")

	// Send delete request to the Nextcloud OIDC user API
	resp, err := deprovisionNextcloudUser(email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to deprovision user in Nextcloud. Error: %s", err), http.StatusInternalServerError)
	} else {
		w.WriteHeader(resp.StatusCode)
	}
}

func provisionNextcloudUser(email string) (*http.Response, error) {
	// Extract the userId from the email
	userId := strings.Split(email, "@")[0]

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, os.Getenv("NEXTCLOUD_OIDC_USER_API"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	prepareNextcloudRequest(req)
	req.Header.Set("Content-Type", "application/json")

	data := url.Values{}
	data.Set("userId", userId)
	data.Set("email", userId)
	data.Set("displayName", userId)
	req.URL.RawQuery = data.Encode()

	return client.Do(req)
}

func deprovisionNextcloudUser(email string) (*http.Response, error) {
	// Extract the userId from the email
	userId := strings.Split(email, "@")[0]

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", os.Getenv("NEXTCLOUD_OIDC_USER_API"), userId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	prepareNextcloudRequest(req)

	return client.Do(req)
}

func prepareNextcloudRequest(req *http.Request) {
	// Create the basic auth header
	auth := fmt.Sprintf("%s:%s", os.Getenv("NEXTCLOUD_ADMIN_USERNAME"), os.Getenv("NEXTCLOUD_ADMIN_PASSWORD"))
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encoded)
	req.Header.Set("ocs-apirequest", "true")
}
