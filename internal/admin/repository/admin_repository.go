package repository

import (
	"database/sql"
	"fmt"

	"github.com/alejandro-albiol/athenai/internal/admin/dto"
	"github.com/alejandro-albiol/athenai/internal/admin/interfaces"
	"github.com/lib/pq"
)

type AdminRepository struct {
	sql *sql.DB
}

func NewAdminRepository(db *sql.DB) interfaces.AdminRepository {
	return &AdminRepository{sql: db}
}

// Exercise management methods
func (r *AdminRepository) CreateExercise(exercise dto.ExerciseCreationDTO) (string, error) {
	var exerciseID string
	query := `
		INSERT INTO public.exercise (
			name, synonyms, muscular_groups, equipment_needed, 
			difficulty_level, exercise_type, instructions, 
			video_url, image_url, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := r.sql.QueryRow(
		query,
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
	).Scan(&exerciseID)

	if err != nil {
		return "", err
	}

	return exerciseID, nil
}

func (r *AdminRepository) GetExerciseByID(id string) (dto.ExerciseResponseDTO, error) {
	var exercise dto.ExerciseResponseDTO
	query := `
		SELECT id, name, synonyms, muscular_groups, equipment_needed,
			   difficulty_level, exercise_type, instructions, video_url,
			   image_url, created_by, is_active, created_at, updated_at
		FROM public.exercise
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.sql.QueryRow(query, id).Scan(
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
		return dto.ExerciseResponseDTO{}, err
	}

	return exercise, nil
}

