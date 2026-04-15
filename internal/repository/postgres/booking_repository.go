package postgres

import (
	"context"
	"fmt"
	"time"

	"mediconnect/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) CountBookingsByDoctorAndDate(ctx context.Context, doctorID string, date time.Time) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("bookings").
		Where("doctor_id = ? AND schedule_date = ? AND status != ? AND status != ?", doctorID, date.Format("2006-01-02"), "CANCELLED", "NO_SHOW").
		Count(&count).Error
	return int(count), err
}

func (r *BookingRepository) GetBookedTimesForDoctorAndDate(ctx context.Context, doctorID string, date time.Time) ([]string, error) {
	var times []string
	err := r.db.WithContext(ctx).Table("bookings").
		Select("schedule_time").
		Where("doctor_id = ? AND schedule_date = ? AND status != ? AND status != ?", doctorID, date.Format("2006-01-02"), "CANCELLED", "NO_SHOW").
		Pluck("schedule_time", &times).Error

	// Format times specifically if needed, gorm returns what db gives
	var formattedTimes []string
	for _, t := range times {
		if len(t) >= 5 {
			formattedTimes = append(formattedTimes, t[:5])
		}
	}
	return formattedTimes, err
}

func (r *BookingRepository) CreateBookingWithLock(ctx context.Context, b *domain.Booking, maxQueue int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingTime string
		err := tx.Table("bookings").Select("schedule_time").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("doctor_id = ? AND schedule_date = ? AND schedule_time = ? AND status NOT IN ?",
				b.DoctorID, b.ScheduleDate.Format("2006-01-02"), b.ScheduleTime, []string{"CANCELLED", "NO_SHOW"}).
			Take(&existingTime).Error

		if err == nil {
			return fmt.Errorf("slot is already booked")
		} else if err != gorm.ErrRecordNotFound {
			return err // Other DB error
		}

		// Insert new booking
		if err := tx.Table("bookings").Create(b).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *BookingRepository) CancelStaleBookings(ctx context.Context) (int64, error) {
	res := r.db.WithContext(ctx).Table("bookings").
		Where("status = ? AND (schedule_date + schedule_time) < (CURRENT_TIMESTAMP - INTERVAL '30 minutes')", "PENDING").
		Updates(map[string]interface{}{"status": "NO_SHOW", "updated_at": gorm.Expr("CURRENT_TIMESTAMP")})

	return res.RowsAffected, res.Error
}

func (r *BookingRepository) GetBookingsByUserID(ctx context.Context, userID string) ([]domain.BookingDetail, error) {
	var results []domain.BookingDetail
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			b.id,
			b.booking_code,
			b.queue_number,
			b.status,
			b.schedule_date,
			b.schedule_time,
			b.facility_id,
			f.name  AS facility_name,
			f.type  AS facility_type,
			b.doctor_id,
			u.full_name  AS doctor_name,
			d.speciality
		FROM bookings b
		LEFT JOIN facilities f ON f.id = b.facility_id
		LEFT JOIN doctors    d ON d.id = b.doctor_id
		LEFT JOIN users      u ON u.id = d.user_id
		WHERE b.user_id = ?
		ORDER BY b.schedule_date DESC, b.schedule_time DESC
	`, userID).Scan(&results).Error
	return results, err
}
