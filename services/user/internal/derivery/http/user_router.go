package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-market/services/user/internal/derivery/http/middleware"
)

func RegisterUserRoutes(r chi.Router, h *UserHandler, secret string) {
	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(secret))

		r.Get("/me", h.GetMe)
		r.Get("/{id}", h.GetByID)
		r.Get("/", h.GetByEmail)
		r.Put("/{id}", h.Update)
		
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Delete("/{id}", h.Delete)
		})
	})
}
