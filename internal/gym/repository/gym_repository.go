package repository

import (
	"database/sql"
	"time"

	"github.com/alejandro-albiol/athenai/internal/gym/dto"
	"github.com/lib/pq"
)

type GymRepository struct {
	db *sql.DB
}

func NewGymRepository(db *sql.DB) *GymRepository {
	return &GymRepository{db: db}
}

func (r *GymRepository) CreateGym(gym dto.GymCreationDTO) (string, error) {
	query := `
		INSERT INTO gyms (name, domain, email, address, contact_name, phone, logo_url,
			is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, true, $8, $8)
		RETURNING id`

	var id string
	now := time.Now()
	err := r.db.QueryRow(
		query,
		gym.Name,
		gym.Domain,
		gym.Email,
		gym.Address,
		gym.ContactName,
		gym.Phone,
		gym.LogoURL,
		now,
	).Scan(&id)

	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *GymRepository) GetGymByID(id string) (dto.GymResponseDTO, error) {
	query := `
		SELECT id, name, domain, email, address, contact_name, phone, logo_url, 
			description, business_hours, social_links, payment_methods, currency, 
			timezone_offset, is_active, created_at, updated_at
		FROM gyms 
		WHERE id = $1 AND deleted_at IS NULL`

	var gym dto.GymResponseDTO
	err := r.db.QueryRow(query, id).Scan(
		&gym.ID,
		&gym.Name,
		&gym.Domain,
		&gym.Email,
		&gym.Address,
		&gym.ContactName,
		&gym.Phone,
		&gym.LogoURL,
		&gym.Description,
		pq.Array(&gym.BusinessHours),
		pq.Array(&gym.SocialLinks),
		pq.Array(&gym.PaymentMethods),
		&gym.Currency,
		&gym.TimezoneOffset,
		&gym.IsActive,
		&gym.CreatedAt,
		&gym.UpdatedAt,
	)
	if err != nil {
		return dto.GymResponseDTO{}, err
	}

	return gym, nil
}

func (r *GymRepository) GetGymByDomain(domain string) (dto.GymResponseDTO, error) {
	query := `
		SELECT id, name, domain, email, address, contact_name, phone, logo_url, 
			description, business_hours, social_links, payment_methods, currency, 
			timezone_offset, is_active, created_at, updated_at
		FROM gyms 
		WHERE domain = $1 AND deleted_at IS NULL`

	var gym dto.GymResponseDTO
	err := r.db.QueryRow(query, domain).Scan(
		&gym.ID,
		&gym.Name,
		&gym.Domain,
		&gym.Email,
		&gym.Address,
		&gym.ContactName,
		&gym.Phone,
		&gym.LogoURL,
		&gym.Description,
		pq.Array(&gym.BusinessHours),
		pq.Array(&gym.SocialLinks),
		pq.Array(&gym.PaymentMethods),
		&gym.Currency,
		&gym.TimezoneOffset,
		&gym.IsActive,
		&gym.CreatedAt,
		&gym.UpdatedAt,
	)

	if err != nil {
		return dto.GymResponseDTO{}, err
	}

	return gym, nil
}

func (r *GymRepository) GetAllGyms() ([]dto.GymResponseDTO, error) {
	query := `
		SELECT id, name, domain, email, address, contact_name, phone, logo_url, 
			description, business_hours, social_links, payment_methods, currency, 
			timezone_offset, is_active, created_at, updated_at
		FROM gyms 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gyms []dto.GymResponseDTO
	for rows.Next() {
		var gym dto.GymResponseDTO
		err := rows.Scan(
			&gym.ID,
			&gym.Name,
			&gym.Domain,
			&gym.Email,
			&gym.Address,
			&gym.ContactName,
			&gym.Phone,
			&gym.LogoURL,
			&gym.Description,
			pq.Array(&gym.BusinessHours),
			pq.Array(&gym.SocialLinks),
			pq.Array(&gym.PaymentMethods),
			&gym.Currency,
			&gym.TimezoneOffset,
			&gym.IsActive,
			&gym.CreatedAt,
			&gym.UpdatedAt,
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
func (r *GymRepository) UpdateGym(id string, gym dto.GymUpdateDTO) (dto.GymResponseDTO, error) {
	query := `
		UPDATE gyms 
		SET name = $1, email = $2, address = $3, contact_name = $4, 
			phone = $5, logo_url = $6, description = $7, business_hours = $8, 
			social_links = $9, payment_methods = $10, currency = $11, 
			timezone_offset = $12, updated_at = $13
		WHERE id = $14 AND deleted_at IS NULL
		RETURNING id, name, email, address, contact_name, phone, logo_url,
			description, business_hours, social_links, payment_methods,
			currency, timezone_offset, is_active, created_at, updated_at`

	var updatedGym dto.GymResponseDTO
	err := r.db.QueryRow(query,
		gym.Name,
		gym.Email,
		gym.Address,
		gym.ContactName,
		gym.Phone,
		gym.LogoURL,
		gym.Description,
		pq.Array(gym.BusinessHours),
		pq.Array(gym.SocialLinks),
		pq.Array(gym.PaymentMethods),
		gym.Currency,
		gym.TimezoneOffset,
		time.Now(),
		id,
	).Scan(
		&updatedGym.ID,
		&updatedGym.Name,
		&updatedGym.Domain,
		&updatedGym.Email,
		&updatedGym.Address,
		&updatedGym.ContactName,
		&updatedGym.Phone,
		&updatedGym.LogoURL,
		&updatedGym.Description,
		pq.Array(&updatedGym.BusinessHours),
		pq.Array(&updatedGym.SocialLinks),
		pq.Array(&updatedGym.PaymentMethods),
		&updatedGym.Currency,
		&updatedGym.TimezoneOffset,
		&updatedGym.IsActive,
		&updatedGym.CreatedAt,
		&updatedGym.UpdatedAt,
	)

	if err != nil {
		return dto.GymResponseDTO{}, err
	}

	return updatedGym, nil
}
func (r *GymRepository) SetGymActive(id string, active bool) error {
	query := `
		UPDATE gyms 
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
		UPDATE gyms 
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
