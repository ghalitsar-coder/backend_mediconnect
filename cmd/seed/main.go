package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgres://mediconnect_user:mediconnect_password@localhost:5433/mediconnect_db?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	rand.Seed(time.Now().UnixNano())

	// Fetch all facilities
	type Facility struct {
		ID   string
		Name string
	}
	var facilities []Facility
	if err := db.Raw("SELECT id, name FROM facilities").Scan(&facilities).Error; err != nil {
		log.Fatalf("Failed to query facilities: %v", err)
	}

	if len(facilities) == 0 {
		log.Println("No facilities found!")
		return
	}

	polis := []string{
		"Poli Umum",
		"Poli Gigi",
		"Poli Anak",
		"Poli Kandungan",
		"Poli Paru",
		"Poli THT",
		"Poli Mata",
	}

	for _, fac := range facilities {
		for _, poli := range polis {
			newUserID := uuid.New().String()
			nik := fmt.Sprintf("32%014d", rand.Int63n(100000000000000))
			email := fmt.Sprintf("doctor_%s_%s@mediconnect.id", uuid.New().String()[:8], newUserID[:4])
			name := fmt.Sprintf("Dr. %s (%s)", poli, fac.Name)
			phone := fmt.Sprintf("0812%08d", rand.Intn(100000000))
			passwordHash := "$2a$10$5EuwgQ.y6L6T1Y1pW3c5MeHhZ.8.YnO8l5b0hA3q9.M26c6Zq8R/2"

			// Insert user
			if err := db.Exec(`
				INSERT INTO users (id, nik, email, password_hash, phone, full_name, role)
				VALUES (?, ?, ?, ?, ?, ?, 'NAKES')
			`, newUserID, nik, email, passwordHash, phone, name).Error; err != nil {
				log.Printf("Failed to insert user for %s: %v", name, err)
				continue
			}

			// Insert doctor
			newDocID := uuid.New().String()
			sip := fmt.Sprintf("SIP/%d/2026", rand.Intn(9999))

			if err := db.Exec(`
				INSERT INTO doctors (id, user_id, facility_id, speciality, sip_number)
				VALUES (?, ?, ?, ?, ?)
			`, newDocID, newUserID, fac.ID, poli, sip).Error; err != nil {
				log.Printf("Failed to insert doctor for %s: %v", name, err)
				continue
			}

			fmt.Printf("Injected doctor %s for facility %s\n", name, fac.Name)
		}
	}

	fmt.Println("Doctor injection completed!")
}
