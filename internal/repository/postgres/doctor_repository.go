package postgres

import (
	"context"
	"mediconnect/internal/domain"

	"gorm.io/gorm"
)

type DoctorRepository struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) domain.DoctorRepository {
	return &DoctorRepository{db: db}
}

func (r *DoctorRepository) GetDoctors(ctx context.Context, facilityID string, poliName string) ([]domain.Doctor, error) {
	var doctors []domain.Doctor

	query := r.db.WithContext(ctx).
		Table("doctors").
		Select("id, facility_id, name, specialization, poli_name, rating, patients_count")

	if facilityID != "" {
		query = query.Where("facility_id = ?", facilityID)
	}

	if poliName != "" {
		query = query.Where("specialization = ?", poliName)
	}

	if err := query.Find(&doctors).Error; err != nil {
		return nil, err
	}

	return doctors, nil
}

func (r *DoctorRepository) GetDoctorSchedules(ctx context.Context, doctorID string, dayOfWeek int) ([]domain.DoctorSchedule, error) {
	// Mocking schedules as there is no schedules table right now in the database
	schedules := []domain.DoctorSchedule{
		{
			ID:                  "sched-1",
			DoctorID:            doctorID,
			DayOfWeek:           dayOfWeek,
			StartTime:           "08:00",
			EndTime:             "16:00",
			SlotDurationMinutes: 30,
			MaxPatients:         20,
		},
	}
	return schedules, nil
}
