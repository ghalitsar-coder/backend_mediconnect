package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ─── DSN ──────────────────────────────────────────────────────────────────────
const dsn = "postgres://mediconnect_user:mediconnect_password@localhost:5433/mediconnect_db?sslmode=disable"

// bcrypt hash for "password123" (cost=10) — same for all seed accounts
// Verified: bcrypt.CompareHashAndPassword(hash, []byte("password123")) == nil
const defaultPasswordHash = "$2a$10$YrLqnCNNBt9L/WW76EHeeOdEN5y9vdZPhE.fmczTur1bM.sMEaQIi"

// ─── Fixed deterministic UUIDs ────────────────────────────────────────────────
// Generated once; aman dipakai ulang (ON CONFLICT DO NOTHING)

var facilityIDs = []string{
	"11100000-0000-0000-0000-000000000001", // Puskesmas Kebayoran Baru
	"11100000-0000-0000-0000-000000000002", // Puskesmas Tebet
	"11100000-0000-0000-0000-000000000003", // Puskesmas Cilandak
	"11100000-0000-0000-0000-000000000004", // Puskesmas Mampang Prapatan
	"11100000-0000-0000-0000-000000000005", // Klinik Sehat Bersama
	"11100000-0000-0000-0000-000000000006", // Klinik Medika Prima
	"11100000-0000-0000-0000-000000000007", // Klinik Husada Jaya
	"11100000-0000-0000-0000-000000000008", // Klinik Permata Ibu
}

var nakesUserIDs = []string{
	"22200000-0000-0000-0000-000000000001",
	"22200000-0000-0000-0000-000000000002",
	"22200000-0000-0000-0000-000000000003",
	"22200000-0000-0000-0000-000000000004",
	"22200000-0000-0000-0000-000000000005",
	"22200000-0000-0000-0000-000000000006",
	"22200000-0000-0000-0000-000000000007",
	"22200000-0000-0000-0000-000000000008",
}

var patientUserIDs = []string{
	"33300000-0000-0000-0000-000000000001",
	"33300000-0000-0000-0000-000000000002",
	"33300000-0000-0000-0000-000000000003",
	"33300000-0000-0000-0000-000000000004",
	"33300000-0000-0000-0000-000000000005",
	"33300000-0000-0000-0000-000000000006",
	"33300000-0000-0000-0000-000000000007",
	"33300000-0000-0000-0000-000000000008",
	"33300000-0000-0000-0000-000000000009",
	"33300000-0000-0000-0000-000000000010",
}

const dinkesUserID = "44400000-0000-0000-0000-000000000001"

var doctorIDs = []string{
	"55500000-0000-0000-0000-000000000001",
	"55500000-0000-0000-0000-000000000002",
	"55500000-0000-0000-0000-000000000003",
	"55500000-0000-0000-0000-000000000004",
	"55500000-0000-0000-0000-000000000005",
	"55500000-0000-0000-0000-000000000006",
	"55500000-0000-0000-0000-000000000007",
	"55500000-0000-0000-0000-000000000008",
}

var appointmentIDs = []string{
	"66600000-0000-0000-0000-000000000001",
	"66600000-0000-0000-0000-000000000002",
	"66600000-0000-0000-0000-000000000003",
	"66600000-0000-0000-0000-000000000004",
	"66600000-0000-0000-0000-000000000005",
	"66600000-0000-0000-0000-000000000006",
	"66600000-0000-0000-0000-000000000007",
	"66600000-0000-0000-0000-000000000008",
	"66600000-0000-0000-0000-000000000009",
	"66600000-0000-0000-0000-000000000010",
}

var medicalRecordIDs = []string{
	"77700000-0000-0000-0000-000000000001",
	"77700000-0000-0000-0000-000000000002",
	"77700000-0000-0000-0000-000000000003",
	"77700000-0000-0000-0000-000000000004",
	"77700000-0000-0000-0000-000000000005",
}

