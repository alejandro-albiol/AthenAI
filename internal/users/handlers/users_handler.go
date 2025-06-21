package handlers

import (
	"encoding/json"
	"net/http"

	e "errors"

	"github.com/alejandro-albiol/athenai/internal/users/interfaces"
	"github.com/alejandro-albiol/athenai/internal/users/services"
	"github.com/alejandro-albiol/athenai/pkg/errors"
	errorsconst "github.com/alejandro-albiol/athenai/pkg/errors/const"
	"github.com/alejandro-albiol/athenai/pkg/responses"
)

type UsersHandler struct {
	service *services.UsersService
}

func NewUsersHandler(service *services.UsersService) *UsersHandler {
	return &UsersHandler{service: service}
}

func (h *UsersHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var dto interfaces.UserCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.APIResponse[any]{
			Status:  "error",
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}
	err := h.service.RegisterUser(dto)
	if err != nil {
		var apiErr errors.APIError
		if e.As(err, &apiErr) {
			status := http.StatusBadRequest
			switch apiErr.Code {
			case errorsconst.CodeNotFound:
				status = http.StatusNotFound
			case errorsconst.CodeConflict:
				status = http.StatusConflict
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(responses.APIResponse[any]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responses.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responses.APIResponse[any]{
		Status:  "success",
		Message: "User registered",
		Data:    dto,
	})
}

func (h *UsersHandler) GetUserByID(w http.ResponseWriter, id string) {
	user, err := h.service.GetUserByID(id)
	if err != nil {
		var apiErr errors.APIError
		if e.As(err, &apiErr) {
			status := http.StatusBadRequest
			if apiErr.Code == errorsconst.CodeNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(responses.APIResponse[any]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responses.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.APIResponse[any]{
		Status:  "success",
		Message: "User found",
		Data:    user,
	})
}

func (h *UsersHandler) GetUserByUsername(w http.ResponseWriter, username string) {
	user, err := h.service.GetUserByUsername(username)
	if err != nil {
		var apiErr errors.APIError
		if e.As(err, &apiErr) {
			status := http.StatusBadRequest
			if apiErr.Code == errorsconst.CodeNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(responses.APIResponse[any]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responses.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.APIResponse[any]{
		Status:  "success",
		Message: "User found",
		Data:    user,
	})
}

func (h *UsersHandler) GetUserByEmail(w http.ResponseWriter, email string) {
	user, err := h.service.GetUserByEmail(email)
	if err != nil {
		var apiErr errors.APIError
		if e.As(err, &apiErr) {
			status := http.StatusBadRequest
			if apiErr.Code == errorsconst.CodeNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(responses.APIResponse[any]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responses.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.APIResponse[any]{
		Status:  "success",
		Message: "User found",
		Data:    user,
	})
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, user interfaces.User) {
	err := h.service.UpdateUser(user)
	if err != nil {
		var apiErr errors.APIError
		if e.As(err, &apiErr) {
			status := http.StatusBadRequest
			switch apiErr.Code {
			case errorsconst.CodeNotFound:
				status = http.StatusNotFound
			case errorsconst.CodeConflict:
				status = http.StatusConflict
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(responses.APIResponse[any]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responses.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.APIResponse[any]{
		Status:  "success",
		Message: "User updated",
		Data:    user,
	})
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, id string) {
	err := h.service.DeleteUser(id)
	if err != nil {
		var apiErr errors.APIError
		if e.As(err, &apiErr) {
			status := http.StatusBadRequest
			if apiErr.Code == errorsconst.CodeNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(responses.APIResponse[any]{
				Status:  "error",
				Message: apiErr.Message,
				Data:    apiErr,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responses.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
