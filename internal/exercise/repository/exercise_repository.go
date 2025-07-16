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
			name, synonyms, muscular_groups, equipment_needed, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, TRUE
		) RETURNING id`
	var id string
	err := r.db.QueryRow(query,
		exercise.Name,
		pq.Array(exercise.Synonyms),
		pq.Array(exercise.MuscularGroups),
		pq.Array(exercise.EquipmentNeeded),
		exercise.DifficultyLevel,
		exercise.ExerciseType,
		exercise.Instructions,
		exercise.VideoURL,
		exercise.ImageURL,
		exercise.CreatedBy,
	).Scan(&id)
	return id, err
}

func (r *ExerciseRepository) GetExerciseByID(id string) (dto.ExerciseResponseDTO, error) {
	var exercise dto.ExerciseResponseDTO
	query := `SELECT id, name, synonyms, muscular_groups, equipment_needed, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE id = $1 AND is_active = TRUE`
	err := r.db.QueryRow(query, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Synonyms,
		&exercise.MuscularGroups,
		&exercise.EquipmentNeeded,
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
	return exercise, err
}

func (r *ExerciseRepository) GetAllExercises() ([]dto.ExerciseResponseDTO, error) {
	query := `SELECT id, name, synonyms, muscular_groups, equipment_needed, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE is_active = TRUE`
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
			&exercise.MuscularGroups,
			&exercise.EquipmentNeeded,
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

func (r *ExerciseRepository) UpdateExercise(id string, update dto.ExerciseUpdateDTO) (dto.ExerciseResponseDTO, error) {
	query := `
		UPDATE public.exercise SET
			name = COALESCE($2, name),
			synonyms = COALESCE($3, synonyms),
			muscular_groups = COALESCE($4, muscular_groups),
			equipment_needed = COALESCE($5, equipment_needed),
			difficulty_level = COALESCE($6, difficulty_level),
			exercise_type = COALESCE($7, exercise_type),
			instructions = COALESCE($8, instructions),
			video_url = COALESCE($9, video_url),
			image_url = COALESCE($10, image_url),
			is_active = COALESCE($11, is_active),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, synonyms, muscular_groups, equipment_needed, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at`
	var exercise dto.ExerciseResponseDTO
	err := r.db.QueryRow(query,
		id,
		update.Name,
		pq.Array(update.Synonyms),
		pq.Array(update.MuscularGroups),
		pq.Array(update.EquipmentNeeded),
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
		&exercise.MuscularGroups,
		&exercise.EquipmentNeeded,
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
	return exercise, err
}

func (r *ExerciseRepository) DeleteExercise(id string) error {
	query := `UPDATE public.exercise SET is_active = FALSE, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ExerciseRepository) GetExercisesByMuscularGroup(muscularGroups []string) ([]dto.ExerciseResponseDTO, error) {
	query := `SELECT id, name, synonyms, muscular_groups, equipment_needed, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE is_active = TRUE AND muscular_groups && $1`
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
			&exercise.MuscularGroups,
			&exercise.EquipmentNeeded,
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
	query := `SELECT id, name, synonyms, muscular_groups, equipment_needed, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE is_active = TRUE AND equipment_needed && $1`
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
			&exercise.MuscularGroups,
			&exercise.EquipmentNeeded,
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
