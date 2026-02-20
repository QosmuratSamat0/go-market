package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	userErr "github.com/go-market/pkg/errs"
	"github.com/go-market/services/user/internal/model"
	"github.com/go-market/services/user/internal/service"
)

type UserHandler struct {
	log *slog.Logger
	svc *service.Service
}

func New(log *slog.Logger, svc *service.Service) *UserHandler {
	return &UserHandler{
		log: log,
		svc: svc,
	}
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar" validate:"omitempty"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar" validate:"omitempty"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.GetByID"
	log := h.log.With(slog.String("op", op))

	id := chi.URLParam(r, "id")

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		log.Error("failed to get user by id", slog.String("error", err.Error()))
		if errors.Is(err, userErr.ErrInvalidID) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, userErr.ErrUserNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{Error: "failed to get user"})
		return
	}

	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	render.JSON(w, r, SuccessResponse{Data: response})
}

func (h *UserHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.GetByEmail"
	log := h.log.With(slog.String("op", op))

	email := chi.URLParam(r, "email")

	user, err := h.svc.GetByEmail(r.Context(), email)
	if err != nil {
		log.Error("failed to get user by email", slog.String("error", err.Error()))
		if errors.Is(err, userErr.ErrInvalidEmail) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, userErr.ErrUserNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{Error: "failed to get user"})
		return
	}

	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	render.JSON(w, r, SuccessResponse{Data: response})
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.Create"
	log := h.log.With(slog.String("op", op))

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: "invalid request body"})
		return
	}

	user := model.User{
		Username:  req.Username,
		Email:     req.Email,
		Avatar:    req.Avatar,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := h.svc.Create(r.Context(), user)
	if err != nil {
		log.Error("failed to create user", slog.String("error", err.Error()))
		if errors.Is(err, userErr.ErrInvalidEmail) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, userErr.ErrUserExists) {
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{Error: "failed to create user"})
		return
	}

	response := map[string]string{
		"id": id,
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, SuccessResponse{Data: response})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.Update"
	log := h.log.With(slog.String("op", op))

	id := chi.URLParam(r, "id")

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: "invalid request body"})
		return
	}

	user := model.User{
		ID:        id,
		Username:  req.Username,
		Email:     req.Email,
		Avatar:    req.Avatar,
		UpdatedAt: time.Now(),
	}

	err := h.svc.Update(r.Context(), user)
	if err != nil {
		log.Error("failed to update user", slog.String("error", err.Error()))
		if errors.Is(err, userErr.ErrInvalidID) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, userErr.ErrUserNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{Error: "failed to update user"})
		return
	}

	render.JSON(w, r, SuccessResponse{Message: "user updated successfully"})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.Delete"
	log := h.log.With(slog.String("op", op))

	id := chi.URLParam(r, "id")

	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		log.Error("failed to delete user", slog.String("error", err.Error()))
		if errors.Is(err, userErr.ErrInvalidID) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, userErr.ErrUserNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{Error: "failed to delete user"})
		return
	}

	render.JSON(w, r, SuccessResponse{Message: "user deleted successfully"})
}
