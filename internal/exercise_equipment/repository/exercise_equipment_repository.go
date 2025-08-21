package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/interfaces"
)

type ExerciseEquipmentRepositoryImpl struct {
	db *sql.DB
}

func NewExerciseEquipmentRepository(db *sql.DB) interfaces.ExerciseEquipmentRepository {
	return &ExerciseEquipmentRepositoryImpl{db: db}
}

func (r *ExerciseEquipmentRepositoryImpl) CreateLink(link dto.ExerciseEquipment) (string, error) {
	query := `INSERT INTO exercise_equipment (exercise_id, equipment_id) VALUES ($1, $2) RETURNING id`
	var id string
	err := r.db.QueryRow(query, link.ExerciseID, link.EquipmentID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *ExerciseEquipmentRepositoryImpl) DeleteLink(id string) error {
	query := `DELETE FROM exercise_equipment WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ExerciseEquipmentRepositoryImpl) RemoveAllLinksForExercise(exerciseID string) error {
	query := `DELETE FROM exercise_equipment WHERE exercise_id = $1`
	_, err := r.db.Exec(query, exerciseID)
	return err
}

func (r *ExerciseEquipmentRepositoryImpl) FindByID(id string) (dto.ExerciseEquipment, error) {
	query := `SELECT exercise_id, equipment_id FROM exercise_equipment WHERE id = $1`
	var link dto.ExerciseEquipment
	err := r.db.QueryRow(query, id).Scan(&link.ExerciseID, &link.EquipmentID)
	if err != nil {
		return dto.ExerciseEquipment{}, err
	}
	return link, nil
}

func (r *ExerciseEquipmentRepositoryImpl) FindByExerciseID(exerciseID string) ([]dto.ExerciseEquipment, error) {
	query := `SELECT id, equipment_id FROM exercise_equipment WHERE exercise_id = $1`
	rows, err := r.db.Query(query, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []dto.ExerciseEquipment
	for rows.Next() {
		var link dto.ExerciseEquipment
		if err := rows.Scan(&link.ID, &link.EquipmentID); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (r *ExerciseEquipmentRepositoryImpl) FindByEquipmentID(equipmentID string) ([]dto.ExerciseEquipment, error) {
	query := `SELECT id, exercise_id FROM exercise_equipment WHERE equipment_id = $1`
	rows, err := r.db.Query(query, equipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []dto.ExerciseEquipment
	for rows.Next() {
		var link dto.ExerciseEquipment
		if err := rows.Scan(&link.ID, &link.ExerciseID); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}
