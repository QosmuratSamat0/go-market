package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, h *UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/{id}", h.GetByID)
		r.Get("/", h.GetByEmail)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
