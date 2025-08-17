package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/workout_generator/dto"
)

// WorkoutGeneratorService connects to the LLM API to generate workouts
type WorkoutGeneratorService struct {
	LLMEndpoint string // e.g. Hugging Face Inference API URL
	APIToken    string // e.g. Bearer token for Hugging Face
}

func NewWorkoutGeneratorService(llmEndpoint, apiToken string) *WorkoutGeneratorService {
	return &WorkoutGeneratorService{
		LLMEndpoint: llmEndpoint,
		APIToken:    apiToken,
	}
}

func (s *WorkoutGeneratorService) GenerateWorkout(req dto.WorkoutGeneratorRequest) (*dto.WorkoutGeneratorResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.LLMEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.APIToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM API returned status: %s", resp.Status)
	}

	var workoutResp dto.WorkoutGeneratorResponse
	if err := json.NewDecoder(resp.Body).Decode(&workoutResp); err != nil {
		return nil, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	return &workoutResp, nil
}
