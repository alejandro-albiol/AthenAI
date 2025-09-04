package service

import (
	"fmt"

	dto "github.com/alejandro-albiol/athenai/internal/user/dto"
	"github.com/alejandro-albiol/athenai/internal/user/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repository interfaces.UserRepository
}

func NewUsersService(repository interfaces.UserRepository) *UsersService {
	return &UsersService{repository: repository}
}

func (s *UsersService) RegisterUser(gymID string, user *dto.UserCreationDTO) (*string, error) {
	existingUsername, err := s.repository.GetUserByUsername(gymID, user.Username)
	if err == nil && existingUsername != nil && existingUsername.ID != "" {
		return nil, apierror.New(errorcode_enum.CodeConflict, fmt.Sprintf("Username %s already exists", user.Username), nil)
	}

	existingEmail, err := s.repository.GetUserByEmail(gymID, user.Email)
	if err == nil && existingEmail != nil && existingEmail.ID != "" {
		return nil, apierror.New(errorcode_enum.CodeConflict, fmt.Sprintf("Email %s already exists", user.Email), nil)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to hash password", err)
	}
	user.Password = string(hashedPassword)

	return s.repository.CreateUser(gymID, user)
}

func (s *UsersService) GetUserByID(gymID, id string) (*dto.UserResponseDTO, error) {
	user, err := s.repository.GetUserByID(gymID, id)
	if err != nil || user == nil {
		return nil, apierror.New(errorcode_enum.CodeNotFound, fmt.Sprintf("User with ID %s not found", id), err)
	}
	return user, nil
}

func (s *UsersService) GetUserByUsername(gymID, username string) (*dto.UserResponseDTO, error) {
	user, err := s.repository.GetUserByUsername(gymID, username)
	if err != nil || user == nil {
		return nil, apierror.New(errorcode_enum.CodeNotFound, fmt.Sprintf("User with username %s not found", username), err)
	}
	return user, nil
}

func (s *UsersService) GetUserByEmail(gymID, email string) (*dto.UserResponseDTO, error) {
	user, err := s.repository.GetUserByEmail(gymID, email)
	if err != nil || user == nil {
		return nil, apierror.New(errorcode_enum.CodeNotFound, fmt.Sprintf("User with email %s not found", email), err)
	}
	return user, nil
}

func (s *UsersService) GetPasswordHashByUsername(gymID, username string) (string, error) {
	passwordHash, err := s.repository.GetPasswordHashByUsername(gymID, username)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeNotFound, fmt.Sprintf("User with username %s not found", username), err)
	}
	return passwordHash, nil
}

func (s *UsersService) GetAllUsers(gymID string) ([]*dto.UserResponseDTO, error) {
	users, err := s.repository.GetAllUsers(gymID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve users", err)
	}
	return users, nil
}

func (s *UsersService) UpdateUser(gymID, id string, user *dto.UserUpdateDTO) error {
	existingUser, err := s.repository.GetUserByID(gymID, id)
	if err != nil || existingUser == nil || existingUser.ID == "" {
		return apierror.New(errorcode_enum.CodeNotFound, fmt.Sprintf("User with ID %s not found", id), err)
	}

	if user.Username != existingUser.Username {
		existingUsername, err := s.repository.GetUserByUsername(gymID, user.Username)
		if err == nil && existingUsername != nil && existingUsername.ID != "" {
			return apierror.New(errorcode_enum.CodeConflict, fmt.Sprintf("Username %s already exists", user.Username), nil)
		}
	}

	if user.Email != existingUser.Email {
		existingEmail, err := s.repository.GetUserByEmail(gymID, user.Email)
		if err == nil && existingEmail != nil && existingEmail.ID != "" {
			return apierror.New(errorcode_enum.CodeConflict, fmt.Sprintf("Email %s already exists", user.Email), nil)
		}
	}

	return s.repository.UpdateUser(gymID, id, user)
}

func (s *UsersService) UpdatePassword(gymID, id, newPassword string) error {
	existingUser, err := s.repository.GetUserByID(gymID, id)
	if err != nil || existingUser == nil || existingUser.ID == "" {
		return apierror.New(errorcode_enum.CodeNotFound, fmt.Sprintf("User with ID %s not found", id), err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to hash password", err)
	}

	return s.repository.UpdatePassword(gymID, id, string(hashedPassword))
}

func (s *UsersService) DeleteUser(gymID, id string) error {
	existingUser, err := s.repository.GetUserByID(gymID, id)
	if err != nil || existingUser == nil || existingUser.ID == "" {
		return apierror.New(errorcode_enum.CodeNotFound, fmt.Sprintf("User with ID %s not found", id), err)
	}
	return s.repository.DeleteUser(gymID, id)
}

func (s *UsersService) VerifyUser(gymID, userID string) error {
	existingUser, err := s.repository.GetUserByID(gymID, userID)
	if err != nil || existingUser == nil || existingUser.ID == "" {
		return apierror.New(errorcode_enum.CodeNotFound, fmt.Sprintf("User with ID %s not found", userID), err)
	}

	if existingUser.Verified {
		return apierror.New(errorcode_enum.CodeConflict, "User is already verified", nil)
	}

	err = s.repository.VerifyUser(gymID, userID)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to verify user", err)
	}

	return nil
}

func (s *UsersService) SetUserActive(gymID, userID string, active bool) error {
	existingUser, err := s.repository.GetUserByID(gymID, userID)
	if err != nil || existingUser == nil || existingUser.ID == "" {
		return apierror.New(errorcode_enum.CodeNotFound, fmt.Sprintf("User with ID %s not found", userID), err)
	}

	if existingUser.IsActive == active {
		status := "active"
		if !active {
			status = "inactive"
		}
		return apierror.New(errorcode_enum.CodeConflict, fmt.Sprintf("User is already %s", status), nil)
	}

	err = s.repository.SetUserActive(gymID, userID, active)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update user status", err)
	}

	return nil
}
