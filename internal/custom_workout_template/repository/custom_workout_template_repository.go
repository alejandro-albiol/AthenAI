package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/dto"
)

type CustomWorkoutTemplateRepositoryImpl struct {
	DB *sql.DB
}

func NewCustomWorkoutTemplateRepository(db *sql.DB) *CustomWorkoutTemplateRepositoryImpl {
	return &CustomWorkoutTemplateRepositoryImpl{DB: db}
}

func (r *CustomWorkoutTemplateRepositoryImpl) Create(gymID string, template dto.CreateCustomWorkoutTemplateDTO) (string, error) {
	query := `INSERT INTO "` + gymID + `".custom_workout_template (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING id`

	var id string
	err := r.DB.QueryRow(
		query,
		template.Name,
		template.Description,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *CustomWorkoutTemplateRepositoryImpl) GetByID(gymID, id string) (dto.ResponseCustomWorkoutTemplateDTO, error) {
	query := `SELECT id, name, description, created_at, updated_at
		FROM "` + gymID + `".custom_workout_template WHERE id = $1 AND deleted_at IS NULL`
	var template dto.ResponseCustomWorkoutTemplateDTO
	err := r.DB.QueryRow(query, id).Scan(
		&template.ID,
		&template.Name,
		&template.Description,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		return dto.ResponseCustomWorkoutTemplateDTO{}, err
	}
	return template, nil
}

func (r *CustomWorkoutTemplateRepositoryImpl) GetByName(gymID, name string) (dto.ResponseCustomWorkoutTemplateDTO, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM "` + gymID + `".custom_workout_template WHERE name = $1 AND deleted_at IS NULL`
	var template dto.ResponseCustomWorkoutTemplateDTO
	err := r.DB.QueryRow(query, name).Scan(
		&template.ID,
		&template.Name,
		&template.Description,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		return dto.ResponseCustomWorkoutTemplateDTO{}, err
	}
	return template, nil
}

func (r *CustomWorkoutTemplateRepositoryImpl) List(gymID string) ([]dto.ResponseCustomWorkoutTemplateDTO, error) {
	query := `SELECT id, name, description, created_at, updated_at
		FROM "` + gymID + `".custom_workout_template WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := []dto.ResponseCustomWorkoutTemplateDTO{}
	for rows.Next() {
		template := dto.ResponseCustomWorkoutTemplateDTO{}
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Description,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *CustomWorkoutTemplateRepositoryImpl) Update(gymID string, template dto.UpdateCustomWorkoutTemplateDTO) error {
	query := `UPDATE "` + gymID + `".custom_workout_template
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL`
	result, err := r.DB.Exec(query,
		template.Name,
		template.Description,
		template.ID,
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

func (r *CustomWorkoutTemplateRepositoryImpl) Delete(gymID, id string) error {
	query := `UPDATE "` + gymID + `".custom_workout_template SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
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
