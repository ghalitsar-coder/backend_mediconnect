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
		Table("doctors d").
		Select("d.id, d.facility_id, u.full_name as name, d.speciality as specialization, 'Poli ' || d.speciality as poli_name, 4.5 as rating, 120 as patients_count").
		Joins("JOIN users u ON d.user_id = u.id").
		Where("d.is_active = ?", true)

	if facilityID != "" {
		query = query.Where("d.facility_id = ?", facilityID)
	}

	// Assuming poliName matches specialization or we mock it. The query needs adaptation based on actual DB structure.
	// For now, we fetch base on facility.

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