func main() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("✅ Connected to database")

	seedFacilities(db)
	seedNakesUsers(db)
	seedPatientUsers(db)
	seedDinkesUser(db)
	seedDoctors(db)
	seedAppointments(db)
	seedBookings(db)
	seedMedicalRecords(db)
	seedLabResults(db)
	seedDiseaseLogs(db)

	fmt.Println("\n🎉 Seed completed successfully!")
}

// ─── FACILITIES ───────────────────────────────────────────────────────────────

func seedFacilities(db *gorm.DB) {
	type row struct {
		idx        int
		name       string
		address    string
		lat        float64
		lng        float64
		ftype      string
		districtID string
	}
	rows := []row{
		{0, "Puskesmas Kebayoran Baru", "Jl. Barito II No. 1, Kebayoran Baru", -6.242500, 106.797200, "PUSKESMAS", "3174020"},
		{1, "Puskesmas Tebet", "Jl. Tebet Timur Dalam Raya No. 3", -6.230800, 106.852200, "PUSKESMAS", "3174010"},
		{2, "Puskesmas Cilandak", "Jl. Cilandak Tengah Raya No. 12", -6.291700, 106.797900, "PUSKESMAS", "3174040"},
		{3, "Puskesmas Mampang Prapatan", "Jl. Mampang Prapatan Raya No. 18", -6.244100, 106.817500, "PUSKESMAS", "3174050"},
		{4, "Klinik Sehat Bersama", "Jl. Sudirman No. 10", -6.221500, 106.804800, "KLINIK", "3174030"},
		{5, "Klinik Medika Prima", "Jl. MH Thamrin Kav 9", -6.195600, 106.822800, "KLINIK", "3171010"},
		{6, "Klinik Husada Jaya", "Jl. Fatmawati No. 5, Cilandak", -6.302100, 106.795300, "KLINIK", "3174040"},
		{7, "Klinik Permata Ibu", "Jl. Raya Pasar Minggu No. 30", -6.258900, 106.843100, "KLINIK", "3174060"},
	}
	for _, r := range rows {
		if err := db.Exec(`
			INSERT INTO facilities (id, name, address, lat, lng, type, district_id, is_active)
			VALUES (?, ?, ?, ?, ?, ?, ?, true)
			ON CONFLICT (id) DO NOTHING
		`, facilityIDs[r.idx], r.name, r.address, r.lat, r.lng, r.ftype, r.districtID).Error; err != nil {
			log.Printf("⚠️  Skip facility %s: %v\n", r.name, err)
		}
	}
	log.Println("✅ Seeded facilities")
}

// ─── USERS: NAKES ─────────────────────────────────────────────────────────────

func seedNakesUsers(db *gorm.DB) {
	type row struct {
		idx   int
		nik   string
		email string
		phone string
		name  string
	}
	rows := []row{
		{0, "3271010101800001", "dr.budi@mediconnect.id", "081200000001", "Dr. Budi Santoso"},
		{1, "3271010101800002", "dr.siti@mediconnect.id", "081200000002", "Dr. Siti Rahayu"},
		{2, "3271010101800003", "dr.andi@mediconnect.id", "081200000003", "Dr. Andi Kurniawan"},
		{3, "3271010101800004", "dr.dewi@mediconnect.id", "081200000004", "Dr. Dewi Lestari"},
		{4, "3271010101800005", "dr.rudi@mediconnect.id", "081200000005", "Dr. Rudi Hermawan"},
		{5, "3271010101800006", "dr.ani@mediconnect.id", "081200000006", "Dr. Ani Wijaya"},
		{6, "3271010101800007", "dr.henry@mediconnect.id", "081200000007", "Dr. Henry Prasetyo"},
		{7, "3271010101800008", "dr.maya@mediconnect.id", "081200000008", "Dr. Maya Indah"},
	}
	for _, r := range rows {
		if err := db.Exec(`
			INSERT INTO users (id, nik, email, password_hash, phone, full_name, role, is_active)
			VALUES (?, ?, ?, ?, ?, ?, 'NAKES', true)
			ON CONFLICT (id) DO NOTHING
		`, nakesUserIDs[r.idx], r.nik, r.email, defaultPasswordHash, r.phone, r.name).Error; err != nil {
			log.Printf("⚠️  Skip nakes user %s: %v\n", r.name, err)
		}
	}
	log.Println("✅ Seeded NAKES users")
}

