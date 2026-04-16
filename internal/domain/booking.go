package domain

import (
	"context"
	"time"
)

type Booking struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	UserID       string    `json:"user_id"`
	FacilityID   string    `json:"facility_id"`
	DoctorID     string    `json:"doctor_id"`
	ScheduleDate time.Time `json:"schedule_date"`
	ScheduleTime string    `json:"schedule_time"`
	BookingCode  string    `json:"booking_code"`
	QueueNumber  string    `json:"queue_number"`
	Status       string    `json:"status" gorm:"default:PENDING"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// BookingDetail is a richer view of a booking with facility and doctor names joined in.
type BookingDetail struct {
	ID           string    `json:"id"`
	BookingCode  string    `json:"booking_code"`
	QueueNumber  string    `json:"queue_number"`
	Status       string    `json:"status"`
	ScheduleDate time.Time `json:"schedule_date"`
	ScheduleTime string    `json:"schedule_time"`
	FacilityID   string    `json:"facility_id"`
	FacilityName string    `json:"facility_name"`
	FacilityType string    `json:"facility_type"`
	DoctorID     string    `json:"doctor_id"`
	DoctorName   string    `json:"doctor_name"`
	Speciality   string    `json:"speciality"`
	CreatedAt    time.Time `json:"created_at"`
}

type BookingRequest struct {
	FacilityID   string `json:"facility_id"`
	DoctorID     string `json:"doctor_id"`
	ScheduleDate string `json:"schedule_date"`
	ScheduleTime string `json:"schedule_time"`
}

type BookingResponse struct {
	BookingID   string `json:"booking_id"`
	Token       string `json:"token"`
	QueueNumber string `json:"no_antrian"`
}

type BookingRepository interface {
	CountBookingsByDoctorAndDate(ctx context.Context, doctorID string, date time.Time) (int, error)
	GetBookedTimesForDoctorAndDate(ctx context.Context, doctorID string, date time.Time) ([]string, error)
	CreateBookingWithLock(ctx context.Context, b *Booking, maxQueue int) error
	CancelStaleBookings(ctx context.Context) (int64, error)
	GetBookingsByUserID(ctx context.Context, userID string) ([]BookingDetail, error)
}

type BookingUsecase interface {
	CreateBooking(ctx context.Context, userID string, req BookingRequest) (*BookingResponse, error)
	GetMyBookings(ctx context.Context, userID string) ([]BookingDetail, error)
}
