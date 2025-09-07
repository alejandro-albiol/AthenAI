package tests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
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

func extractUUIDFromJWT(token string, t *testing.T) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return "", io.ErrUnexpectedEOF
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", err
	}
	// Print the decoded JWT payload for debugging
	payloadJson, _ := json.MarshalIndent(claims, "", "  ")
	t.Logf("Decoded JWT payload: %s", payloadJson)
	// Try common claim keys
	if sub, ok := claims["sub"].(string); ok {
		return sub, nil
	}
	if uid, ok := claims["user_id"].(string); ok {
		return uid, nil
	}
	return "", io.ErrUnexpectedEOF
}

func loginAsAdminWithGymID(t *testing.T, client *http.Client, baseURL, email, password, gymID string) (string, string) {
	loginPayload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if gymID != "" {
		req.Header.Set("X-Gym-ID", gymID)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("login failed, status: %d. Response: %s", resp.StatusCode, string(b))
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
	uuid, err := extractUUIDFromJWT(result.Data.AccessToken, t)
	if err != nil {
		t.Fatalf("failed to extract admin UUID from token: %v", err)
	}
	return result.Data.AccessToken, uuid
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

	var gymID, equipmentID, muscularGroupID, exerciseID, userID, templateID, blockID string
	var headers map[string]string
	var adminUuid string
	var adminToken string

	// Cleanup logic: delete in reverse dependency order
	defer func() {
		deleteIfExists := func(url string) {
			req, _ := http.NewRequest("DELETE", url, nil)
			for k, v := range headers {
				req.Header.Set(k, v)
			}
			resp, err := client.Do(req)
			if err == nil {
				resp.Body.Close()
			}
		}
		if blockID != "" {
			deleteIfExists(baseURL + "/api/v1/template-block/" + blockID)
		}
		if templateID != "" {
			deleteIfExists(baseURL + "/api/v1/workout-template/" + templateID)
		}
		if userID != "" {
			deleteIfExists(baseURL + "/api/v1/users/" + userID)
		}
		if exerciseID != "" {
			deleteIfExists(baseURL + "/api/v1/exercise/" + exerciseID)
		}
		if muscularGroupID != "" {
			deleteIfExists(baseURL + "/api/v1/muscular-group/" + muscularGroupID)
		}
		if equipmentID != "" {
			deleteIfExists(baseURL + "/api/v1/equipment/" + equipmentID)
		}
		if gymID != "" {
			deleteIfExists(baseURL + "/api/v1/gym/" + gymID)
		}
	}()

	// 1. Login as admin (no gym ID) to get JWT for gym creation
	t.Run("Login as Platform Admin", func(t *testing.T) {
		token, uuid := loginAsAdminWithGymID(t, client, baseURL, adminEmail, adminPassword, "")
		adminToken = token
		adminUuid = uuid
		headers = map[string]string{"Authorization": "Bearer " + adminToken}
	})

	// 2. Create Gym using admin JWT
	t.Run("Create Gym", func(t *testing.T) {
		payload := map[string]any{
			"name":    "E2E Gym " + time.Now().Format("20060102150405"),
			"address": "E2E Street",
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
			"created_by":  adminUuid,
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
			"name":        "E2E Muscular Group " + time.Now().Format("20060102150405"),
			"description": "E2E test group",
			"body_part":   "full_body",
			"created_by":  adminUuid,
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
			"name":               "E2E Exercise " + time.Now().Format("20060102150405"),
			"description":        "E2E test exercise",
			"difficulty_level":   "beginner",
			"exercise_type":      "strength",
			"synonyms":           []string{"E2E Exo", "E2E Move"},
			"equipment_ids":      []string{equipmentID},
			"muscular_group_ids": []string{muscularGroupID},
			"created_by":         adminUuid,
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

	// After gym creation, use the admin token to register a user for this gym.
	// The request must include the gym UUID in the X-Gym-ID header.
	t.Run("Create User in Gym", func(t *testing.T) {
		payload := map[string]any{
			"username":          "e2euser_" + gymID,
			"email":             "e2euser+" + gymID + "@test.com",
			"password":          "E2Epass123!",
			"role":              "user",
			"description":       "E2E test user",
			"training_phase":    "maintenance",
			"motivation":        "wellbeing",
			"special_situation": "none",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/user", bytes.NewBuffer(body))
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		req.Header.Set("X-Gym-ID", gymID) // Add gym context for registration
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
			"name":             "E2E Template",
			"description":      "E2E test template",
			"difficulty_level": "beginner",
			"created_by":       adminUuid,
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
			t.Fatalf("expected 201, got %d. Response: %s", resp.StatusCode, string(b))
		}
		var res struct {
			Data struct {
				ID string `json:"id"`
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode template creation response: %v", err)
		}
		if res.Data.ID == "" {
			t.Fatalf("no template ID returned. Full response: %+v", res)
		}
		templateID = res.Data.ID
	})

	// Ensure IDs are set before dependent calls
	if templateID == "" {
		t.Fatal("templateID is empty before creating template block")
	}
	if exerciseID == "" {
		t.Fatal("exerciseID is empty before creating template block")
	}

	t.Run("Add Block to Template", func(t *testing.T) {
		// Ensure templateID and exerciseID are set and not empty
		if templateID == "" {
			t.Fatal("templateID is empty before creating template block")
		}
		if exerciseID == "" {
			t.Fatal("exerciseID is empty before creating template block")
		}
		payload := map[string]any{
			"template_id":    templateID,
			"block_name":     "E2E Block",
			"block_type":     "main",
			"block_order":    1,
			"exercise_count": 1,
			"exercise_id":    exerciseID,
			"created_by":     adminUuid,
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