// ─── USERS: PATIENT ───────────────────────────────────────────────────────────

func seedPatientUsers(db *gorm.DB) {
	type row struct {
		idx   int
		nik   string
		email string
		phone string
		name  string
	}
	rows := []row{
		{0, "3271020202900001", "ahmad.dhani@gmail.com", "082100000001", "Ahmad Dhani"},
		{1, "3271020202900002", "bela.sari@gmail.com", "082100000002", "Bela Sari"},
		{2, "3271020202900003", "candra.putra@gmail.com", "082100000003", "Candra Putra"},
		{3, "3271020202900004", "dina.marlina@gmail.com", "082100000004", "Dina Marlina"},
		{4, "3271020202900005", "eko.prasetyo@gmail.com", "082100000005", "Eko Prasetyo"},
		{5, "3271020202900006", "fira.aulia@gmail.com", "082100000006", "Fira Aulia"},
		{6, "3271020202900007", "galih.wicaksono@gmail.com", "082100000007", "Galih Wicaksono"},
		{7, "3271020202900008", "hani.safitri@gmail.com", "082100000008", "Hani Safitri"},
		{8, "3271020202900009", "irwan.susanto@gmail.com", "082100000009", "Irwan Susanto"},
		{9, "3271020202900010", "julia.rahmawati@gmail.com", "082100000010", "Julia Rahmawati"},
	}
	for _, r := range rows {
		if err := db.Exec(`
			INSERT INTO users (id, nik, email, password_hash, phone, full_name, role, is_active)
			VALUES (?, ?, ?, ?, ?, ?, 'PATIENT', true)
			ON CONFLICT (id) DO NOTHING
		`, patientUserIDs[r.idx], r.nik, r.email, defaultPasswordHash, r.phone, r.name).Error; err != nil {
			log.Printf("⚠️  Skip patient user %s: %v\n", r.name, err)
		}
	}
	log.Println("✅ Seeded PATIENT users")
}

// ─── USER: DINKES ─────────────────────────────────────────────────────────────

func seedDinkesUser(db *gorm.DB) {
	if err := db.Exec(`
		INSERT INTO users (id, nik, email, password_hash, phone, full_name, role, is_active)
		VALUES (?, '3171030303800001', 'admin.dinkes@mediconnect.id', ?, '081300000001', 'Admin Dinas Kesehatan', 'DINKES', true)
		ON CONFLICT (id) DO NOTHING
	`, dinkesUserID, defaultPasswordHash).Error; err != nil {
		log.Printf("⚠️  Skip dinkes user: %v\n", err)
	}
	log.Println("✅ Seeded DINKES user")
}

// ─── DOCTORS ──────────────────────────────────────────────────────────────────

