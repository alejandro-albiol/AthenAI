package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/alejandro-albiol/athenai/internal/gym/dto"
)

type GymRepository struct {
	db *sql.DB
}

func NewGymRepository(db *sql.DB) *GymRepository {
	return &GymRepository{db: db}
}

func (r *GymRepository) CreateGym(gym dto.GymCreationDTO) (string, error) {
	query := `
		INSERT INTO gym (name, domain, email, address, phone,
			business_hours, social_links, payment_methods,
			is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, '[]'::jsonb, '[]'::jsonb, '[]'::jsonb, true, $6, $6)
		RETURNING id`

	var id string
	now := time.Now()
	err := r.db.QueryRow(
		query,
		gym.Name,
		gym.Domain,
		gym.Email,
		gym.Address,
		gym.Phone,
		now,
	).Scan(&id)

	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *GymRepository) GetGymByID(id string) (dto.GymResponseDTO, error) {
	query := `
		SELECT id, name, domain, email, address, phone,
			business_hours, social_links, payment_methods, is_active, created_at, updated_at
		FROM gym 
		WHERE id = $1 AND deleted_at IS NULL`

	var gym dto.GymResponseDTO
	var businessHours, socialLinks, paymentMethods []byte

	err := r.db.QueryRow(query, id).Scan(
		&gym.ID,
		&gym.Name,
		&gym.Domain,
		&gym.Email,
		&gym.Address,
		&gym.Phone,
		&businessHours,
		&socialLinks,
		&paymentMethods,
		&gym.IsActive,
		&gym.CreatedAt,
		&gym.UpdatedAt,
	)
	if err != nil {
		return dto.GymResponseDTO{}, err
	}

	// Parse JSONB fields
	if err := json.Unmarshal(businessHours, &gym.BusinessHours); err != nil {
		gym.BusinessHours = []string{}
	}
	if err := json.Unmarshal(socialLinks, &gym.SocialLinks); err != nil {
		gym.SocialLinks = []string{}
	}
	if err := json.Unmarshal(paymentMethods, &gym.PaymentMethods); err != nil {
		gym.PaymentMethods = []string{}
	}

	return gym, nil
}

func (r *GymRepository) GetGymByDomain(domain string) (dto.GymResponseDTO, error) {
	query := `
		SELECT id, name, domain, email, address, phone,
			business_hours, social_links, payment_methods, is_active, created_at, updated_at
		FROM gym 
		WHERE domain = $1 AND deleted_at IS NULL`

	var gym dto.GymResponseDTO
	var businessHours, socialLinks, paymentMethods []byte

	err := r.db.QueryRow(query, domain).Scan(
		&gym.ID,
		&gym.Name,
		&gym.Domain,
		&gym.Email,
		&gym.Address,
		&gym.Phone,
		&businessHours,
		&socialLinks,
		&paymentMethods,
		&gym.IsActive,
		&gym.CreatedAt,
		&gym.UpdatedAt,
	)

	if err != nil {
		return dto.GymResponseDTO{}, err
	}

	// Parse JSONB fields
	if err := json.Unmarshal(businessHours, &gym.BusinessHours); err != nil {
		gym.BusinessHours = []string{}
	}
	if err := json.Unmarshal(socialLinks, &gym.SocialLinks); err != nil {
		gym.SocialLinks = []string{}
	}
	if err := json.Unmarshal(paymentMethods, &gym.PaymentMethods); err != nil {
		gym.PaymentMethods = []string{}
	}

	return gym, nil
}

func (r *GymRepository) GetAllGyms() ([]dto.GymResponseDTO, error) {
	query := `
		SELECT id, name, domain, email, address, phone,
			business_hours, social_links, payment_methods, is_active, created_at, updated_at
		FROM gym 
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
		var businessHours, socialLinks, paymentMethods []byte

		err := rows.Scan(
			&gym.ID,
			&gym.Name,
			&gym.Domain,
			&gym.Email,
			&gym.Address,
			&gym.Phone,
			&businessHours,
			&socialLinks,
			&paymentMethods,
			&gym.IsActive,
			&gym.CreatedAt,
			&gym.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSONB fields
		if err := json.Unmarshal(businessHours, &gym.BusinessHours); err != nil {
			gym.BusinessHours = []string{}
		}
		if err := json.Unmarshal(socialLinks, &gym.SocialLinks); err != nil {
			gym.SocialLinks = []string{}
		}
		if err := json.Unmarshal(paymentMethods, &gym.PaymentMethods); err != nil {
			gym.PaymentMethods = []string{}
		}

		gyms = append(gyms, gym)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gyms, nil
}
func (r *GymRepository) UpdateGym(id string, gym dto.GymUpdateDTO) (dto.GymResponseDTO, error) {
	// Convert arrays to JSON for JSONB fields
	businessHoursJSON, _ := json.Marshal(gym.BusinessHours)
	socialLinksJSON, _ := json.Marshal(gym.SocialLinks)
	paymentMethodsJSON, _ := json.Marshal(gym.PaymentMethods)

	query := `
		UPDATE gym 
		SET name = $1, email = $2, address = $3, 
			phone = $4, business_hours = $5, social_links = $6, payment_methods = $7,
			updated_at = $8
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING id, name, domain, email, address, phone,
			business_hours, social_links, payment_methods, is_active, created_at, updated_at`

	var updatedGym dto.GymResponseDTO
	var businessHours, socialLinks, paymentMethods []byte

	err := r.db.QueryRow(query,
		gym.Name,
		gym.Email,
		gym.Address,
		gym.Phone,
		businessHoursJSON,
		socialLinksJSON,
		paymentMethodsJSON,
		time.Now(),
		id,
	).Scan(
		&updatedGym.ID,
		&updatedGym.Name,
		&updatedGym.Domain,
		&updatedGym.Email,
		&updatedGym.Address,
		&updatedGym.Phone,
		&businessHours,
		&socialLinks,
		&paymentMethods,
		&updatedGym.IsActive,
		&updatedGym.CreatedAt,
		&updatedGym.UpdatedAt,
	)

	if err != nil {
		return dto.GymResponseDTO{}, err
	}

	// Parse JSONB fields
	if err := json.Unmarshal(businessHours, &updatedGym.BusinessHours); err != nil {
		updatedGym.BusinessHours = []string{}
	}
	if err := json.Unmarshal(socialLinks, &updatedGym.SocialLinks); err != nil {
		updatedGym.SocialLinks = []string{}
	}
	if err := json.Unmarshal(paymentMethods, &updatedGym.PaymentMethods); err != nil {
		updatedGym.PaymentMethods = []string{}
	}

	return updatedGym, nil
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
		DELETE FROM gym 
		WHERE id = $1`

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
