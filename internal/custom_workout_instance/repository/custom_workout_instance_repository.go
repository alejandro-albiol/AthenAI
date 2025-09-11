package repository

import (
	"database/sql"
	"fmt"
	"time"

	exerciseDTO "github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"
	"github.com/lib/pq"
)

type CustomWorkoutInstanceRepositoryImpl struct {
	DB *sql.DB
}

func NewCustomWorkoutInstanceRepository(db *sql.DB) *CustomWorkoutInstanceRepositoryImpl {
	return &CustomWorkoutInstanceRepositoryImpl{DB: db}
}

func (r *CustomWorkoutInstanceRepositoryImpl) Create(gymID string, createdBy string, instance *dto.CreateCustomWorkoutInstanceDTO) (*string, error) {
	schema := pq.QuoteIdentifier(gymID)
	query := fmt.Sprintf(`
		INSERT INTO %s.custom_workout_instance (
			created_by, name, description, template_source, public_template_id, gym_template_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) 
		RETURNING id`, schema)

	var id string
	err := r.DB.QueryRow(query,
		createdBy,
		instance.Name,
		instance.Description,
		instance.TemplateSource,
		instance.PublicTemplateID,
		instance.GymTemplateID,
	).Scan(&id)

	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *CustomWorkoutInstanceRepositoryImpl) GetByID(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error) {
	// First get the basic workout instance
	instance, err := r.getBasicWorkoutInstance(gymID, id)
	if err != nil {
		return nil, err
	}

	// Get exercises and calculate stats
	exercises, err := r.getWorkoutExercises(gymID, id)
	if err != nil {
		return nil, err
	}

	// Calculate all the dynamic fields
	r.calculateWorkoutStats(instance, exercises)
	instance.Exercises = exercises

	return instance, nil
}

func (r *CustomWorkoutInstanceRepositoryImpl) GetSummaryByID(gymID, id string) (*dto.SummaryCustomWorkoutInstanceDTO, error) {
	// Get basic instance
	instance, err := r.getBasicWorkoutInstance(gymID, id)
	if err != nil {
		return nil, err
	}

	// Get exercises for calculations (but don't include in response)
	exercises, err := r.getWorkoutExercises(gymID, id)
	if err != nil {
		return nil, err
	}

	// Calculate stats
	r.calculateWorkoutStats(instance, exercises)

	// Convert to summary
	summary := &dto.SummaryCustomWorkoutInstanceDTO{
		ID:                       instance.ID,
		Name:                     instance.Name,
		Description:              instance.Description,
		TemplateSource:           instance.TemplateSource,
		DifficultyLevel:          instance.DifficultyLevel,
		EstimatedDurationMinutes: instance.EstimatedDurationMinutes,
		TotalExercises:           instance.TotalExercises,
		TotalSets:                instance.TotalSets,
		CreatedAt:                instance.CreatedAt,
		UpdatedAt:                instance.UpdatedAt,
	}

	// Set primary exercise type and muscle groups
	if len(instance.ExerciseTypes) > 0 {
		summary.PrimaryExerciseType = instance.ExerciseTypes[0]
	}
	if len(instance.MuscularGroups) > 3 {
		summary.PrimaryMuscularGroups = instance.MuscularGroups[:3]
	} else {
		summary.PrimaryMuscularGroups = instance.MuscularGroups
	}

	return summary, nil
}

