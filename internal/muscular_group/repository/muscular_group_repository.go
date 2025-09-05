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

func (r *MuscularGroupRepository) CreateMuscularGroup(mg *dto.CreateMuscularGroupDTO) (*string, error) {
	query := `INSERT INTO public.muscular_group (name, description, body_part, is_active) VALUES ($1, $2, $3, TRUE) RETURNING id`
	var id string
	err := r.db.QueryRow(query, mg.Name, mg.Description, mg.BodyPart).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *MuscularGroupRepository) GetAllMuscularGroups() ([]*dto.MuscularGroupResponseDTO, error) {
	query := `SELECT id, name, description, body_part, is_active FROM public.muscular_group WHERE is_active = TRUE`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []*dto.MuscularGroupResponseDTO
	for rows.Next() {
		mg := &dto.MuscularGroupResponseDTO{}
		err := rows.Scan(&mg.ID, &mg.Name, &mg.Description, &mg.BodyPart, &mg.IsActive)
		if err != nil {
			return nil, err
		}
		groups = append(groups, mg)
	}
	return groups, nil
}

func (r *MuscularGroupRepository) GetMuscularGroupByID(id string) (*dto.MuscularGroupResponseDTO, error) {
	mg := &dto.MuscularGroupResponseDTO{}
	query := `SELECT id, name, description, body_part, is_active FROM public.muscular_group WHERE id = $1 AND is_active = TRUE`
	err := r.db.QueryRow(query, id).Scan(&mg.ID, &mg.Name, &mg.Description, &mg.BodyPart, &mg.IsActive)
	if err != nil {
		return nil, err
	}
	return mg, nil
}

func (r *MuscularGroupRepository) GetMuscularGroupByName(name string) (*dto.MuscularGroupResponseDTO, error) {
	mg := &dto.MuscularGroupResponseDTO{}
	query := `SELECT id, name, description, body_part, is_active FROM public.muscular_group WHERE name = $1 AND is_active = TRUE`
	err := r.db.QueryRow(query, name).Scan(&mg.ID, &mg.Name, &mg.Description, &mg.BodyPart, &mg.IsActive)
	if err != nil {
		return nil, err
	}
	return mg, nil
}

func (r *MuscularGroupRepository) UpdateMuscularGroup(id string, mg *dto.UpdateMuscularGroupDTO) (*dto.MuscularGroupResponseDTO, error) {
	query := `UPDATE public.muscular_group SET name = COALESCE($2, name), description = COALESCE($3, description), body_part = COALESCE($4, body_part) WHERE id = $1 RETURNING id, name, description, body_part, is_active`
	updated := &dto.MuscularGroupResponseDTO{}
	err := r.db.QueryRow(query, id, mg.Name, mg.Description, mg.BodyPart).Scan(&updated.ID, &updated.Name, &updated.Description, &updated.BodyPart, &updated.IsActive)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *MuscularGroupRepository) DeleteMuscularGroup(id string) error {
	query := `UPDATE public.muscular_group SET is_active = FALSE WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