func (r *AdminRepository) GetAllExercises() ([]dto.ExerciseResponseDTO, error) {
	query := `
		SELECT id, name, synonyms, muscular_groups, equipment_needed,
			   difficulty_level, exercise_type, instructions, video_url,
			   image_url, created_by, is_active, created_at, updated_at
		FROM public.exercise
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.sql.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []dto.ExerciseResponseDTO
	for rows.Next() {
		var exercise dto.ExerciseResponseDTO
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

func (r *AdminRepository) UpdateExercise(id string, exercise dto.ExerciseUpdateDTO) (dto.ExerciseResponseDTO, error) {
	// Build dynamic query based on provided fields
	query := "UPDATE public.exercise SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 0

	if exercise.Name != nil {
		argCount++
		query += ", name = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *exercise.Name)
	}

	if exercise.Synonyms != nil {
		argCount++
		query += ", synonyms = $" + fmt.Sprintf("%d", argCount)
		args = append(args, pq.Array(exercise.Synonyms))
	}

	if exercise.MuscularGroups != nil {
		argCount++
		query += ", muscular_groups = $" + fmt.Sprintf("%d", argCount)
		args = append(args, pq.Array(exercise.MuscularGroups))
	}

	if exercise.EquipmentNeeded != nil {
		argCount++
		query += ", equipment_needed = $" + fmt.Sprintf("%d", argCount)
		args = append(args, pq.Array(exercise.EquipmentNeeded))
	}

	if exercise.DifficultyLevel != nil {
		argCount++
		query += ", difficulty_level = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *exercise.DifficultyLevel)
	}

	if exercise.ExerciseType != nil {
		argCount++
		query += ", exercise_type = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *exercise.ExerciseType)
	}

	if exercise.Instructions != nil {
		argCount++
		query += ", instructions = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *exercise.Instructions)
	}

	if exercise.VideoURL != nil {
		argCount++
		query += ", video_url = $" + fmt.Sprintf("%d", argCount)
		args = append(args, exercise.VideoURL)
	}

	if exercise.ImageURL != nil {
		argCount++
		query += ", image_url = $" + fmt.Sprintf("%d", argCount)
		args = append(args, exercise.ImageURL)
	}

	if exercise.IsActive != nil {
		argCount++
		query += ", is_active = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *exercise.IsActive)
	}

	argCount++
	query += " WHERE id = $" + fmt.Sprintf("%d", argCount) + " AND deleted_at IS NULL"
	args = append(args, id)

	_, err := r.sql.Exec(query, args...)
	if err != nil {
		return dto.ExerciseResponseDTO{}, err
	}

	// Return updated exercise
	return r.GetExerciseByID(id)
}

func (r *AdminRepository) DeleteExercise(id string) error {
	query := `
		UPDATE public.exercise 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.sql.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AdminRepository) GetExercisesByMuscularGroup(muscularGroups []string) ([]dto.ExerciseResponseDTO, error) {
	query := `
		SELECT id, name, synonyms, muscular_groups, equipment_needed,
			   difficulty_level, exercise_type, instructions, video_url,
			   image_url, created_by, is_active, created_at, updated_at
		FROM public.exercise
		WHERE muscular_groups && $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.sql.Query(query, pq.Array(muscularGroups))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []dto.ExerciseResponseDTO
	for rows.Next() {
		var exercise dto.ExerciseResponseDTO
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

func (r *AdminRepository) GetExercisesByEquipment(equipment []string) ([]dto.ExerciseResponseDTO, error) {
	query := `
		SELECT id, name, synonyms, muscular_groups, equipment_needed,
			   difficulty_level, exercise_type, instructions, video_url,
			   image_url, created_by, is_active, created_at, updated_at
		FROM public.exercise
		WHERE equipment_needed && $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.sql.Query(query, pq.Array(equipment))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []dto.ExerciseResponseDTO
	for rows.Next() {
		var exercise dto.ExerciseResponseDTO
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

// Equipment management methods
func (r *AdminRepository) CreateEquipment(equipment dto.EquipmentCreationDTO) (string, error) {
	var equipmentID string
	query := `
		INSERT INTO public.equipment (name, description, category)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.sql.QueryRow(
		query,
		equipment.Name,
		equipment.Description,
		equipment.Category,
	).Scan(&equipmentID)

	if err != nil {
		return "", err
	}

	return equipmentID, nil
}

func (r *AdminRepository) GetEquipmentByID(id string) (dto.EquipmentResponseDTO, error) {
	var equipment dto.EquipmentResponseDTO
	query := `
		SELECT id, name, description, category, is_active, created_at, updated_at
		FROM public.equipment
		WHERE id = $1
	`

	err := r.sql.QueryRow(query, id).Scan(
		&equipment.ID,
		&equipment.Name,
		&equipment.Description,
		&equipment.Category,
		&equipment.IsActive,
		&equipment.CreatedAt,
		&equipment.UpdatedAt,
	)

	if err != nil {
		return dto.EquipmentResponseDTO{}, err
	}

	return equipment, nil
}

func (r *AdminRepository) GetAllEquipment() ([]dto.EquipmentResponseDTO, error) {
	query := `
		SELECT id, name, description, category, is_active, created_at, updated_at
		FROM public.equipment
		ORDER BY created_at DESC
	`

	rows, err := r.sql.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipments []dto.EquipmentResponseDTO
	for rows.Next() {
		var equipment dto.EquipmentResponseDTO
		err := rows.Scan(
			&equipment.ID,
			&equipment.Name,
			&equipment.Description,
			&equipment.Category,
			&equipment.IsActive,
			&equipment.CreatedAt,
			&equipment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

func (r *AdminRepository) UpdateEquipment(id string, equipment dto.EquipmentUpdateDTO) (dto.EquipmentResponseDTO, error) {
	// Build dynamic query based on provided fields
	query := "UPDATE public.equipment SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 0

	if equipment.Name != nil {
		argCount++
		query += ", name = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *equipment.Name)
	}

	if equipment.Description != nil {
		argCount++
		query += ", description = $" + fmt.Sprintf("%d", argCount)
		args = append(args, equipment.Description)
	}

	if equipment.Category != nil {
		argCount++
		query += ", category = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *equipment.Category)
	}

	if equipment.IsActive != nil {
		argCount++
		query += ", is_active = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *equipment.IsActive)
	}

	argCount++
	query += " WHERE id = $" + fmt.Sprintf("%d", argCount)
	args = append(args, id)

	_, err := r.sql.Exec(query, args...)
	if err != nil {
		return dto.EquipmentResponseDTO{}, err
	}

	// Return updated equipment
	return r.GetEquipmentByID(id)
}

func (r *AdminRepository) DeleteEquipment(id string) error {
	query := `DELETE FROM public.equipment WHERE id = $1`

	result, err := r.sql.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Muscular group management methods
func (r *AdminRepository) CreateMuscularGroup(group dto.MuscularGroupCreationDTO) (string, error) {
	var groupID string
	query := `
		INSERT INTO public.muscular_group (name, description, body_part)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.sql.QueryRow(
		query,
		group.Name,
		group.Description,
		group.BodyPart,
	).Scan(&groupID)

	if err != nil {
		return "", err
	}

	return groupID, nil
}

func (r *AdminRepository) GetMuscularGroupByID(id string) (dto.MuscularGroupResponseDTO, error) {
	var group dto.MuscularGroupResponseDTO
	query := `
		SELECT id, name, description, body_part, is_active, created_at, updated_at
		FROM public.muscular_group
		WHERE id = $1
	`

	err := r.sql.QueryRow(query, id).Scan(
		&group.ID,
		&group.Name,
		&group.Description,
		&group.BodyPart,
		&group.IsActive,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if err != nil {
		return dto.MuscularGroupResponseDTO{}, err
	}

	return group, nil
}

func (r *AdminRepository) GetAllMuscularGroups() ([]dto.MuscularGroupResponseDTO, error) {
	query := `
		SELECT id, name, description, body_part, is_active, created_at, updated_at
		FROM public.muscular_group
		ORDER BY created_at DESC
	`

	rows, err := r.sql.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []dto.MuscularGroupResponseDTO
	for rows.Next() {
		var group dto.MuscularGroupResponseDTO
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.Description,
			&group.BodyPart,
			&group.IsActive,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (r *AdminRepository) UpdateMuscularGroup(id string, group dto.MuscularGroupUpdateDTO) (dto.MuscularGroupResponseDTO, error) {
	// Build dynamic query based on provided fields
	query := "UPDATE public.muscular_group SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 0

	if group.Name != nil {
		argCount++
		query += ", name = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *group.Name)
	}

	if group.Description != nil {
		argCount++
		query += ", description = $" + fmt.Sprintf("%d", argCount)
		args = append(args, group.Description)
	}

	if group.BodyPart != nil {
		argCount++
		query += ", body_part = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *group.BodyPart)
	}

	if group.IsActive != nil {
		argCount++
		query += ", is_active = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *group.IsActive)
	}

	argCount++
	query += " WHERE id = $" + fmt.Sprintf("%d", argCount)
	args = append(args, id)

	_, err := r.sql.Exec(query, args...)
	if err != nil {
		return dto.MuscularGroupResponseDTO{}, err
	}

	// Return updated group
	return r.GetMuscularGroupByID(id)
}

func (r *AdminRepository) DeleteMuscularGroup(id string) error {
	query := `DELETE FROM public.muscular_group WHERE id = $1`

	result, err := r.sql.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
