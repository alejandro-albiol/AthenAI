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

func (r *CustomExerciseRepository) Create(gymID string, exercise *dto.CustomExerciseCreationDTO) error {
	query := `INSERT INTO "%s".custom_exercise (created_by, name, synonyms, difficulty_level, exercise_type, instructions, video_url, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	schema := gymID
	_, err := r.DB.Exec(
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
	return err
}

func (r *CustomExerciseRepository) GetByID(gymID, id string) (*dto.CustomExerciseResponseDTO, error) {
	query := `SELECT id, created_by, name, synonyms, difficulty_level, exercise_type, instructions, video_url, image_url, is_active FROM "%s".custom_exercise WHERE id = $1`
	schema := gymID
	row := r.DB.QueryRow(fmt.Sprintf(query, schema), id)
	var res dto.CustomExerciseResponseDTO
	if err := row.Scan(&res.ID, &res.CreatedBy, &res.Name, &res.Synonyms, &res.DifficultyLevel, &res.ExerciseType, &res.Instructions, &res.VideoURL, &res.ImageURL, &res.IsActive); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *CustomExerciseRepository) List(gymID string) ([]*dto.CustomExerciseResponseDTO, error) {
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

func (r *CustomExerciseRepository) Update(gymID string, exercise *dto.CustomExerciseUpdateDTO) error {
	query := `UPDATE "%s".custom_exercise SET name = $1, synonyms = $2, difficulty_level = $3, exercise_type = $4, instructions = $5, video_url = $6, image_url = $7, is_active = $8 WHERE id = $9`
	schema := gymID
	_, err := r.DB.Exec(
		fmt.Sprintf(query, schema),
		exercise.Name,
		exercise.Synonyms,
		exercise.DifficultyLevel,
		exercise.ExerciseType,
		exercise.Instructions,
		exercise.VideoURL,
		exercise.ImageURL,
		exercise.IsActive,
		exercise.ID,
	)
	return err
}

func (r *CustomExerciseRepository) Delete(gymID, id string) error {
	query := `DELETE FROM "%s".custom_exercise WHERE id = $1`
	schema := gymID
	_, err := r.DB.Exec(fmt.Sprintf(query, schema), id)
	return err
}
