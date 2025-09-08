package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"
)

type ExerciseMuscularGroupRepository struct {
	db *sql.DB
}

func NewExerciseMuscularGroupRepository(db *sql.DB) *ExerciseMuscularGroupRepository {
	return &ExerciseMuscularGroupRepository{db: db}
}

func (r *ExerciseMuscularGroupRepository) CreateLink(link *dto.ExerciseMuscularGroup) (*string, error) {
	query := `INSERT INTO exercise_muscular_group (exercise_id, muscular_group_id) VALUES ($1, $2) RETURNING exercise_id`
	var id string
	err := r.db.QueryRow(query, link.ExerciseID, link.MuscularGroupID).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *ExerciseMuscularGroupRepository) DeleteLink(id string) error {
	query := `DELETE FROM exercise_muscular_group WHERE exercise_id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ExerciseMuscularGroupRepository) RemoveAllLinksForExercise(exerciseID string) error {
	query := `DELETE FROM exercise_muscular_group WHERE exercise_id = $1`
	_, err := r.db.Exec(query, exerciseID)
	return err
}

func (r *ExerciseMuscularGroupRepository) FindByID(id string) (*dto.ExerciseMuscularGroup, error) {
	query := `SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE exercise_id = $1`
	var link dto.ExerciseMuscularGroup
	err := r.db.QueryRow(query, id).Scan(&link.ExerciseID, &link.MuscularGroupID)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *ExerciseMuscularGroupRepository) FindByExerciseID(exerciseID string) ([]*dto.ExerciseMuscularGroup, error) {
	query := `SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE exercise_id = $1`
	rows, err := r.db.Query(query, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var links []*dto.ExerciseMuscularGroup
	for rows.Next() {
		link := &dto.ExerciseMuscularGroup{}
		err := rows.Scan(&link.ExerciseID, &link.MuscularGroupID)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (r *ExerciseMuscularGroupRepository) FindByMuscularGroupID(muscularGroupID string) ([]*dto.ExerciseMuscularGroup, error) {
	query := `SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE muscular_group_id = $1`
	rows, err := r.db.Query(query, muscularGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var links []*dto.ExerciseMuscularGroup
	for rows.Next() {
		link := &dto.ExerciseMuscularGroup{}
		err := rows.Scan(&link.ExerciseID, &link.MuscularGroupID)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}
