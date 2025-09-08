package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/dto"
)

type CustomMemberWorkoutRepository struct {
	DB *sql.DB
}

func NewCustomMemberWorkoutRepository(db *sql.DB) *CustomMemberWorkoutRepository {
	return &CustomMemberWorkoutRepository{DB: db}
}

func (r *CustomMemberWorkoutRepository) Create(gymID string, memberWorkout *dto.CreateCustomMemberWorkoutDTO) (*string, error) {
	query := `INSERT INTO "` + gymID + `".custom_member_workout (
		created_by, member_id, workout_instance_id, scheduled_date, notes, rating, status
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`

	// For now, set created_by to member_id (or pass as param if available)
	var id string
	err := r.DB.QueryRow(
		query,
		memberWorkout.MemberID, // created_by
		memberWorkout.MemberID,
		memberWorkout.WorkoutInstanceID,
		memberWorkout.ScheduledDate,
		memberWorkout.Notes,
		memberWorkout.Rating,
		"scheduled", // default status
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *CustomMemberWorkoutRepository) GetByID(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error) {
	query := `SELECT id, member_id, workout_instance_id, scheduled_date, started_at, completed_at, status, notes, rating, created_at, updated_at
		FROM "` + gymID + `".custom_member_workout WHERE id = $1`
	row := r.DB.QueryRow(query, id)
	var res dto.ResponseCustomMemberWorkoutDTO
	err := row.Scan(
		&res.ID,
		&res.MemberID,
		&res.WorkoutInstanceID,
		&res.ScheduledDate,
		&res.StartedAt,
		&res.CompletedAt,
		&res.Status,
		&res.Notes,
		&res.Rating,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *CustomMemberWorkoutRepository) ListByMemberID(gymID, memberID string) ([]*dto.ResponseCustomMemberWorkoutDTO, error) {
	query := `SELECT id, member_id, workout_instance_id, scheduled_date, started_at, completed_at, status, notes, rating, created_at, updated_at
		FROM "` + gymID + `".custom_member_workout WHERE member_id = $1 ORDER BY scheduled_date DESC`
	rows, err := r.DB.Query(query, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*dto.ResponseCustomMemberWorkoutDTO
	for rows.Next() {
		var res dto.ResponseCustomMemberWorkoutDTO
		err := rows.Scan(
			&res.ID,
			&res.MemberID,
			&res.WorkoutInstanceID,
			&res.ScheduledDate,
			&res.StartedAt,
			&res.CompletedAt,
			&res.Status,
			&res.Notes,
			&res.Rating,
			&res.CreatedAt,
			&res.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &res)
	}
	return result, nil
}

func (r *CustomMemberWorkoutRepository) Update(gymID string, memberWorkout *dto.UpdateCustomMemberWorkoutDTO) error {
	query := `UPDATE "` + gymID + `".custom_member_workout SET
	       started_at = COALESCE($1, started_at),
	       completed_at = COALESCE($2, completed_at),
	       status = COALESCE($3, status),
	       notes = COALESCE($4, notes),
	       rating = COALESCE($5, rating),
	       updated_at = NOW()
       WHERE id = $6`
	res, err := r.DB.Exec(
		query,
		memberWorkout.StartedAt,
		memberWorkout.CompletedAt,
		memberWorkout.Status,
		memberWorkout.Notes,
		memberWorkout.Rating,
		memberWorkout.ID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CustomMemberWorkoutRepository) Delete(gymID, id string) error {
	query := `DELETE FROM "` + gymID + `".custom_member_workout WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
