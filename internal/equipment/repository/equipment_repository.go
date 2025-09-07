package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/equipment/dto"
)

type EquipmentRepository struct {
	db *sql.DB
}

func NewEquipmentRepository(db *sql.DB) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}
func (r *EquipmentRepository) CreateEquipment(equipment *dto.EquipmentCreationDTO) (*string, error) {
	query := `INSERT INTO public.equipment (name, description, category) VALUES ($1, $2, $3) RETURNING id`
	var id string
	err := r.db.QueryRow(query, equipment.Name, equipment.Description, equipment.Category).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
func (r *EquipmentRepository) GetEquipmentByID(id string) (*dto.EquipmentResponseDTO, error) {
	query := `SELECT id, name, description, category, is_active, created_at, updated_at FROM public.equipment WHERE id = $1`
	var equipment dto.EquipmentResponseDTO
	err := r.db.QueryRow(query, id).Scan(
		&equipment.ID,
		&equipment.Name,
		&equipment.Description,
		&equipment.Category,
		&equipment.IsActive,
		&equipment.CreatedAt,
		&equipment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &equipment, nil
}
func (r *EquipmentRepository) GetAllEquipment() ([]*dto.EquipmentResponseDTO, error) {
	query := `SELECT id, name, description, category, is_active, created_at, updated_at FROM public.equipment`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var equipmentList []*dto.EquipmentResponseDTO
	for rows.Next() {
		equipment := dto.EquipmentResponseDTO{}
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
		equipmentList = append(equipmentList, &equipment)
	}
	return equipmentList, nil
}
func (r *EquipmentRepository) UpdateEquipment(id string, update *dto.EquipmentUpdateDTO) (*dto.EquipmentResponseDTO, error) {
	query := `UPDATE public.equipment SET name = COALESCE($2, name), description = COALESCE($3, description), category = COALESCE($4, category), is_active = COALESCE($5, is_active), updated_at = NOW() WHERE id = $1 RETURNING id, name, description, category, is_active, created_at, updated_at`
	var equipment dto.EquipmentResponseDTO
	err := r.db.QueryRow(query, id, update.Name, update.Description, update.Category, update.IsActive).Scan(
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
	return &equipment, nil
}
func (r *EquipmentRepository) DeleteEquipment(id string) error {
	query := `DELETE FROM public.equipment WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
