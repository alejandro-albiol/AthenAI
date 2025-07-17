package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/muscular_group/dto"
)

type MuscularGroupRepository struct {
	db *sql.DB
}

func NewMuscularGroupRepository(db *sql.DB) *MuscularGroupRepository {
	return &MuscularGroupRepository{db: db}
}

func (r *MuscularGroupRepository) CreateMuscularGroup(mg dto.MuscularGroup) (string, error) {
	query := `INSERT INTO public.muscular_group (name, description, body_part, is_active) VALUES ($1, $2, $3, TRUE) RETURNING id`
	var id string
	err := r.db.QueryRow(query, mg.Name, mg.Description, mg.BodyPart).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *MuscularGroupRepository) GetAllMuscularGroups() ([]dto.MuscularGroup, error) {
	query := `SELECT id, name, description, body_part, is_active FROM public.muscular_group WHERE is_active = TRUE`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []dto.MuscularGroup
	for rows.Next() {
		var mg dto.MuscularGroup
		err := rows.Scan(&mg.ID, &mg.Name, &mg.Description, &mg.BodyPart, &mg.IsActive)
		if err != nil {
			return nil, err
		}
		groups = append(groups, mg)
	}
	return groups, nil
}

func (r *MuscularGroupRepository) GetMuscularGroupByID(id string) (dto.MuscularGroup, error) {
	var mg dto.MuscularGroup
	query := `SELECT id, name, description, body_part, is_active FROM public.muscular_group WHERE id = $1 AND is_active = TRUE`
	err := r.db.QueryRow(query, id).Scan(&mg.ID, &mg.Name, &mg.Description, &mg.BodyPart, &mg.IsActive)
	if err != nil {
		return dto.MuscularGroup{}, err
	}
	return mg, nil
}

func (r *MuscularGroupRepository) GetMuscularGroupByName(name string) (*dto.MuscularGroup, error) {
	var mg dto.MuscularGroup
	query := `SELECT id, name, description, body_part, is_active FROM public.muscular_group WHERE name = $1 AND is_active = TRUE`
	err := r.db.QueryRow(query, name).Scan(&mg.ID, &mg.Name, &mg.Description, &mg.BodyPart, &mg.IsActive)
	if err != nil {
		return nil, err
	}
	return &mg, nil
}

func (r *MuscularGroupRepository) UpdateMuscularGroup(id string, mg dto.MuscularGroup) (dto.MuscularGroup, error) {
	query := `UPDATE public.muscular_group SET name = COALESCE($2, name), description = COALESCE($3, description), body_part = COALESCE($4, body_part) WHERE id = $1 RETURNING id, name, description, body_part, is_active`
	var updated dto.MuscularGroup
	err := r.db.QueryRow(query, id, mg.Name, mg.Description, mg.BodyPart).Scan(&updated.ID, &updated.Name, &updated.Description, &updated.BodyPart, &updated.IsActive)
	if err != nil {
		return dto.MuscularGroup{}, err
	}
	return updated, nil
}

func (r *MuscularGroupRepository) DeleteMuscularGroup(id string) error {
	query := `UPDATE public.muscular_group SET is_active = FALSE WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *MuscularGroupRepository) FindByID(id string) (dto.MuscularGroup, error) {
	return r.GetMuscularGroupByID(id)
}
