package http

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"mediconnect/internal/delivery/http/handler"
	"mediconnect/internal/delivery/http/middleware"
)

// NewRouter constructs the chi router with all middleware and routes registered.
func NewRouter(facilityHandler *handler.FacilityHandler, authHandler *handler.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	// ── Global Middleware Stack ──────────────────────────────────────────────
	r.Use(chimw.RequestID) // inject X-Request-Id header
	r.Use(chimw.RealIP)    // trust X-Real-IP / X-Forwarded-For
	r.Use(chimw.Logger)    // structured request logging
	r.Use(chimw.Recoverer) // recover from panics, return 500

	// ── Routes ──────────────────────────────────────────────────────────────
	r.Get("/api/v1/health", handler.HealthHandler)

	r.Route("/api/v1", func(r chi.Router) {
		// Public Auth Routes
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/logout", authHandler.Logout)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth) // authentication layer

			// Facilities (Protected if you want, or put outside this group if public)
			r.Get("/facilities", facilityHandler.GetFacilities)
		})
	})

	return r
}
