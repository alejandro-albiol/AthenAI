package repository

import (
	"database/sql"
	"fmt"

	"github.com/alejandro-albiol/athenai/internal/custom_equipment/dto"
)

// CustomEquipmentRepositoryImpl implements interfaces.CustomEquipmentRepository

type CustomEquipmentRepository struct {
	DB *sql.DB
}

func NewCustomEquipmentRepository(db *sql.DB) *CustomEquipmentRepository {
	return &CustomEquipmentRepository{DB: db}
}

func (r *CustomEquipmentRepository) Create(gymID string, equipment *dto.CreateCustomEquipmentDTO) (*string, error) {
	query := `INSERT INTO "%s".custom_equipment (created_by, name, description, category, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string
	err := r.DB.QueryRow(
		fmt.Sprintf(query, gymID),
		equipment.CreatedBy,
		equipment.Name,
		equipment.Description,
		equipment.Category,
		equipment.IsActive,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *CustomEquipmentRepository) GetByID(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error) {
	query := `SELECT id, created_by, name, description, category, is_active FROM "%s".custom_equipment WHERE id = $1`
	schema := gymID
	row := r.DB.QueryRow(fmt.Sprintf(query, schema), id)
	var res dto.ResponseCustomEquipmentDTO
	if err := row.Scan(&res.ID, &res.CreatedBy, &res.Name, &res.Description, &res.Category, &res.IsActive); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *CustomEquipmentRepository) GetByName(gymID, name string) (*dto.ResponseCustomEquipmentDTO, error) {
	query := `SELECT id, created_by, name, description, category, is_active FROM "%s".custom_equipment WHERE name = $1`
	schema := gymID
	row := r.DB.QueryRow(fmt.Sprintf(query, schema), name)
	var res dto.ResponseCustomEquipmentDTO
	if err := row.Scan(&res.ID, &res.CreatedBy, &res.Name, &res.Description, &res.Category, &res.IsActive); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *CustomEquipmentRepository) List(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error) {
	query := `SELECT id, created_by, name, description, category, is_active FROM "%s".custom_equipment`
	schema := gymID
	rows, err := r.DB.Query(fmt.Sprintf(query, schema))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*dto.ResponseCustomEquipmentDTO
	for rows.Next() {
		var res dto.ResponseCustomEquipmentDTO
		if err := rows.Scan(&res.ID, &res.CreatedBy, &res.Name, &res.Description, &res.Category, &res.IsActive); err != nil {
			return nil, err
		}
		result = append(result, &res)
	}
	return result, nil
}

func (r *CustomEquipmentRepository) Update(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error {
	query := `UPDATE "%s".custom_equipment SET name = $1, description = $2, category = $3, is_active = $4 WHERE id = $5`
	schema := gymID
	_, err := r.DB.Exec(
		fmt.Sprintf(query, schema),
		equipment.Name,
		equipment.Description,
		equipment.Category,
		equipment.IsActive,
		equipment.ID,
	)
	return err
}

func (r *CustomEquipmentRepository) Delete(gymID, id string) error {
	query := `DELETE FROM "%s".custom_equipment WHERE id = $1`
	schema := gymID
	_, err := r.DB.Exec(fmt.Sprintf(query, schema), id)
	return err
}
