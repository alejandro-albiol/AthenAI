package repository

import (
	dbinterfaces "github.com/alejandro-albiol/athenai/internal/databases/interfaces"
	"github.com/alejandro-albiol/athenai/internal/users/interfaces"
)

type UsersRepository struct {
	db dbinterfaces.DBService
}

func NewUsersRepository(db dbinterfaces.DBService) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) CreateUser(dto interfaces.UserCreationDTO) error {
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, dto.Username, dto.Email, dto.Password)
	return err
}

func (r *UsersRepository) GetUserByID(id string) (interfaces.User, error) {
	query := "SELECT id, username, email FROM users WHERE id = ?"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return interfaces.User{}, err
	}
	defer rows.Close()

	var user interfaces.User
	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return interfaces.User{}, err
		}
	}
	return user, nil
}

func (r *UsersRepository) GetUserByUsername(username string) (interfaces.User, error) {
	query := "SELECT id, username, email FROM users WHERE username = ?"
	rows, err := r.db.Query(query, username)
	if err != nil {
		return interfaces.User{}, err
	}
	defer rows.Close()

	var user interfaces.User
	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return interfaces.User{}, err
		}
	}
	return user, nil
}

func (r *UsersRepository) GetUserByEmail(email string) (interfaces.User, error) {
	query := "SELECT id, username, email FROM users WHERE email = ?"
	rows, err := r.db.Query(query, email)
	if err != nil {
		return interfaces.User{}, err
	}
	defer rows.Close()

	var user interfaces.User
	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return interfaces.User{}, err
		}
	}
	return user, nil
}

func (r *UsersRepository) UpdateUser(user interfaces.User) error {
	query := "UPDATE users SET username = ?, email = ? WHERE id = ?"
	_, err := r.db.Exec(query, user.Username, user.Email, user.ID)
	return err
}

func (r *UsersRepository) DeleteUser(id string) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}