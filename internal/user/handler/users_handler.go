package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
	"github.com/alejandro-albiol/athenai/internal/user/service"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type UsersHandler struct {
	service *service.UsersService
}

func NewUsersHandler(service *service.UsersService) *UsersHandler {
	return &UsersHandler{service: service}
}

func (h *UsersHandler) RegisterUser(w http.ResponseWriter, r *http.Request, gymID string) {
	var dto dto.UserCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}
	err := h.service.RegisterUser(gymID, dto)
	if err != nil {
		var apiErr apierror.APIError
		if errors.As(err, &apiErr) {
			status := http.StatusBadRequest
			switch apiErr.Code {
			case errorcode_enum.CodeNotFound:
				status = http.StatusNotFound
			case errorcode_enum.CodeConflict:
				status = http.StatusConflict
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(response.APIResponse[any]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.APIResponse[any]{
		Status:  "success",
		Message: "User registered",
		Data:    dto,
	})
}

func (h *UsersHandler) GetAllUsers(w http.ResponseWriter, gymID string) {
	users, err := h.service.GetAllUsers(gymID)
	if err != nil {
		var apiErr apierror.APIError
		if errors.As(err, &apiErr) {
			status := http.StatusBadRequest
			if apiErr.Code == errorcode_enum.CodeNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(response.APIResponse[apierror.APIError]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.APIResponse[[]dto.UserResponseDTO]{
		Status:  "success",
		Message: "Users retrieved",
		Data:    users,
	})
}

func (h *UsersHandler) GetUserByID(w http.ResponseWriter, gymID, id string) {
	user, err := h.service.GetUserByID(gymID, id)
	if err != nil {
		var apiErr apierror.APIError
		if errors.As(err, &apiErr) {
			status := http.StatusBadRequest
			if apiErr.Code == errorcode_enum.CodeNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(response.APIResponse[apierror.APIError]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.APIResponse[dto.UserResponseDTO]{
		Status:  "success",
		Message: "User found",
		Data:    user,
	})
}

func (h *UsersHandler) GetUserByUsername(w http.ResponseWriter, gymID, username string) {
	user, err := h.service.GetUserByUsername(gymID, username)
	if err != nil {
		var apiErr apierror.APIError
		if errors.As(err, &apiErr) {
			status := http.StatusBadRequest
			if apiErr.Code == errorcode_enum.CodeNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(response.APIResponse[apierror.APIError]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.APIResponse[dto.UserResponseDTO]{
		Status:  "success",
		Message: "User found",
		Data:    user,
	})
}

func (h *UsersHandler) GetUserByEmail(w http.ResponseWriter, gymID, email string) {
	user, err := h.service.GetUserByEmail(gymID, email)
	if err != nil {
		var apiErr apierror.APIError
		if errors.As(err, &apiErr) {
			status := http.StatusBadRequest
			if apiErr.Code == errorcode_enum.CodeNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(response.APIResponse[apierror.APIError]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.APIResponse[dto.UserResponseDTO]{
		Status:  "success",
		Message: "User found",
		Data:    user,
	})
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, gymID, id string, userDTO dto.UserUpdateDTO) {
	err := h.service.UpdateUser(gymID, id, userDTO)
	if err != nil {
		var apiErr apierror.APIError
		if errors.As(err, &apiErr) {
			status := http.StatusBadRequest
			switch apiErr.Code {
			case errorcode_enum.CodeNotFound:
				status = http.StatusNotFound
			case errorcode_enum.CodeConflict:
				status = http.StatusConflict
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(response.APIResponse[any]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.APIResponse[any]{
		Status:  "success",
		Message: "User updated",
		Data:    nil,
	})
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, gymID, id string) {
	err := h.service.DeleteUser(gymID, id)
	if err != nil {
		var apiErr apierror.APIError
		if errors.As(err, &apiErr) {
			status := http.StatusBadRequest
			if apiErr.Code == errorcode_enum.CodeNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(response.APIResponse[any]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
