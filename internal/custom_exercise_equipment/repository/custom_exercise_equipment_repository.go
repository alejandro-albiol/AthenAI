package repository

import (
	"database/sql"
	"fmt"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/interfaces"
)

type CustomExerciseEquipmentRepositoryImpl struct {
	db *sql.DB
}

func NewCustomExerciseEquipmentRepository(db *sql.DB) interfaces.CustomExerciseEquipmentRepository {
	return &CustomExerciseEquipmentRepositoryImpl{db: db}
}

func (r *CustomExerciseEquipmentRepositoryImpl) CreateLink(gymID string, link *dto.CustomExerciseEquipment) (*string, error) {
	schema := gymID
	query := fmt.Sprintf("INSERT INTO %s.custom_exercise_equipment (custom_exercise_id, equipment_id) VALUES ($1, $2) RETURNING id", schema)
	var id string
	err := r.db.QueryRow(query, link.CustomExerciseID, link.EquipmentID).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *CustomExerciseEquipmentRepositoryImpl) DeleteLink(gymID, id string) error {
	schema := gymID
	query := fmt.Sprintf("DELETE FROM %s.custom_exercise_equipment WHERE id = $1", schema)
	result, err := r.db.Exec(query, id)
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

func (r *CustomExerciseEquipmentRepositoryImpl) RemoveAllLinksForExercise(gymID, customExerciseID string) error {
	schema := gymID
	query := fmt.Sprintf("DELETE FROM %s.custom_exercise_equipment WHERE custom_exercise_id = $1", schema)
	_, err := r.db.Exec(query, customExerciseID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CustomExerciseEquipmentRepositoryImpl) FindByID(gymID, id string) (*dto.CustomExerciseEquipment, error) {
	schema := gymID
	query := fmt.Sprintf("SELECT custom_exercise_id, equipment_id FROM %s.custom_exercise_equipment WHERE id = $1", schema)
	var link dto.CustomExerciseEquipment
	err := r.db.QueryRow(query, id).Scan(&link.CustomExerciseID, &link.EquipmentID)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *CustomExerciseEquipmentRepositoryImpl) FindByCustomExerciseID(gymID, customExerciseID string) ([]*dto.CustomExerciseEquipment, error) {
	schema := gymID
	query := fmt.Sprintf("SELECT id, equipment_id FROM %s.custom_exercise_equipment WHERE custom_exercise_id = $1", schema)
	rows, err := r.db.Query(query, customExerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*dto.CustomExerciseEquipment
	for rows.Next() {
		var link dto.CustomExerciseEquipment
		if err := rows.Scan(&link.ID, &link.EquipmentID); err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, nil
}

func (r *CustomExerciseEquipmentRepositoryImpl) FindByEquipmentID(gymID, equipmentID string) ([]*dto.CustomExerciseEquipment, error) {
	schema := gymID
	query := fmt.Sprintf("SELECT id, custom_exercise_id FROM %s.custom_exercise_equipment WHERE equipment_id = $1", schema)
	rows, err := r.db.Query(query, equipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*dto.CustomExerciseEquipment
	for rows.Next() {
		var link dto.CustomExerciseEquipment
		if err := rows.Scan(&link.ID, &link.CustomExerciseID); err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, nil
}
