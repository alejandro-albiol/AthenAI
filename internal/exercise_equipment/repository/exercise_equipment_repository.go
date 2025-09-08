package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"
)

type ExerciseEquipmentRepository struct {
	db *sql.DB
}

func NewExerciseEquipmentRepository(db *sql.DB) *ExerciseEquipmentRepository {
	return &ExerciseEquipmentRepository{db: db}
}

func (r *ExerciseEquipmentRepository) CreateLink(link *dto.ExerciseEquipment) (*string, error) {
	query := `INSERT INTO exercise_equipment (exercise_id, equipment_id) VALUES ($1, $2) RETURNING id`
	var id string
	err := r.db.QueryRow(query, link.ExerciseID, link.EquipmentID).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *ExerciseEquipmentRepository) DeleteLink(id string) error {
	query := `DELETE FROM exercise_equipment WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ExerciseEquipmentRepository) RemoveAllLinksForExercise(exerciseID string) error {
	query := `DELETE FROM exercise_equipment WHERE exercise_id = $1`
	_, err := r.db.Exec(query, exerciseID)
	return err
}

func (r *ExerciseEquipmentRepository) FindByID(id string) (*dto.ExerciseEquipment, error) {
	query := `SELECT exercise_id, equipment_id FROM exercise_equipment WHERE id = $1`
	link := &dto.ExerciseEquipment{}
	err := r.db.QueryRow(query, id).Scan(&link.ExerciseID, &link.EquipmentID)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (r *ExerciseEquipmentRepository) FindByExerciseID(exerciseID string) ([]*dto.ExerciseEquipment, error) {
	query := `SELECT id, equipment_id FROM exercise_equipment WHERE exercise_id = $1`
	rows, err := r.db.Query(query, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*dto.ExerciseEquipment
	for rows.Next() {
		link := &dto.ExerciseEquipment{}
		if err := rows.Scan(&link.ID, &link.EquipmentID); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (r *ExerciseEquipmentRepository) FindByEquipmentID(equipmentID string) ([]*dto.ExerciseEquipment, error) {
	query := `SELECT id, exercise_id FROM exercise_equipment WHERE equipment_id = $1`
	rows, err := r.db.Query(query, equipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*dto.ExerciseEquipment
	for rows.Next() {
		link := &dto.ExerciseEquipment{}
		if err := rows.Scan(&link.ID, &link.ExerciseID); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}
