package repository

import (
	"database/sql"
	"fmt"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/interfaces"
)

type CustomExerciseMuscularGroupRepositoryImpl struct {
	db *sql.DB
}

func NewCustomExerciseMuscularGroupRepository(db *sql.DB) interfaces.CustomExerciseMuscularGroupRepository {
	return &CustomExerciseMuscularGroupRepositoryImpl{db: db}
}

func (r *CustomExerciseMuscularGroupRepositoryImpl) CreateLink(gymID string, link dto.CustomExerciseMuscularGroup) (string, error) {
	schema := gymID
	query := fmt.Sprintf("INSERT INTO %s.custom_exercise_muscular_group (custom_exercise_id, muscular_group_id) VALUES ($1, $2) RETURNING id", schema)
	var id string
	err := r.db.QueryRow(query, link.CustomExerciseID, link.MuscularGroupID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *CustomExerciseMuscularGroupRepositoryImpl) DeleteLink(gymID, id string) error {
	schema := gymID
	query := fmt.Sprintf("DELETE FROM %s.custom_exercise_muscular_group WHERE id = $1", schema)
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

func (r *CustomExerciseMuscularGroupRepositoryImpl) RemoveAllLinksForExercise(gymID, customExerciseID string) error {
	schema := gymID
	query := fmt.Sprintf("DELETE FROM %s.custom_exercise_muscular_group WHERE custom_exercise_id = $1", schema)
	_, err := r.db.Exec(query, customExerciseID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CustomExerciseMuscularGroupRepositoryImpl) FindByID(gymID, id string) (dto.CustomExerciseMuscularGroup, error) {
	schema := gymID
	query := fmt.Sprintf("SELECT custom_exercise_id, muscular_group_id FROM %s.custom_exercise_muscular_group WHERE id = $1", schema)
	var link dto.CustomExerciseMuscularGroup
	err := r.db.QueryRow(query, id).Scan(&link.CustomExerciseID, &link.MuscularGroupID)
	if err != nil {
		return dto.CustomExerciseMuscularGroup{}, err
	}
	return link, nil
}

func (r *CustomExerciseMuscularGroupRepositoryImpl) FindByCustomExerciseID(gymID, customExerciseID string) ([]dto.CustomExerciseMuscularGroup, error) {
	schema := gymID
	query := fmt.Sprintf("SELECT id, muscular_group_id FROM %s.custom_exercise_muscular_group WHERE custom_exercise_id = $1", schema)
	rows, err := r.db.Query(query, customExerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []dto.CustomExerciseMuscularGroup
	for rows.Next() {
		var link dto.CustomExerciseMuscularGroup
		if err := rows.Scan(&link.ID, &link.MuscularGroupID); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (r *CustomExerciseMuscularGroupRepositoryImpl) FindByMuscularGroupID(gymID, muscularGroupID string) ([]dto.CustomExerciseMuscularGroup, error) {
	schema := gymID
	query := fmt.Sprintf("SELECT id, custom_exercise_id FROM %s.custom_exercise_muscular_group WHERE muscular_group_id = $1", schema)
	rows, err := r.db.Query(query, muscularGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []dto.CustomExerciseMuscularGroup
	for rows.Next() {
		var link dto.CustomExerciseMuscularGroup
		if err := rows.Scan(&link.ID, &link.CustomExerciseID); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}
