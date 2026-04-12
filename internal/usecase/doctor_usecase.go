package usecase

import (
	"context"
	"time"

	"mediconnect/internal/domain"
)

type DoctorUsecase struct {
	repo domain.DoctorRepository
}

func NewDoctorUsecase(repo domain.DoctorRepository) domain.DoctorUsecase {
	return &DoctorUsecase{repo: repo}
}

func (u *DoctorUsecase) GetDoctors(ctx context.Context, facilityID string, poliName string) ([]domain.Doctor, error) {
	return u.repo.GetDoctors(ctx, facilityID, poliName)
}

func (u *DoctorUsecase) GetAvailableSlots(ctx context.Context, doctorID string, date time.Time) ([]domain.Slot, error) {
	// Basic mock implementation for slots
	slots := []domain.Slot{
		{Time: "08:00", IsAvailable: true},
		{Time: "08:30", IsAvailable: true},
		{Time: "09:00", IsAvailable: false},
		{Time: "09:30", IsAvailable: true},
	}
	return slots, nil
}
