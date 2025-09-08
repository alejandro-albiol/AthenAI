package repository

import (
	"database/sql"
	"fmt"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"
)

type CustomExerciseRepository struct {
	DB *sql.DB
}

func NewCustomExerciseRepository(db *sql.DB) *CustomExerciseRepository {
	return &CustomExerciseRepository{DB: db}
}

func (r *CustomExerciseRepository) CreateCustomExercise(gymID string, exercise *dto.CustomExerciseCreationDTO) (*string, error) {
	query := `INSERT INTO "%s".custom_exercise (created_by, name, synonyms, difficulty_level, exercise_type, instructions, video_url, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	schema := gymID
	row := r.DB.QueryRow(
		fmt.Sprintf(query, schema),
		exercise.CreatedBy,
		exercise.Name,
		exercise.Synonyms,
		exercise.DifficultyLevel,
		exercise.ExerciseType,
		exercise.Instructions,
		exercise.VideoURL,
		exercise.ImageURL,
	)
	var id string
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *CustomExerciseRepository) UpdateCustomExercise(gymID, id string, update *dto.CustomExerciseUpdateDTO) error {
	query := `UPDATE "%s".custom_exercise SET name = $1, synonyms = $2, difficulty_level = $3, exercise_type = $4, instructions = $5, video_url = $6, image_url = $7, is_active = $8 WHERE id = $9`
	schema := gymID
	_, err := r.DB.Exec(
		fmt.Sprintf(query, schema),
		update.Name,
		update.Synonyms,
		update.DifficultyLevel,
		update.ExerciseType,
		update.Instructions,
		update.VideoURL,
		update.ImageURL,
		update.IsActive,
		id,
	)
	return err
}

func (r *CustomExerciseRepository) GetCustomExerciseByID(gymID, id string) (*dto.CustomExerciseResponseDTO, error) {
	query := `SELECT id, created_by, name, synonyms, difficulty_level, exercise_type, instructions, video_url, image_url, is_active FROM "%s".custom_exercise WHERE id = $1`
	schema := gymID
	row := r.DB.QueryRow(fmt.Sprintf(query, schema), id)
	var res dto.CustomExerciseResponseDTO
	err := row.Scan(&res.ID, &res.CreatedBy, &res.Name, &res.Synonyms, &res.DifficultyLevel, &res.ExerciseType, &res.Instructions, &res.VideoURL, &res.ImageURL, &res.IsActive)
	return &res, err
}

func (r *CustomExerciseRepository) ListCustomExercises(gymID string) ([]*dto.CustomExerciseResponseDTO, error) {
	query := `SELECT id, created_by, name, synonyms, difficulty_level, exercise_type, instructions, video_url, image_url, is_active FROM "%s".custom_exercise`
	schema := gymID
	rows, err := r.DB.Query(fmt.Sprintf(query, schema))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*dto.CustomExerciseResponseDTO
	for rows.Next() {
		var res dto.CustomExerciseResponseDTO
		if err := rows.Scan(&res.ID, &res.CreatedBy, &res.Name, &res.Synonyms, &res.DifficultyLevel, &res.ExerciseType, &res.Instructions, &res.VideoURL, &res.ImageURL, &res.IsActive); err != nil {
			return nil, err
		}
		result = append(result, &res)
	}
	return result, nil
}

func (r *CustomExerciseRepository) DeleteCustomExercise(gymID, id string) error {
	// Soft delete: set deleted_at
	query := `UPDATE "%s".custom_exercise SET deleted_at = NOW() WHERE id = $1`
	schema := gymID
	_, err := r.DB.Exec(fmt.Sprintf(query, schema), id)
	return err
}
