package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates a new chi router with specified options.
func New() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

// GroupRouter creates a new chi router group within the parent router.
func GroupRouter(parent chi.Router, pattern string, fn func(r chi.Router)) chi.Router {
	return parent.Route(pattern, fn)
}

// UseMiddleware adds middleware to the router.
func UseMiddleware(r chi.Router, m ...func(http.Handler) http.Handler) {
	r.Use(m...)
}

// MountRouter mounts a sub-router onto the parent router.
func MountRouter(parent chi.Router, pattern string, handler http.Handler) {
	parent.Mount(pattern, handler)
}