func seedDoctors(db *gorm.DB) {
	type row struct {
		docIdx      int
		userIdx     int
		facilityIdx int
		speciality  string
		sipNumber   string
	}
	rows := []row{
		{0, 0, 0, "Poli Umum", "SIP/0001/2026"},
		{1, 1, 0, "Poli Gigi", "SIP/0002/2026"},
		{2, 2, 1, "Poli Anak", "SIP/0003/2026"},
		{3, 3, 1, "Poli Kandungan", "SIP/0004/2026"},
		{4, 4, 2, "Poli Paru", "SIP/0005/2026"},
		{5, 5, 2, "Poli THT", "SIP/0006/2026"},
		{6, 6, 3, "Poli Mata", "SIP/0007/2026"},
		{7, 7, 4, "Poli Umum", "SIP/0008/2026"},
	}
	for _, r := range rows {
		if err := db.Exec(`
			INSERT INTO doctors (id, user_id, facility_id, speciality, sip_number, is_active)
			VALUES (?, ?, ?, ?, ?, true)
			ON CONFLICT (id) DO NOTHING
		`, doctorIDs[r.docIdx], nakesUserIDs[r.userIdx], facilityIDs[r.facilityIdx], r.speciality, r.sipNumber).Error; err != nil {
			log.Printf("⚠️  Skip doctor idx %d: %v\n", r.docIdx, err)
		}
	}
	log.Println("✅ Seeded doctors")
}

// ─── APPOINTMENTS ─────────────────────────────────────────────────────────────

func seedAppointments(db *gorm.DB) {
	type row struct {
		apptIdx     int
		userIdx     int
		facilityIdx int
		doctorIdx   int
		poli        string
		daysOffset  int
		status      string
	}
	rows := []row{
		{0, 0, 0, 0, "Poli Umum", -5, "DONE"},
		{1, 1, 0, 0, "Poli Umum", -3, "DONE"},
		{2, 2, 1, 2, "Poli Anak", -2, "CONFIRMED"},
		{3, 3, 1, 3, "Poli Kandungan", -1, "CONFIRMED"},
		{4, 4, 2, 4, "Poli Paru", 1, "PENDING"},
		{5, 5, 2, 5, "Poli THT", 2, "PENDING"},
		{6, 6, 3, 6, "Poli Mata", 3, "PENDING"},
		{7, 7, 4, 7, "Poli Umum", -10, "CANCELLED"},
		{8, 8, 0, 0, "Poli Umum", -7, "DONE"},
		{9, 9, 1, 2, "Poli Anak", 5, "PENDING"},
	}

	now := time.Now()
	for _, r := range rows {
		scheduled := now.AddDate(0, 0, r.daysOffset).Truncate(24*time.Hour).Add(9 * time.Hour)
		qrToken := fmt.Sprintf("QR-%d-%d", r.apptIdx+1, time.Now().UnixNano())

		var cancelledAt interface{} = nil
		var cancelReason interface{} = nil
		if r.status == "CANCELLED" {
			t := scheduled.Add(1 * time.Hour)
			cancelledAt = t
			cancelReason = "Pasien membatalkan janji"
		}

		if err := db.Exec(`
			INSERT INTO appointments
				(id, user_id, facility_id, doctor_id, poli, scheduled_at, status, qr_token, cancelled_at, cancel_reason)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING
		`,
			appointmentIDs[r.apptIdx],
			patientUserIDs[r.userIdx],
			facilityIDs[r.facilityIdx],
			doctorIDs[r.doctorIdx],
			r.poli,
			scheduled,
			r.status,
			qrToken,
			cancelledAt,
			cancelReason,
		).Error; err != nil {
			log.Printf("⚠️  Skip appointment idx %d: %v\n", r.apptIdx, err)
		}
	}
	log.Println("✅ Seeded appointments")
}

// ─── BOOKINGS ─────────────────────────────────────────────────────────────────

var bookingIDs = []string{
	"88800000-0000-0000-0000-000000000001",
	"88800000-0000-0000-0000-000000000002",
	"88800000-0000-0000-0000-000000000003",
	"88800000-0000-0000-0000-000000000004",
	"88800000-0000-0000-0000-000000000005",
}

