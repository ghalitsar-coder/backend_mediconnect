package domain

import (
	"context"
	"time"
)

type Booking struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	FacilityID   string    `json:"facility_id"`
	DoctorID     string    `json:"doctor_id"`
	ScheduleDate time.Time `json:"schedule_date"`
	ScheduleTime string    `json:"schedule_time"`
	BookingCode  string    `json:"booking_code"`
	QueueNumber  string    `json:"queue_number"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
}

type BookingUsecase interface {
	CreateBooking(ctx context.Context, userID string, req BookingRequest) (*BookingResponse, error)
}
