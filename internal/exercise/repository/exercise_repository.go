package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/exercise/dto"
	"github.com/lib/pq"
)

type ExerciseRepository struct {
	db *sql.DB
}

func NewExerciseRepository(db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) CreateExercise(exercise dto.ExerciseCreationDTO) (string, error) {
	query := `
			   INSERT INTO public.exercise (
					   name, synonyms, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active
			   ) VALUES (
					   $1, $2, $3, $4, $5, $6, $7, $8, TRUE
			   ) RETURNING id`
	var id string
	err := r.db.QueryRow(query,
		exercise.Name,
		pq.Array(exercise.Synonyms), // store as TEXT[]
		exercise.DifficultyLevel,
		exercise.ExerciseType,
		exercise.Instructions,
		exercise.VideoURL,
		exercise.ImageURL,
		exercise.CreatedBy,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *ExerciseRepository) GetExerciseByID(id string) (dto.ExerciseResponseDTO, error) {
	var exercise dto.ExerciseResponseDTO
	query := `SELECT id, name, synonyms, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE id = $1 AND is_active = TRUE`
	err := r.db.QueryRow(query, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Synonyms,
		&exercise.DifficultyLevel,
		&exercise.ExerciseType,
		&exercise.Instructions,
		&exercise.VideoURL,
		&exercise.ImageURL,
		&exercise.CreatedBy,
		&exercise.IsActive,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)
	if err != nil {
		return exercise, err
	}
	return exercise, nil
}

func (r *ExerciseRepository) GetExerciseByName(name string) (dto.ExerciseResponseDTO, error) {
	var exercise dto.ExerciseResponseDTO
	query := `SELECT id, name, synonyms, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE name = $1 AND is_active = TRUE`
	err := r.db.QueryRow(query, name).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Synonyms,
		&exercise.DifficultyLevel,
		&exercise.ExerciseType,
		&exercise.Instructions,
		&exercise.VideoURL,
		&exercise.ImageURL,
		&exercise.CreatedBy,
		&exercise.IsActive,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)
	if err != nil {
		return exercise, err
	}
	return exercise, nil
}

func (r *ExerciseRepository) GetAllExercises() ([]dto.ExerciseResponseDTO, error) {
	query := `SELECT id, name, synonyms, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE is_active = TRUE`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var exercises []dto.ExerciseResponseDTO
	for rows.Next() {
		exercise := dto.ExerciseResponseDTO{}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Synonyms,
			&exercise.DifficultyLevel,
			&exercise.ExerciseType,
			&exercise.Instructions,
			&exercise.VideoURL,
			&exercise.ImageURL,
			&exercise.CreatedBy,
			&exercise.IsActive,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return exercises, nil
}

func (r *ExerciseRepository) UpdateExercise(id string, update dto.ExerciseUpdateDTO) (dto.ExerciseResponseDTO, error) {
	query := `
			   UPDATE public.exercise SET
					   name = COALESCE($2, name),
					   synonyms = COALESCE($3, synonyms),
					   difficulty_level = COALESCE($4, difficulty_level),
					   exercise_type = COALESCE($5, exercise_type),
					   instructions = COALESCE($6, instructions),
					   video_url = COALESCE($7, video_url),
					   image_url = COALESCE($8, image_url),
					   is_active = COALESCE($9, is_active),
					   updated_at = NOW()
			   WHERE id = $1
			   RETURNING id, name, synonyms, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at`
	var exercise dto.ExerciseResponseDTO
	err := r.db.QueryRow(query,
		id,
		update.Name,
		pq.Array(update.Synonyms),
		update.DifficultyLevel,
		update.ExerciseType,
		update.Instructions,
		update.VideoURL,
		update.ImageURL,
		update.IsActive,
	).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Synonyms,
		&exercise.DifficultyLevel,
		&exercise.ExerciseType,
		&exercise.Instructions,
		&exercise.VideoURL,
		&exercise.ImageURL,
		&exercise.CreatedBy,
		&exercise.IsActive,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)
	if err != nil {
		return exercise, err
	}
	return exercise, nil
}

func (r *ExerciseRepository) DeleteExercise(id string) error {
	query := `UPDATE public.exercise SET is_active = FALSE, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ExerciseRepository) GetExercisesByMuscularGroup(muscularGroups []string) ([]dto.ExerciseResponseDTO, error) {
	// Join with exercise_muscular_group
	query := `SELECT e.id, e.name, e.synonyms, e.difficulty_level, e.exercise_type, e.instructions, e.video_url, e.image_url, e.created_by, e.is_active, e.created_at, e.updated_at
				 FROM public.exercise e
				 JOIN public.exercise_muscular_group emg ON e.id = emg.exercise_id
				 WHERE e.is_active = TRUE AND emg.muscular_group_id = ANY($1)`
	rows, err := r.db.Query(query, pq.Array(muscularGroups))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var exercises []dto.ExerciseResponseDTO
	for rows.Next() {
		exercise := dto.ExerciseResponseDTO{}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Synonyms,
			&exercise.DifficultyLevel,
			&exercise.ExerciseType,
			&exercise.Instructions,
			&exercise.VideoURL,
			&exercise.ImageURL,
			&exercise.CreatedBy,
			&exercise.IsActive,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, nil
}

func (r *ExerciseRepository) GetExercisesByEquipment(equipment []string) ([]dto.ExerciseResponseDTO, error) {
	// Join with exercise_equipment
	query := `SELECT e.id, e.name, e.synonyms, e.difficulty_level, e.exercise_type, e.instructions, e.video_url, e.image_url, e.created_by, e.is_active, e.created_at, e.updated_at
				 FROM public.exercise e
				 JOIN public.exercise_equipment ee ON e.id = ee.exercise_id
				 WHERE e.is_active = TRUE AND ee.equipment_id = ANY($1)`
	rows, err := r.db.Query(query, pq.Array(equipment))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var exercises []dto.ExerciseResponseDTO
	for rows.Next() {
		exercise := dto.ExerciseResponseDTO{}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Synonyms,
			&exercise.DifficultyLevel,
			&exercise.ExerciseType,
			&exercise.Instructions,
			&exercise.VideoURL,
			&exercise.ImageURL,
			&exercise.CreatedBy,
			&exercise.IsActive,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, nil
}