func seedBookings(db *gorm.DB) {
	// Cek apakah tabel bookings ada
	var exists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'bookings')").Scan(&exists)
	if !exists {
		log.Println("ℹ️  Tabel bookings belum ada, skip seed bookings")
		return
	}

	type row struct {
		bookingIdx  int
		userIdx     int
		facilityIdx int
		doctorIdx   int
		daysOffset  int
		schedTime   string
		bookingCode string
		queueNumber string
		status      string
	}

	now := time.Now()
	rows := []row{
		{0, 0, 0, 0, -2, "08:00", "BKG-20260413-001", "A001", "DONE"},
		{1, 1, 0, 0, -1, "09:00", "BKG-20260414-002", "A002", "DONE"},
		{2, 2, 1, 2, 1, "10:00", "BKG-20260416-003", "B001", "PENDING"},
		{3, 3, 1, 3, 2, "11:00", "BKG-20260417-004", "C001", "PENDING"},
		{4, 4, 2, 4, 3, "08:30", "BKG-20260418-005", "D001", "PENDING"},
	}

	for _, r := range rows {
		schedDate := now.AddDate(0, 0, r.daysOffset).Truncate(24 * time.Hour)
		if err := db.Exec(`
			INSERT INTO bookings
				(id, user_id, facility_id, doctor_id, schedule_date, schedule_time,
				 booking_code, queue_number, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`,
			bookingIDs[r.bookingIdx],
			patientUserIDs[r.userIdx],
			facilityIDs[r.facilityIdx],
			doctorIDs[r.doctorIdx],
			schedDate,
			r.schedTime,
			r.bookingCode,
			r.queueNumber,
			r.status,
		).Error; err != nil {
			log.Printf("⚠️  Skip booking idx %d: %v\n", r.bookingIdx, err)
		}
	}
	log.Println("✅ Seeded bookings")
}

// ─── MEDICAL RECORDS ──────────────────────────────────────────────────────────

func seedMedicalRecords(db *gorm.DB) {
	type row struct {
		mrIdx        int
		userIdx      int
		apptIdx      int
		createdByIdx int
		diagnosisEnc string
		notesEnc     string
		icd10Code    string
	}
	rows := []row{
		{0, 0, 0, 0, "ENC::Hipertensi derajat 1", "ENC::Tekanan darah 150/90, disarankan diet rendah garam", "I10"},
		{1, 1, 1, 0, "ENC::Infeksi saluran pernapasan atas", "ENC::Pasien batuk pilek 3 hari, diberikan antibiotik", "J06"},
		{2, 2, 2, 2, "ENC::Diare akut tanpa dehidrasi", "ENC::BAB cair 5x sehari, ORS diberikan", "A09"},
		{3, 3, 3, 3, "ENC::Anemia defisiensi besi pada kehamilan", "ENC::Hb 9.5, suplemen zat besi diberikan", "O99"},
		{4, 8, 8, 0, "ENC::Dispepsia fungsional", "ENC::Nyeri ulu hati, antasida dan PPI diberikan", "K30"},
	}
	for _, r := range rows {
		if err := db.Exec(`
			INSERT INTO medical_records
				(id, user_id, appointment_id, diagnosis_enc, notes_enc, icd10_code, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING
		`,
			medicalRecordIDs[r.mrIdx],
			patientUserIDs[r.userIdx],
			appointmentIDs[r.apptIdx],
			r.diagnosisEnc,
			r.notesEnc,
			r.icd10Code,
			nakesUserIDs[r.createdByIdx],
		).Error; err != nil {
			log.Printf("⚠️  Skip medical record idx %d: %v\n", r.mrIdx, err)
		}
	}
	log.Println("✅ Seeded medical records")
}

// ─── LAB RESULTS ──────────────────────────────────────────────────────────────

var labResultIDs = []string{
	"99900000-0000-0000-0000-000000000001",
	"99900000-0000-0000-0000-000000000002",
	"99900000-0000-0000-0000-000000000003",
	"99900000-0000-0000-0000-000000000004",
	"99900000-0000-0000-0000-000000000005",
	"99900000-0000-0000-0000-000000000006",
	"99900000-0000-0000-0000-000000000007",
	"99900000-0000-0000-0000-000000000008",
	"99900000-0000-0000-0000-000000000009",
	"99900000-0000-0000-0000-000000000010",
}

