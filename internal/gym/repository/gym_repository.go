package repository

import (
	"database/sql"
	"time"

	"github.com/alejandro-albiol/athenai/internal/gym/dto"
)

type GymRepository struct {
	db *sql.DB
}

func NewGymRepository(db *sql.DB) *GymRepository {
	return &GymRepository{db: db}
}

func (r *GymRepository) CreateGym(gym *dto.GymCreationDTO) (*string, error) {
	query := `
		INSERT INTO gym (name, email, address, phone, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, true, $5, $5)
		RETURNING id`

	var id string
	now := time.Now()
	err := r.db.QueryRow(
		query,
		gym.Name,
		gym.Email,
		gym.Address,
		gym.Phone,
		now,
	).Scan(&id)

	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *GymRepository) GetGymByID(id string) (*dto.GymResponseDTO, error) {
	query := `
		SELECT id, name, email, address, phone, is_active, created_at, updated_at
		FROM gym 
		WHERE id = $1 AND deleted_at IS NULL`

	gym := &dto.GymResponseDTO{}

	err := r.db.QueryRow(query, id).Scan(
		&gym.ID,
		&gym.Name,
		&gym.Email,
		&gym.Address,
		&gym.Phone,
		&gym.IsActive,
		&gym.CreatedAt,
		&gym.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return gym, nil
}

func (r *GymRepository) GetGymByName(name string) (*dto.GymResponseDTO, error) {
	query := `
		SELECT id, name, email, address, phone, is_active, created_at, updated_at
		FROM gym 
		WHERE name = $1 AND deleted_at IS NULL`

	gym := &dto.GymResponseDTO{}
	err := r.db.QueryRow(query, name).Scan(
		&gym.ID,
		&gym.Name,
		&gym.Email,
		&gym.Address,
		&gym.Phone,
		&gym.IsActive,
		&gym.CreatedAt,
		&gym.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return gym, nil
}

func (r *GymRepository) GetAllGyms() ([]*dto.GymResponseDTO, error) {
	query := `
		SELECT id, name, email, address, phone, is_active, created_at, updated_at, deleted_at
		FROM gym 
		ORDER BY 
			CASE WHEN deleted_at IS NULL THEN 0 ELSE 1 END,
			created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gyms []*dto.GymResponseDTO
	for rows.Next() {
		gym := &dto.GymResponseDTO{}

		err := rows.Scan(
			&gym.ID,
			&gym.Name,
			&gym.Email,
			&gym.Address,
			&gym.Phone,
			&gym.IsActive,
			&gym.CreatedAt,
			&gym.UpdatedAt,
			&gym.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		gyms = append(gyms, gym)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gyms, nil
}

func (r *GymRepository) UpdateGym(id string, gym *dto.GymUpdateDTO) (*dto.GymResponseDTO, error) {
	query := `
		UPDATE gym 
		SET name = $1, email = $2, address = $3, phone = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, name, email, address, phone, is_active, created_at, updated_at`

	var updatedGym dto.GymResponseDTO
	err := r.db.QueryRow(query,
		gym.Name,
		gym.Email,
		gym.Address,
		gym.Phone,
		time.Now(),
		id,
	).Scan(
		&updatedGym.ID,
		&updatedGym.Name,
		&updatedGym.Email,
		&updatedGym.Address,
		&updatedGym.Phone,
		&updatedGym.IsActive,
		&updatedGym.CreatedAt,
		&updatedGym.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updatedGym, nil
}

func (r *GymRepository) SetGymActive(id string, active bool) error {
	query := `
		UPDATE gym 
		SET is_active = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, active, time.Now(), id)
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

func (r *GymRepository) DeleteGym(id string) error {
	query := `
		UPDATE gym 
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, time.Now(), id)
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