// Helper method to get basic workout instance without calculated fields
func (r *CustomWorkoutInstanceRepositoryImpl) getBasicWorkoutInstance(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error) {
	schema := pq.QuoteIdentifier(gymID)
	query := fmt.Sprintf(`
		SELECT id, name, description, template_source, public_template_id, gym_template_id, created_at, updated_at
		FROM %s.custom_workout_instance 
		WHERE id = $1`, schema)

	instance := &dto.ResponseCustomWorkoutInstanceDTO{}
	var createdAt, updatedAt time.Time

	err := r.DB.QueryRow(query, id).Scan(
		&instance.ID,
		&instance.Name,
		&instance.Description,
		&instance.TemplateSource,
		&instance.PublicTemplateID,
		&instance.GymTemplateID,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	instance.CreatedAt = createdAt.Format(time.RFC3339)
	instance.UpdatedAt = updatedAt.Format(time.RFC3339)

	return instance, nil
}

// Helper method to get workout exercises with exercise details
func (r *CustomWorkoutInstanceRepositoryImpl) getWorkoutExercises(gymID, workoutInstanceID string) ([]exerciseDTO.ResponseCustomWorkoutExerciseDTO, error) {
	schema := pq.QuoteIdentifier(gymID)
	query := fmt.Sprintf(`
		SELECT 
			cwe.id, cwe.created_by, cwe.workout_instance_id, cwe.exercise_source,
			cwe.public_exercise_id, cwe.gym_exercise_id, cwe.block_name, cwe.exercise_order,
			cwe.sets, cwe.reps_min, cwe.reps_max, cwe.weight_kg, cwe.duration_seconds,
			cwe.rest_seconds, cwe.notes, cwe.created_at, cwe.updated_at
		FROM %s.custom_workout_exercise cwe
		WHERE cwe.workout_instance_id = $1
		ORDER BY cwe.block_name, cwe.exercise_order`, schema)

	rows, err := r.DB.Query(query, workoutInstanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []exerciseDTO.ResponseCustomWorkoutExerciseDTO
	for rows.Next() {
		var exercise exerciseDTO.ResponseCustomWorkoutExerciseDTO
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&exercise.ID,
			&exercise.CreatedBy,
			&exercise.WorkoutInstanceID,
			&exercise.ExerciseSource,
			&exercise.PublicExerciseID,
			&exercise.GymExerciseID,
			&exercise.BlockName,
			&exercise.ExerciseOrder,
			&exercise.Sets,
			&exercise.RepsMin,
			&exercise.RepsMax,
			&exercise.WeightKg,
			&exercise.DurationSeconds,
			&exercise.RestSeconds,
			&exercise.Notes,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		exercise.CreatedAt = createdAt.Format(time.RFC3339)
		exercise.UpdatedAt = updatedAt.Format(time.RFC3339)
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

// Helper method to calculate workout statistics and populate calculated fields
func (r *CustomWorkoutInstanceRepositoryImpl) calculateWorkoutStats(instance *dto.ResponseCustomWorkoutInstanceDTO, exercises []exerciseDTO.ResponseCustomWorkoutExerciseDTO) {
	if len(exercises) == 0 {
		instance.DifficultyLevel = "beginner"
		instance.EstimatedDurationMinutes = 0
		instance.TotalExercises = 0
		instance.TotalSets = 0
		instance.ExerciseTypes = []string{}
		instance.MuscularGroups = []string{}
		instance.EquipmentNeeded = []string{}
		return
	}

	// Calculate basic counts
	instance.TotalExercises = len(exercises)
	totalSets := 0
	totalEstimatedWeight := 0.0
	totalReps := 0
	repCount := 0
	estimatedDuration := 0

	// Maps for unique values and counting
	exerciseTypeMap := make(map[string]int)
	muscularGroupMap := make(map[string]int)
	equipmentMap := make(map[string]bool)
	difficultyMap := make(map[string]int)

	// Get exercise details from public/gym exercises and calculate stats
	for _, exercise := range exercises {
		// Count sets
		if exercise.Sets != nil {
			totalSets += *exercise.Sets
		}

		// Calculate estimated weight
		if exercise.WeightKg != nil && exercise.Sets != nil {
			totalEstimatedWeight += float64(*exercise.Sets) * *exercise.WeightKg
		}

		// Calculate average reps
		if exercise.RepsMin != nil && exercise.RepsMax != nil {
			avgReps := (*exercise.RepsMin + *exercise.RepsMax) / 2
			totalReps += avgReps
			repCount++
		}

		// Calculate duration
		if exercise.DurationSeconds != nil {
			estimatedDuration += *exercise.DurationSeconds
		}
		if exercise.RestSeconds != nil {
			estimatedDuration += *exercise.RestSeconds
		}

		// TODO: Get exercise details from public.exercise or gym.custom_exercise
		// For now, we'll use placeholder logic. In a real implementation,
		// you'd query the exercise tables to get difficulty, type, muscle groups, equipment

		// Placeholder difficulty calculation (could be enhanced with actual exercise data)
		if exercise.WeightKg != nil && *exercise.WeightKg > 50 {
			difficultyMap["advanced"]++
		} else if exercise.WeightKg != nil && *exercise.WeightKg > 20 {
			difficultyMap["intermediate"]++
		} else {
			difficultyMap["beginner"]++
		}
	}

	instance.TotalSets = totalSets

	// Calculate difficulty level (most common difficulty)
	maxCount := 0
	instance.DifficultyLevel = "beginner"
	for difficulty, count := range difficultyMap {
		if count > maxCount {
			maxCount = count
			instance.DifficultyLevel = difficulty
		}
	}

	// Convert duration from seconds to minutes
	instance.EstimatedDurationMinutes = estimatedDuration / 60

	// Convert maps to slices
	instance.ExerciseTypes = make([]string, 0, len(exerciseTypeMap))
	for exerciseType := range exerciseTypeMap {
		instance.ExerciseTypes = append(instance.ExerciseTypes, exerciseType)
	}

	instance.MuscularGroups = make([]string, 0, len(muscularGroupMap))
	for muscularGroup := range muscularGroupMap {
		instance.MuscularGroups = append(instance.MuscularGroups, muscularGroup)
	}

	instance.EquipmentNeeded = make([]string, 0, len(equipmentMap))
	for equipment := range equipmentMap {
		instance.EquipmentNeeded = append(instance.EquipmentNeeded, equipment)
	}

	// Calculate workout stats
	instance.WorkoutStats = &dto.WorkoutStatsDTO{
		TotalEstimatedWeight:   totalEstimatedWeight,
		DifficultyBreakdown:    difficultyMap,
		ExerciseTypeBreakdown:  exerciseTypeMap,
		AverageSetsPerExercise: float64(totalSets) / float64(len(exercises)),
	}

	if repCount > 0 {
		instance.WorkoutStats.AverageRepsPerSet = float64(totalReps) / float64(repCount)
	}
}

func (r *CustomWorkoutInstanceRepositoryImpl) GetByUserID(gymID, userID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	schema := pq.QuoteIdentifier(gymID)
	query := fmt.Sprintf(`
		SELECT id, name, description, template_source, public_template_id, gym_template_id, created_at, updated_at
		FROM %s.custom_workout_instance 
		WHERE created_by = $1 
		ORDER BY created_at DESC`, schema)

	return r.getWorkoutInstanceList(gymID, query, userID)
}

func (r *CustomWorkoutInstanceRepositoryImpl) GetSummariesByUserID(gymID, userID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error) {
	instances, err := r.GetByUserID(gymID, userID)
	if err != nil {
		return nil, err
	}

	summaries := make([]*dto.SummaryCustomWorkoutInstanceDTO, len(instances))
	for i, instance := range instances {
		summary := &dto.SummaryCustomWorkoutInstanceDTO{
			ID:                       instance.ID,
			Name:                     instance.Name,
			Description:              instance.Description,
			TemplateSource:           instance.TemplateSource,
			DifficultyLevel:          instance.DifficultyLevel,
			EstimatedDurationMinutes: instance.EstimatedDurationMinutes,
			TotalExercises:           instance.TotalExercises,
			TotalSets:                instance.TotalSets,
			CreatedAt:                instance.CreatedAt,
			UpdatedAt:                instance.UpdatedAt,
		}

		if len(instance.ExerciseTypes) > 0 {
			summary.PrimaryExerciseType = instance.ExerciseTypes[0]
		}
		if len(instance.MuscularGroups) > 3 {
			summary.PrimaryMuscularGroups = instance.MuscularGroups[:3]
		} else {
			summary.PrimaryMuscularGroups = instance.MuscularGroups
		}

		summaries[i] = summary
	}

	return summaries, nil
}

func (r *CustomWorkoutInstanceRepositoryImpl) GetLastsByUserID(gymID, userID string, numberOfWorkouts int) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	schema := pq.QuoteIdentifier(gymID)
	query := fmt.Sprintf(`
		SELECT id, name, description, template_source, public_template_id, gym_template_id, created_at, updated_at
		FROM %s.custom_workout_instance 
		WHERE created_by = $1 
		ORDER BY created_at DESC 
		LIMIT $2`, schema)

	return r.getWorkoutInstanceList(gymID, query, userID, numberOfWorkouts)
}

func (r *CustomWorkoutInstanceRepositoryImpl) List(gymID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	schema := pq.QuoteIdentifier(gymID)
	query := fmt.Sprintf(`
		SELECT id, name, description, template_source, public_template_id, gym_template_id, created_at, updated_at
		FROM %s.custom_workout_instance 
		ORDER BY created_at DESC`, schema)

	return r.getWorkoutInstanceList(gymID, query)
}

func (r *CustomWorkoutInstanceRepositoryImpl) ListSummaries(gymID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error) {
	instances, err := r.List(gymID)
	if err != nil {
		return nil, err
	}

	summaries := make([]*dto.SummaryCustomWorkoutInstanceDTO, len(instances))
	for i, instance := range instances {
		summary := &dto.SummaryCustomWorkoutInstanceDTO{
			ID:                       instance.ID,
			Name:                     instance.Name,
			Description:              instance.Description,
			TemplateSource:           instance.TemplateSource,
			DifficultyLevel:          instance.DifficultyLevel,
			EstimatedDurationMinutes: instance.EstimatedDurationMinutes,
			TotalExercises:           instance.TotalExercises,
			TotalSets:                instance.TotalSets,
			CreatedAt:                instance.CreatedAt,
			UpdatedAt:                instance.UpdatedAt,
		}

		if len(instance.ExerciseTypes) > 0 {
			summary.PrimaryExerciseType = instance.ExerciseTypes[0]
		}
		if len(instance.MuscularGroups) > 3 {
			summary.PrimaryMuscularGroups = instance.MuscularGroups[:3]
		} else {
			summary.PrimaryMuscularGroups = instance.MuscularGroups
		}

		summaries[i] = summary
	}

	return summaries, nil
}

// Helper method to get workout instance list and calculate stats
func (r *CustomWorkoutInstanceRepositoryImpl) getWorkoutInstanceList(gymID, query string, args ...interface{}) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []*dto.ResponseCustomWorkoutInstanceDTO
	for rows.Next() {
		instance := &dto.ResponseCustomWorkoutInstanceDTO{}
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&instance.ID,
			&instance.Name,
			&instance.Description,
			&instance.TemplateSource,
			&instance.PublicTemplateID,
			&instance.GymTemplateID,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		instance.CreatedAt = createdAt.Format(time.RFC3339)
		instance.UpdatedAt = updatedAt.Format(time.RFC3339)

		// Get exercises and calculate stats for each instance
		exercises, err := r.getWorkoutExercises(gymID, instance.ID)
		if err != nil {
			return nil, err
		}

		r.calculateWorkoutStats(instance, exercises)
		instances = append(instances, instance)
	}

	return instances, nil
}

func (r *CustomWorkoutInstanceRepositoryImpl) Update(gymID, id string, instance *dto.UpdateCustomWorkoutInstanceDTO) error {
	schema := pq.QuoteIdentifier(gymID)
	query := fmt.Sprintf(`
		UPDATE %s.custom_workout_instance SET
			name = COALESCE($1, name),
			description = COALESCE($2, description),
			template_source = COALESCE($3, template_source),
			public_template_id = COALESCE($4, public_template_id),
			gym_template_id = COALESCE($5, gym_template_id),
			updated_at = NOW()
		WHERE id = $6`, schema)

	result, err := r.DB.Exec(query,
		instance.Name,
		instance.Description,
		instance.TemplateSource,
		instance.PublicTemplateID,
		instance.GymTemplateID,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *CustomWorkoutInstanceRepositoryImpl) Delete(gymID, id string) error {
	schema := pq.QuoteIdentifier(gymID)
	query := fmt.Sprintf(`DELETE FROM %s.custom_workout_instance WHERE id = $1`, schema)

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
