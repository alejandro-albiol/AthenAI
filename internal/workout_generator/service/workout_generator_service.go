package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	exdto "github.com/alejandro-albiol/athenai/internal/exercise/dto"
	exerciseIF "github.com/alejandro-albiol/athenai/internal/exercise/interfaces"
	tbdto "github.com/alejandro-albiol/athenai/internal/template_block/dto"
	templateIF "github.com/alejandro-albiol/athenai/internal/template_block/interfaces"
	userIF "github.com/alejandro-albiol/athenai/internal/user/interfaces"
	"github.com/alejandro-albiol/athenai/internal/workout_generator/dto"
	wtdto "github.com/alejandro-albiol/athenai/internal/workout_template/dto"
	workoutTemplateIF "github.com/alejandro-albiol/athenai/internal/workout_template/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type WorkoutGeneratorService struct {
	LLMEndpoint     string
	APIToken        string
	ExerciseService exerciseIF.ExerciseService
	TemplateService workoutTemplateIF.WorkoutTemplateService
	BlockService    templateIF.TemplateBlockService
	UserService     userIF.UserService
}

func NewWorkoutGeneratorService(
	llmEndpoint, apiToken string,
	exerciseService exerciseIF.ExerciseService,
	templateService workoutTemplateIF.WorkoutTemplateService,
	blockService templateIF.TemplateBlockService,
	userService userIF.UserService,
) *WorkoutGeneratorService {
	return &WorkoutGeneratorService{
		LLMEndpoint:     llmEndpoint,
		APIToken:        apiToken,
		ExerciseService: exerciseService,
		TemplateService: templateService,
		BlockService:    blockService,
		UserService:     userService,
	}
}

func (s *WorkoutGeneratorService) GenerateWorkout(req dto.WorkoutGeneratorRequest) (*dto.WorkoutGeneratorResponse, error) {
	// 1. Get user context
	user, err := s.UserService.GetUserByID("", req.UserID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "User not found", err)
	}

	// 2. Get available exercises (public + tenant)
	publicExercises, err := s.ExerciseService.GetAllExercises()
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get public exercises", err)
	}
	// For tenant-specific, you may need to pass gymID or tenantID, here assumed as user.GymID
	var tenantExercises []exdto.ExerciseResponseDTO
	if user.GymID != "" {
		// TODO: implement GetTenantExercises if needed
		// tenantExercises, err = s.ExerciseService.GetTenantExercises(user.GymID)
		// if err != nil {
		//     return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get tenant exercises", err)
		// }
	}
	exercises := mergeExercises(publicExercises, tenantExercises)

	// 3. Get workout template
	template, err := s.TemplateService.GetWorkoutTemplateByName(req.TemplateName)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "Workout template not found", err)
	}
	blocks, err := s.BlockService.ListTemplateBlocksByTemplateID(template.ID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get template blocks", err)
	}

	// 4. Build LLM prompt
	prompt := buildPrompt(req, exercises, template, blocks)

	// 5. Call LLM
	llmReq := map[string]interface{}{"prompt": prompt}
	payload, err := json.Marshal(llmReq)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to marshal LLM request", err)
	}
	httpReq, err := http.NewRequest("POST", s.LLMEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to create HTTP request", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.APIToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to call LLM API", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apierror.New(errorcode_enum.CodeInternal, "LLM API returned status: "+resp.Status, nil)
	}

	var workoutResp dto.WorkoutGeneratorResponse
	if err := json.NewDecoder(resp.Body).Decode(&workoutResp); err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to decode LLM response", err)
	}
	return &workoutResp, nil
}

// mergeExercises combines public and tenant exercises, prioritizing tenant-specific
func mergeExercises(public, tenant []exdto.ExerciseResponseDTO) []exdto.ExerciseResponseDTO {
	// TODO: implement deduplication/override logic if needed
	return append(public, tenant...)
}

// buildPrompt creates a rich prompt for the LLM
func buildPrompt(req dto.WorkoutGeneratorRequest, exercises []exdto.ExerciseResponseDTO, template wtdto.WorkoutTemplateDTO, blocks []tbdto.TemplateBlockDTO) string {
	var sb strings.Builder
	sb.WriteString("User Context:\n")
	sb.WriteString(fmt.Sprintf("ID: %s\n", req.UserID))
	sb.WriteString(fmt.Sprintf("Tags: %v\n", req.Tags))
	sb.WriteString(fmt.Sprintf("Training Phase: %s\n", req.TrainingPhase))
	sb.WriteString(fmt.Sprintf("Motivation: %s\n", req.Motivation))
	sb.WriteString(fmt.Sprintf("Special Situation: %s\n", req.SpecialSituation))
	sb.WriteString("\nAvailable Exercises:\n")
	for _, ex := range exercises {
		sb.WriteString(fmt.Sprintf("- %s (%s)\n", ex.Name, ex.ExerciseType))
	}
	sb.WriteString("\nWorkout Template:\n")
	sb.WriteString(fmt.Sprintf("Name: %s\n", template.Name))
	sb.WriteString(fmt.Sprintf("Description: %s\n", derefString(template.Description)))
	sb.WriteString(fmt.Sprintf("Difficulty: %s\n", template.DifficultyLevel))
	sb.WriteString("Blocks:\n")
	for _, block := range blocks {
		sb.WriteString(fmt.Sprintf("- %s (%s): %d exercises\n", block.Name, block.Type, block.ExerciseCount))
	}
	sb.WriteString("\nPlease generate a workout plan for this user based on the above context, available exercises, and template structure.")
	return sb.String()
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
