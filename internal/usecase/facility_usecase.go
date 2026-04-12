package usecase

import (
	"context"

	"mediconnect/internal/domain"
)

type facilityUsecase struct {
	facilityRepo domain.FacilityRepository
}

// NewFacilityUsecase returns a domain.FacilityUsecase backed by the given repository.
func NewFacilityUsecase(repo domain.FacilityRepository) domain.FacilityUsecase {
	return &facilityUsecase{facilityRepo: repo}
}

func (u *facilityUsecase) GetFacilities(ctx context.Context, filter domain.FacilityFilter) ([]domain.Facility, error) {
	// Business rules can be added here (e.g. audit, caching, validation)
	// before delegating to the repository.
	return u.facilityRepo.GetFacilities(ctx, filter)
}
