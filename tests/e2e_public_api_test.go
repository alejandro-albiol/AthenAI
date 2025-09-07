package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"

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

	var gymID, equipmentID, muscularGroupID, exerciseID, userID, templateID, blockID string

	t.Run("Create Gym", func(t *testing.T) {
		payload := map[string]any{
			"name":       "E2E Gym",
			"address":    "E2E Street",
			"created_by": adminEmail,
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/gym", bytes.NewBuffer(body))
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("create gym failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 201 {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("expected 201, got %d. %s", resp.StatusCode, string(b))
		}
		var res struct {
			Data struct {
				ID string `json:"id"`
			}
		}
		_ = json.NewDecoder(resp.Body).Decode(&res)
		if res.Data.ID == "" {
			t.Fatal("no gym ID returned")
		}
		gymID = res.Data.ID
	})

	t.Run("Create Equipment", func(t *testing.T) {
		payload := map[string]any{
			"name":        "E2E Equipment",
			"description": "E2E test equipment",
			"category":    "free_weights",
			"created_by":  adminEmail,
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/equipment", bytes.NewBuffer(body))
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("create equipment failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 201 {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("expected 201, got %d. %s", resp.StatusCode, string(b))
		}
		var res struct {
			Data struct {
				ID string `json:"id"`
			}
		}
		_ = json.NewDecoder(resp.Body).Decode(&res)
		if res.Data.ID == "" {
			t.Fatal("no equipment ID returned")
		}
		equipmentID = res.Data.ID
	})

	t.Run("Create Muscular Group", func(t *testing.T) {
		payload := map[string]any{
			"name":        "E2E Muscular Group",
			"description": "E2E test group",
			"body_part":   "full_body",
			"created_by":  adminEmail,
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/muscular-group", bytes.NewBuffer(body))
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("create muscular group failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 201 {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("expected 201, got %d. %s", resp.StatusCode, string(b))
		}
		var res struct {
			Data struct {
				ID string `json:"id"`
			}
		}
		_ = json.NewDecoder(resp.Body).Decode(&res)
		if res.Data.ID == "" {
			t.Fatal("no muscular group ID returned")
		}
		muscularGroupID = res.Data.ID
	})

	t.Run("Create Exercise", func(t *testing.T) {
		payload := map[string]any{
			"name":               "E2E Exercise",
			"description":        "E2E test exercise",
			"difficulty_level":    "beginner",
			"exercise_type":       "strength",
			"equipment_ids":      []string{equipmentID},
			"muscular_group_ids": []string{muscularGroupID},
			"created_by":         adminEmail,
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/exercise", bytes.NewBuffer(body))
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("create exercise failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 201 {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("expected 201, got %d. %s", resp.StatusCode, string(b))
		}
		var res struct {
			Data struct {
				ID string `json:"id"`
			}
		}
		_ = json.NewDecoder(resp.Body).Decode(&res)
		if res.Data.ID == "" {
			t.Fatal("no exercise ID returned")
		}
		exerciseID = res.Data.ID
	})

	t.Run("Create User in Gym", func(t *testing.T) {
		payload := map[string]any{
			"email":      "e2euser+" + gymID + "@test.com",
			"password":   "E2Epass123!",
			"first_name": "E2E",
			"last_name":  "User",
			"role":       "member",
			"created_by": adminEmail,
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/user", bytes.NewBuffer(body))
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		req.Header.Set("X-Gym-ID", gymID)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("create user failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 201 {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("expected 201, got %d. %s", resp.StatusCode, string(b))
		}
		var res struct {
			Data struct {
				ID string `json:"id"`
			}
		}
		_ = json.NewDecoder(resp.Body).Decode(&res)
		if res.Data.ID == "" {
			t.Fatal("no user ID returned")
		}
		userID = res.Data.ID
	})

	t.Run("Create Workout Template", func(t *testing.T) {
		payload := map[string]any{
			"name":        "E2E Template",
			"description": "E2E test template",
			"created_by":  adminEmail,
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/workout-template", bytes.NewBuffer(body))
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("create template failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 201 {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("expected 201, got %d. %s", resp.StatusCode, string(b))
		}
		var res struct {
			Data struct {
				ID string `json:"id"`
			}
		}
		_ = json.NewDecoder(resp.Body).Decode(&res)
		if res.Data.ID == "" {
			t.Fatal("no template ID returned")
		}
		templateID = res.Data.ID
	})

	t.Run("Add Block to Template", func(t *testing.T) {
		payload := map[string]any{
			"workout_template_id": templateID,
			"exercise_id":         exerciseID,
			"order":               1,
			"created_by":          adminEmail,
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/template-block", bytes.NewBuffer(body))
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("create block failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 201 {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("expected 201, got %d. %s", resp.StatusCode, string(b))
		}
		var res struct {
			Data struct {
				ID string `json:"id"`
			}
		}
		_ = json.NewDecoder(resp.Body).Decode(&res)
		if res.Data.ID == "" {
			t.Fatal("no block ID returned")
		}
		blockID = res.Data.ID
	})

	_ = blockID
	_ = userID
}
