package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"
)

type CustomWorkoutInstanceRepositoryImpl struct {
	DB *sql.DB
}

func NewCustomWorkoutInstanceRepository(db *sql.DB) *CustomWorkoutInstanceRepositoryImpl {
	return &CustomWorkoutInstanceRepositoryImpl{DB: db}
}

func (r *CustomWorkoutInstanceRepositoryImpl) Create(gymID string, instance dto.CreateCustomWorkoutInstanceDTO) (string, error) {
	query := `INSERT INTO "` + gymID + `".custom_workout_instance (
		name, description, template_source, public_template_id, gym_template_id, difficulty_level, estimated_duration_minutes, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING id`
	var id string
	err := r.DB.QueryRow(query,
		instance.Name,
		instance.Description,
		instance.TemplateSource,
		instance.PublicTemplateID,
		instance.GymTemplateID,
		instance.DifficultyLevel,
		instance.EstimatedDurationMinutes,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *CustomWorkoutInstanceRepositoryImpl) GetByID(gymID, id string) (dto.ResponseCustomWorkoutInstanceDTO, error) {
	query := `SELECT id, name, description, template_source, public_template_id, gym_template_id, difficulty_level, estimated_duration_minutes, created_at, updated_at
		FROM "` + gymID + `".custom_workout_instance WHERE id = $1 AND deleted_at IS NULL`
	var instance dto.ResponseCustomWorkoutInstanceDTO
	err := r.DB.QueryRow(query, id).Scan(
		&instance.ID,
		&instance.Name,
		&instance.Description,
		&instance.TemplateSource,
		&instance.PublicTemplateID,
		&instance.GymTemplateID,
		&instance.DifficultyLevel,
		&instance.EstimatedDurationMinutes,
		&instance.CreatedAt,
		&instance.UpdatedAt,
	)
	if err != nil {
		return dto.ResponseCustomWorkoutInstanceDTO{}, err
	}
	return instance, nil
}

func (r *CustomWorkoutInstanceRepositoryImpl) List(gymID string) ([]dto.ResponseCustomWorkoutInstanceDTO, error) {
	query := `SELECT id, name, description, template_source, public_template_id, gym_template_id, difficulty_level, estimated_duration_minutes, created_at, updated_at
		FROM "` + gymID + `".custom_workout_instance WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var instances []dto.ResponseCustomWorkoutInstanceDTO
	for rows.Next() {
		instance := dto.ResponseCustomWorkoutInstanceDTO{}
		err := rows.Scan(
			&instance.ID,
			&instance.Name,
			&instance.Description,
			&instance.TemplateSource,
			&instance.PublicTemplateID,
			&instance.GymTemplateID,
			&instance.DifficultyLevel,
			&instance.EstimatedDurationMinutes,
			&instance.CreatedAt,
			&instance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *CustomWorkoutInstanceRepositoryImpl) Update(gymID string, instance dto.UpdateCustomWorkoutInstanceDTO) error {
	query := `UPDATE "` + gymID + `".custom_workout_instance SET
		name = COALESCE($1, name),
		description = COALESCE($2, description),
		template_source = COALESCE($3, template_source),
		public_template_id = COALESCE($4, public_template_id),
		gym_template_id = COALESCE($5, gym_template_id),
		difficulty_level = COALESCE($6, difficulty_level),
		estimated_duration_minutes = COALESCE($7, estimated_duration_minutes),
		updated_at = NOW()
		WHERE id = $8 AND deleted_at IS NULL`
	result, err := r.DB.Exec(query,
		instance.Name,
		instance.Description,
		instance.TemplateSource,
		instance.PublicTemplateID,
		instance.GymTemplateID,
		instance.DifficultyLevel,
		instance.EstimatedDurationMinutes,
		instance.ID,
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
	query := `UPDATE "` + gymID + `".custom_workout_instance SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
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
