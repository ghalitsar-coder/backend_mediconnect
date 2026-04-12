package domain

import (
	"context"
	"time"
)

// ─── Entities ────────────────────────────────────────────────────────────────

// Facility represents a healthcare facility (Puskesmas or Klinik).
type Facility struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	Type       string    `json:"type"`
	DistrictID string    `json:"district_id"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FacilityFilter holds optional query filters when listing facilities.
type FacilityFilter struct {
	DistrictID string
	Type       string
}

// ─── Repository Contract ─────────────────────────────────────────────────────

// FacilityRepository defines the data-access contract for facility persistence.
// Implementations live in internal/repository/postgres/.
type FacilityRepository interface {
	GetFacilities(ctx context.Context, filter FacilityFilter) ([]Facility, error)
}

// ─── Usecase Contract ────────────────────────────────────────────────────────

// FacilityUsecase defines the business-logic contract for facility operations.
// Implementations live in internal/usecase/.
type FacilityUsecase interface {
	GetFacilities(ctx context.Context, filter FacilityFilter) ([]Facility, error)
}
