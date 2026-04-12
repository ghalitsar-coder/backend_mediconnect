package domain

import (
	"context"
	"time"
)

type Doctor struct {
	ID             string  `json:"id"`
	FacilityID     string  `json:"facility_id"`
	Name           string  `json:"name"`
	Specialization string  `json:"spec"`
	PoliName       string  `json:"poli"`
	Rating         float64 `json:"rating"`
	PatientsCount  int     `json:"patients"`
}

type DoctorSchedule struct {
	ID                  string `json:"id"`
	DoctorID            string `json:"doctor_id"`
	DayOfWeek           int    `json:"day_of_week"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	SlotDurationMinutes int    `json:"slot_duration_minutes"`
	MaxPatients         int    `json:"max_patients"`
}

type Slot struct {
	Time        string `json:"time"`
	IsAvailable bool   `json:"isAvailable"`
}

type DoctorRepository interface {
	GetDoctors(ctx context.Context, facilityID string, poliName string) ([]Doctor, error)
	GetDoctorSchedules(ctx context.Context, doctorID string, dayOfWeek int) ([]DoctorSchedule, error)
}

type DoctorUsecase interface {
	GetDoctors(ctx context.Context, facilityID string, poliName string) ([]Doctor, error)
	GetAvailableSlots(ctx context.Context, doctorID string, date time.Time) ([]Slot, error)
}
