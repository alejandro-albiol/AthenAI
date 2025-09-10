package repository

import (
	"database/sql"
	"fmt"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"
)

type CustomWorkoutExerciseRepository struct {
	DB *sql.DB
}

func NewCustomWorkoutExerciseRepository(db *sql.DB) *CustomWorkoutExerciseRepository {
	return &CustomWorkoutExerciseRepository{DB: db}
}

func (r *CustomWorkoutExerciseRepository) Create(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (*string, error) {
	query := `INSERT INTO "%s".custom_workout_exercise 
		(created_by, workout_instance_id, exercise_source, public_exercise_id, gym_exercise_id, block_name, exercise_order, sets, reps_min, reps_max, weight_kg, duration_seconds, rest_seconds, notes) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) 
		RETURNING id`

	var id string
	err := r.DB.QueryRow(
		fmt.Sprintf(query, gymID),
		exercise.CreatedBy,
		exercise.WorkoutInstanceID,
		exercise.ExerciseSource,
		exercise.PublicExerciseID,
		exercise.GymExerciseID,
		exercise.BlockName,
		exercise.ExerciseOrder,
		exercise.Sets,
		exercise.RepsMin,
		exercise.RepsMax,
		exercise.WeightKg,
		exercise.DurationSeconds,
		exercise.RestSeconds,
		exercise.Notes,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *CustomWorkoutExerciseRepository) GetByID(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error) {
	query := `SELECT id, created_by, workout_instance_id, exercise_source, public_exercise_id, gym_exercise_id, block_name, 
		exercise_order, sets, reps_min, reps_max, weight_kg, duration_seconds, rest_seconds, notes, 
		created_at, updated_at 
		FROM "%s".custom_workout_exercise WHERE id = $1`

	row := r.DB.QueryRow(fmt.Sprintf(query, gymID), id)

	var res dto.ResponseCustomWorkoutExerciseDTO
	var createdAt, updatedAt sql.NullString

	err := row.Scan(
		&res.ID,
		&res.CreatedBy,
		&res.WorkoutInstanceID,
		&res.ExerciseSource,
		&res.PublicExerciseID,
		&res.GymExerciseID,
		&res.BlockName,
		&res.ExerciseOrder,
		&res.Sets,
		&res.RepsMin,
		&res.RepsMax,
		&res.WeightKg,
		&res.DurationSeconds,
		&res.RestSeconds,
		&res.Notes,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	if createdAt.Valid {
		res.CreatedAt = createdAt.String
	}
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.String
	}

	return &res, nil
}

func (r *CustomWorkoutExerciseRepository) ListByWorkoutInstanceID(gymID, workoutInstanceID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	query := `SELECT id, created_by, workout_instance_id, exercise_source, public_exercise_id, gym_exercise_id, block_name, 
		exercise_order, sets, reps_min, reps_max, weight_kg, duration_seconds, rest_seconds, notes, 
		created_at, updated_at 
		FROM "%s".custom_workout_exercise WHERE workout_instance_id = $1 ORDER BY block_name, exercise_order`

	rows, err := r.DB.Query(fmt.Sprintf(query, gymID), workoutInstanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*dto.ResponseCustomWorkoutExerciseDTO
	for rows.Next() {
		var res dto.ResponseCustomWorkoutExerciseDTO
		var createdAt, updatedAt sql.NullString

		err := rows.Scan(
			&res.ID,
			&res.CreatedBy,
			&res.WorkoutInstanceID,
			&res.ExerciseSource,
			&res.PublicExerciseID,
			&res.GymExerciseID,
			&res.BlockName,
			&res.ExerciseOrder,
			&res.Sets,
			&res.RepsMin,
			&res.RepsMax,
			&res.WeightKg,
			&res.DurationSeconds,
			&res.RestSeconds,
			&res.Notes,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			res.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			res.UpdatedAt = updatedAt.String
		}

		result = append(result, &res)
	}

	return result, nil
}

func (r *CustomWorkoutExerciseRepository) ListByMuscularGroupID(gymID, muscularGroupID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	query := `SELECT cwe.id, cwe.created_by, cwe.workout_instance_id, cwe.exercise_source, cwe.public_exercise_id, cwe.gym_exercise_id, 
		cwe.block_name, cwe.exercise_order, cwe.sets, cwe.reps_min, cwe.reps_max, cwe.weight_kg, 
		cwe.duration_seconds, cwe.rest_seconds, cwe.notes, cwe.created_at, cwe.updated_at 
		FROM "%s".custom_workout_exercise cwe
		LEFT JOIN public.exercise_muscular_group emg ON cwe.public_exercise_id = emg.exercise_id
		LEFT JOIN "%s".custom_exercise_muscular_group cemg ON cwe.gym_exercise_id = cemg.exercise_id
		WHERE emg.muscular_group_id = $1 OR cemg.muscular_group_id = $1
		ORDER BY cwe.block_name, cwe.exercise_order`

	rows, err := r.DB.Query(fmt.Sprintf(query, gymID, gymID), muscularGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*dto.ResponseCustomWorkoutExerciseDTO
	for rows.Next() {
		var res dto.ResponseCustomWorkoutExerciseDTO
		var createdAt, updatedAt sql.NullString

		err := rows.Scan(
			&res.ID,
			&res.CreatedBy,
			&res.WorkoutInstanceID,
			&res.ExerciseSource,
			&res.PublicExerciseID,
			&res.GymExerciseID,
			&res.BlockName,
			&res.ExerciseOrder,
			&res.Sets,
			&res.RepsMin,
			&res.RepsMax,
			&res.WeightKg,
			&res.DurationSeconds,
			&res.RestSeconds,
			&res.Notes,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			res.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			res.UpdatedAt = updatedAt.String
		}

		result = append(result, &res)
	}

	return result, nil
}

func (r *CustomWorkoutExerciseRepository) ListByEquipmentID(gymID, equipmentID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	query := `SELECT cwe.id, cwe.created_by, cwe.workout_instance_id, cwe.exercise_source, cwe.public_exercise_id, cwe.gym_exercise_id, 
		cwe.block_name, cwe.exercise_order, cwe.sets, cwe.reps_min, cwe.reps_max, cwe.weight_kg, 
		cwe.duration_seconds, cwe.rest_seconds, cwe.notes, cwe.created_at, cwe.updated_at 
		FROM "%s".custom_workout_exercise cwe
		LEFT JOIN public.exercise_equipment ee ON cwe.public_exercise_id = ee.exercise_id
		LEFT JOIN "%s".custom_exercise_equipment cee ON cwe.gym_exercise_id = cee.exercise_id
		WHERE ee.equipment_id = $1 OR cee.equipment_id = $1
		ORDER BY cwe.block_name, cwe.exercise_order`

	rows, err := r.DB.Query(fmt.Sprintf(query, gymID, gymID), equipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*dto.ResponseCustomWorkoutExerciseDTO
	for rows.Next() {
		var res dto.ResponseCustomWorkoutExerciseDTO
		var createdAt, updatedAt sql.NullString

		err := rows.Scan(
			&res.ID,
			&res.CreatedBy,
			&res.WorkoutInstanceID,
			&res.ExerciseSource,
			&res.PublicExerciseID,
			&res.GymExerciseID,
			&res.BlockName,
			&res.ExerciseOrder,
			&res.Sets,
			&res.RepsMin,
			&res.RepsMax,
			&res.WeightKg,
			&res.DurationSeconds,
			&res.RestSeconds,
			&res.Notes,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			res.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			res.UpdatedAt = updatedAt.String
		}

		result = append(result, &res)
	}

	return result, nil
}

func (r *CustomWorkoutExerciseRepository) Update(gymID string, exercise *dto.UpdateCustomWorkoutExerciseDTO) error {
	query := `UPDATE "%s".custom_workout_exercise SET 
		sets = $1, reps_min = $2, reps_max = $3, weight_kg = $4, 
		duration_seconds = $5, rest_seconds = $6, notes = $7, updated_at = NOW()
		WHERE id = $8`

	_, err := r.DB.Exec(
		fmt.Sprintf(query, gymID),
		exercise.Sets,
		exercise.RepsMin,
		exercise.RepsMax,
		exercise.WeightKg,
		exercise.DurationSeconds,
		exercise.RestSeconds,
		exercise.Notes,
		exercise.ID,
	)
	return err
}

func (r *CustomWorkoutExerciseRepository) Delete(gymID, id string) error {
	query := `DELETE FROM "%s".custom_workout_exercise WHERE id = $1`
	_, err := r.DB.Exec(fmt.Sprintf(query, gymID), id)
	return err
}