func seedLabResults(db *gorm.DB) {
	type row struct {
		labIdx      int
		recordIdx   int
		testName    string
		resultEnc   string
		unit        string
		normalRange string
		isReady     bool
	}
	rows := []row{
		{0, 0, "Tekanan Darah Sistolik", "ENC::150", "mmHg", "90-120", true},
		{1, 0, "Tekanan Darah Diastolik", "ENC::90", "mmHg", "60-80", true},
		{2, 0, "Kolesterol Total", "ENC::220", "mg/dL", "<200", true},
		{3, 1, "LED (Laju Endap Darah)", "ENC::35", "mm/jam", "<20", true},
		{4, 1, "Leukosit", "ENC::11500", "/µL", "4000-10000", true},
		{5, 2, "Natrium", "ENC::138", "mEq/L", "136-145", true},
		{6, 2, "Kalium", "ENC::3.8", "mEq/L", "3.5-5.0", true},
		{7, 3, "Hemoglobin", "ENC::9.5", "g/dL", "12-16", true},
		{8, 3, "Serum Ferritin", "ENC::8", "ng/mL", "12-150", true},
		{9, 4, "Amilase", "ENC::95", "U/L", "28-100", false},
	}
	readyAt := time.Now().Add(-2 * time.Hour)
	for _, r := range rows {
		var ra interface{} = nil
		if r.isReady {
			ra = readyAt
		}
		if err := db.Exec(`
			INSERT INTO lab_results
				(id, record_id, test_name, result_enc, unit, normal_range, is_ready, ready_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING
		`,
			labResultIDs[r.labIdx],
			medicalRecordIDs[r.recordIdx],
			r.testName, r.resultEnc, r.unit, r.normalRange, r.isReady, ra,
		).Error; err != nil {
			log.Printf("⚠️  Skip lab result idx %d: %v\n", r.labIdx, err)
		}
	}
	log.Println("✅ Seeded lab results")
}

// ─── DISEASE LOGS ─────────────────────────────────────────────────────────────

var diseaseLogIDs = []string{
	"aaaa0000-0000-0000-0000-000000000001",
	"aaaa0000-0000-0000-0000-000000000002",
	"aaaa0000-0000-0000-0000-000000000003",
	"aaaa0000-0000-0000-0000-000000000004",
	"aaaa0000-0000-0000-0000-000000000005",
	"aaaa0000-0000-0000-0000-000000000006",
	"aaaa0000-0000-0000-0000-000000000007",
	"aaaa0000-0000-0000-0000-000000000008",
}

func seedDiseaseLogs(db *gorm.DB) {
	type row struct {
		logIdx      int
		recordIdx   int
		icd10Code   string
		facilityIdx int
		districtID  string
	}
	rows := []row{
		{0, 0, "I10", 0, "3174020"},
		{1, 1, "J06", 0, "3174020"},
		{2, 2, "A09", 1, "3174010"},
		{3, 3, "O99", 1, "3174010"},
		{4, 4, "K30", 2, "3174040"},
		{5, 0, "I10", 2, "3174040"},
		{6, 1, "J06", 3, "3174050"},
		{7, 2, "A09", 4, "3174030"},
	}
	for _, r := range rows {
		if err := db.Exec(`
			INSERT INTO disease_logs
				(id, record_id, icd10_code, facility_id, district_id)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING
		`,
			diseaseLogIDs[r.logIdx],
			medicalRecordIDs[r.recordIdx],
			r.icd10Code,
			facilityIDs[r.facilityIdx],
			r.districtID,
		).Error; err != nil {
			log.Printf("⚠️  Skip disease log idx %d: %v\n", r.logIdx, err)
		}
	}
	log.Println("✅ Seeded disease logs")
}
