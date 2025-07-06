package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
	"github.com/alejandro-albiol/athenai/internal/user/interfaces"
)

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) interfaces.UserRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) CreateUser(dto dto.UserCreationDTO, gymID string) error {
	query := "INSERT INTO users (username, email, password_hash, role, gym_id) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, dto.Username, dto.Email, dto.Password, dto.Role, gymID)
	return err
}

func (r *usersRepository) GetUserByID(gymID, id string) (dto.UserResponseDTO, error) {
	query := "SELECT id, username, email, role FROM users WHERE id = $1 AND gym_id = $2"
	row := r.db.QueryRow(query, id, gymID)
	var user dto.UserResponseDTO
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}
	return user, nil
}

func (r *usersRepository) GetUserByUsername(gymID, username string) (dto.UserResponseDTO, error) {
	query := "SELECT id, username, email FROM users WHERE username = $1 AND gym_id = $2"
	row := r.db.QueryRow(query, username, gymID)
	var user dto.UserResponseDTO
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}
	return user, nil
}

func (r *usersRepository) GetUserByEmail(gymID, email string) (dto.UserResponseDTO, error) {
	query := "SELECT id, username, email FROM users WHERE email = $1 AND gym_id = $2"
	row := r.db.QueryRow(query, email, gymID)
	var user dto.UserResponseDTO
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}
	return user, nil
}

func (r *usersRepository) GetAllUsers(gymID string) ([]dto.UserResponseDTO, error) {
	query := "SELECT id, username, email, role FROM users WHERE gym_id = $1"
	rows, err := r.db.Query(query, gymID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []dto.UserResponseDTO
	for rows.Next() {
		var user dto.UserResponseDTO
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *usersRepository) GetPasswordHashByUsername(gymID, username string) (string, error) {
	query := "SELECT password_hash FROM users WHERE username = $1 AND gym_id = $2"
	row := r.db.QueryRow(query, username, gymID)
	var passwordHash string
	err := row.Scan(&passwordHash)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (r *usersRepository) UpdateUser(gymID string, id string, user dto.UserUpdateDTO) error {
	query := "UPDATE users SET username = $1, email = $2 WHERE id = $3 AND gym_id = $4"
	_, err := r.db.Exec(query, user.Username, user.Email, id, gymID)
	return err
}

func (r *usersRepository) UpdatePassword(gymID, userID string, newPasswordHash string) error {
	query := "UPDATE users SET password_hash = $1 WHERE id = $2 AND gym_id = $3"
	_, err := r.db.Exec(query, newPasswordHash, userID, gymID)
	return err
}

func (r *usersRepository) DeleteUser(gymID, id string) error {
	query := "DELETE FROM users WHERE id = $1 AND gym_id = $2"
	_, err := r.db.Exec(query, id, gymID)
	return err
}

