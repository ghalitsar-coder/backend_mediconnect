package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"mediconnect/internal/domain"
	"mediconnect/internal/repository/postgres"
	"mediconnect/pkg/messaging"
)

type BookingUsecase struct {
	repo   *postgres.BookingRepository // Using the struct directly for simplicity here
	rabbit *messaging.RabbitMQ
}

func NewBookingUsecase(repo *postgres.BookingRepository, rabbit *messaging.RabbitMQ) *BookingUsecase {
	return &BookingUsecase{repo: repo, rabbit: rabbit}
}

func (u *BookingUsecase) CreateBooking(ctx context.Context, userID string, req domain.BookingRequest) (*domain.BookingResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule_date format")
	}

	queueRand := rand.Intn(100) + 1
	b := &domain.Booking{
		ID:           uuid.New().String(),
		UserID:       userID,
		FacilityID:   req.FacilityID,
		DoctorID:     req.DoctorID,
		ScheduleDate: parsedDate,
		ScheduleTime: req.ScheduleTime,
		BookingCode:  fmt.Sprintf("MC-%X", rand.Uint32()),
		QueueNumber:  fmt.Sprintf("A-%d", queueRand),
		Status:       "PENDING",
	}

	// Double Booking Check inside the repository transaction (concurrency lock)
	err = u.repo.CreateBookingWithLock(ctx, b, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %v", err)
	}

	// Generate response
	res := &domain.BookingResponse{
		BookingID:   b.ID,
		Token:       b.BookingCode,
		QueueNumber: b.QueueNumber,
	}

	// RabbitMQ Notification Offloading
	if u.rabbit != nil {
		msgBytes, _ := json.Marshal(map[string]interface{}{
			"user_id":    userID,
			"booking_id": b.ID,
			"booking_cd": b.BookingCode,
			"queue":      b.QueueNumber,
			"action":     "SEND_REMINDER",
			"sch_date":   req.ScheduleDate,
			"sch_time":   req.ScheduleTime,
		})

		err := u.rabbit.Channel.PublishWithContext(ctx,
			"notifications_exchange",
			"booking_routing_key",
			false, false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        msgBytes,
			},
		)
		if err != nil {
			fmt.Println("Warning: Failed to publish RabbitMQ message:", err)
		}
	}

	return res, nil
}
