package services

import (
	"fmt"
	"github.com/alejandro-albiol/athenai/internal/users/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/errors"
	errorsconst "github.com/alejandro-albiol/athenai/pkg/errors/const"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repository interfaces.UsersRepository
}

func NewUsersService(repository interfaces.UsersRepository) *UsersService {
	return &UsersService{repository: repository}
}

func (s *UsersService) RegisterUser(user interfaces.UserCreationDTO) error {
	existingUsername, err := s.repository.GetUserByUsername(user.Username)
	if err == nil && existingUsername.ID != "" {
		return errors.New(errorsconst.CodeConflict, fmt.Sprintf("Username %s already exists", user.Username))
	}

	existingEmail, err := s.repository.GetUserByEmail(user.Email)
	if err == nil && existingEmail.ID != "" {
		return errors.New(errorsconst.CodeConflict, fmt.Sprintf("Email %s already exists", user.Email))
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.repository.CreateUser(user)
}

func (s *UsersService) GetUserByID(id string) (interfaces.User, error) {
	user, err := s.repository.GetUserByID(id)
	if err != nil {
		return interfaces.User{}, errors.New(errorsconst.CodeNotFound, fmt.Sprintf("User with ID %s not found", id))
	}
	return user, nil
}

func (s *UsersService) GetUserByUsername(username string) (interfaces.User, error) {
	user, err := s.repository.GetUserByUsername(username)
	if err != nil {
		return interfaces.User{}, errors.New(errorsconst.CodeNotFound, fmt.Sprintf("User with username %s not found", username))
	}
	return user, nil
}

func (s *UsersService) GetUserByEmail(email string) (interfaces.User, error) {
	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return interfaces.User{}, errors.New(errorsconst.CodeNotFound, fmt.Sprintf("User with email %s not found", email))
	}
	return user, nil
}

func (s *UsersService) UpdateUser(user interfaces.User) error {
	existingUser, err := s.repository.GetUserByID(user.ID)
	if err != nil || existingUser.ID == "" {
		return errors.New(errorsconst.CodeNotFound, fmt.Sprintf("User with ID %s not found", user.ID))
	}

	if user.Username != existingUser.Username {
		existingUsername, err := s.repository.GetUserByUsername(user.Username)
		if err == nil && existingUsername.ID != "" {
			return errors.New(errorsconst.CodeConflict, fmt.Sprintf("Username %s already exists", user.Username))
		}
	}

	if user.Email != existingUser.Email {
		existingEmail, err := s.repository.GetUserByEmail(user.Email)
		if err == nil && existingEmail.ID != "" {
			return errors.New(errorsconst.CodeConflict, fmt.Sprintf("Email %s already exists", user.Email))
		}
	}

	return s.repository.UpdateUser(user)
}

func (s *UsersService) DeleteUser(id string) error {
	existingUser, err := s.repository.GetUserByID(id)
	if err != nil || existingUser.ID == "" {
		return errors.New(errorsconst.CodeNotFound, fmt.Sprintf("User with ID %s not found", id))
	}
	return s.repository.DeleteUser(id)
}
