package middleware

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	RoleKey contextKey = "role"
)

func RequireRole(required string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value(RoleKey).(string)
			if !ok || role != required {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, "forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
