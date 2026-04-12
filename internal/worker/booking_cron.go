package worker

import (
	"context"
	"fmt"
	"time"

	"mediconnect/internal/repository/postgres"
)

type BookingCron struct {
	repo *postgres.BookingRepository
}

func NewBookingCron(repo *postgres.BookingRepository) *BookingCron {
	return &BookingCron{repo: repo}
}

// StartAutoCancellation runs a ticker every 5 minutes to cancel stale bookings
func (c *BookingCron) StartAutoCancellation(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				cancelledCount, err := c.repo.CancelStaleBookings(context.Background())
				if err != nil {
					fmt.Println("Cron Error: Failed to auto-cancel bookings", err)
				} else if cancelledCount > 0 {
					fmt.Printf("Cron: Successfully auto-cancelled %d stale NO_SHOW bookings\n", cancelledCount)
				}
			}
		}
	}()
}
