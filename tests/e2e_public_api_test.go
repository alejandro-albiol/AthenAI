package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func loadEnvOrSkip(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Skip(".env file not found, skipping E2E test")
	}
}

func getAdminCredentials() (string, string) {
	return os.Getenv("ADMIN_EMAIL"), os.Getenv("ADMIN_PASSWORD")
}

// Setup/teardown helpers, test data, and E2E test functions will go here.

func loginAsAdmin(t *testing.T, client *http.Client, baseURL, email, password string) string {
	loginPayload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(loginPayload)
	resp, err := client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("login failed, status: %d", resp.StatusCode)
	}
	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			AccessToken string `json:"access_token"`
			// ...other fields if needed
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode login response: %v", err)
	}
	if result.Data.AccessToken == "" {
		t.Fatal("no access_token returned from login")
	}
	return result.Data.AccessToken
}

func TestPublicAPI_EndToEnd(t *testing.T) {
	loadEnvOrSkip(t)
	adminEmail, adminPassword := getAdminCredentials()
	baseURL := os.Getenv("E2E_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080" // default
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// 1. Login as admin
	token := loginAsAdmin(t, client, baseURL, adminEmail, adminPassword)
	headers := map[string]string{"Authorization": "Bearer " + token}

	// 2. Create equipment
	eqPayload := map[string]any{
		"name":        "E2E Equipment " + fmt.Sprintf("%d", time.Now().UnixNano()),
		"description": "E2E test equipment",
		"category":    "free_weights",
		"created_by":  adminEmail,
	}
	eqBody, _ := json.Marshal(eqPayload)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/equipment", bytes.NewBuffer(eqBody))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("create equipment failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}
	var createRes struct {
		Data struct {
			ID string `json:"id"`
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&createRes); err != nil {
		t.Fatalf("decode create equipment: %v", err)
	}
	if createRes.Data.ID == "" {
		t.Fatal("no equipment ID returned")
	}
	eqID := createRes.Data.ID

	// 3. Get equipment by ID
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/equipment/%s", baseURL, eqID), nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("get equipment failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// 4. Update equipment
	updPayload := map[string]any{"description": "Updated by E2E"}
	updBody, _ := json.Marshal(updPayload)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/equipment/%s", baseURL, eqID), bytes.NewBuffer(updBody))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("update equipment failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// 5. Delete equipment
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/equipment/%s", baseURL, eqID), nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("delete equipment failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// TODO: Repeat similar CRUD for exercise, gym, user, etc.
}
