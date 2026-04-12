-- ============================================================
--  MediConnect ID — Database Schema (PostgreSQL 16)
--  PRD v2.0 | Target Rilis Q3 2026
--  Migration: 001_init_schema
-- ============================================================

-- Extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";  -- untuk gen_random_uuid()

-- ============================================================
-- 1. ENUM TYPES
-- ============================================================

CREATE TYPE user_role AS ENUM ('PATIENT', 'NAKES', 'DINKES');

CREATE TYPE appointment_status AS ENUM (
    'PENDING',
    'CONFIRMED',
    'DONE',
    'CANCELLED'
);

CREATE TYPE facility_type AS ENUM ('PUSKESMAS', 'KLINIK');

-- ============================================================
-- 2. TABEL: users
-- ============================================================

CREATE TABLE users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    nik           VARCHAR(16) NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone         VARCHAR(20),
    full_name     VARCHAR(255) NOT NULL,
    role          user_role   NOT NULL DEFAULT 'PATIENT',
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_nik   ON users (nik);
CREATE INDEX idx_users_email ON users (email);

COMMENT ON TABLE  users              IS 'Pengguna sistem: pasien, tenaga kesehatan, dan Dinkes';
COMMENT ON COLUMN users.nik          IS 'Nomor Induk Kependudukan (16 digit), divalidasi via API Dukcapil';
COMMENT ON COLUMN users.role         IS 'PATIENT | NAKES | DINKES';

-- ============================================================
-- 3. TABEL: facilities
-- ============================================================

CREATE TABLE facilities (
    id          UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255)    NOT NULL,
    address     TEXT            NOT NULL,
    lat         DECIMAL(10, 7)  NOT NULL,
    lng         DECIMAL(10, 7)  NOT NULL,
    type        facility_type   NOT NULL,
    district_id VARCHAR(10)     NOT NULL,  -- kode wilayah BPS
    is_active   BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_facilities_district ON facilities (district_id);
CREATE INDEX idx_facilities_type     ON facilities (type);

COMMENT ON TABLE  facilities             IS 'Fasilitas kesehatan: Puskesmas dan Klinik';
COMMENT ON COLUMN facilities.lat         IS 'Latitude koordinat GPS';
COMMENT ON COLUMN facilities.lng         IS 'Longitude koordinat GPS';
COMMENT ON COLUMN facilities.district_id IS 'Kode wilayah BPS (kecamatan/kelurahan)';

-- ============================================================
-- 4. TABEL: doctors
-- ============================================================

CREATE TABLE doctors (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    facility_id UUID        NOT NULL REFERENCES facilities(id) ON DELETE CASCADE,
    speciality  VARCHAR(100),
    sip_number  VARCHAR(50),   -- Surat Izin Praktek
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_doctors_user_id     ON doctors (user_id);
CREATE INDEX idx_doctors_facility_id ON doctors (facility_id);

-- ============================================================
-- 5. TABEL: appointments
-- ============================================================

CREATE TABLE appointments (
    id            UUID               PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID               NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    facility_id   UUID               NOT NULL REFERENCES facilities(id),
    doctor_id     UUID               REFERENCES doctors(id),
    poli          VARCHAR(100)       NOT NULL,
    scheduled_at  TIMESTAMPTZ        NOT NULL,
    status        appointment_status NOT NULL DEFAULT 'PENDING',
    qr_token      VARCHAR(255)       NOT NULL UNIQUE,
    checked_in_at TIMESTAMPTZ,
    cancelled_at  TIMESTAMPTZ,
    cancel_reason VARCHAR(255),
    created_at    TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_appointments_user_id      ON appointments (user_id);
CREATE INDEX idx_appointments_scheduled_at ON appointments (scheduled_at);
CREATE INDEX idx_appointments_status       ON appointments (status);
CREATE INDEX idx_appointments_facility_id  ON appointments (facility_id);

COMMENT ON TABLE  appointments          IS 'Booking antrian pasien ke fasilitas kesehatan';
COMMENT ON COLUMN appointments.qr_token IS 'Token unik untuk QR Code check-in di lokasi';
COMMENT ON COLUMN appointments.poli     IS 'Nama poli/unit layanan yang dituju';

-- ============================================================
-- 6. TABEL: medical_records
-- ============================================================

CREATE TABLE medical_records (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    appointment_id UUID        REFERENCES appointments(id),
    diagnosis_enc  TEXT        NOT NULL,   -- AES-256 encrypted
    notes_enc      TEXT,                   -- AES-256 encrypted
    icd10_code     VARCHAR(10),
    created_by     UUID        NOT NULL REFERENCES users(id),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_medical_records_user_id    ON medical_records (user_id);
CREATE INDEX idx_medical_records_created_at ON medical_records (created_at);

COMMENT ON TABLE  medical_records               IS 'Rekam medis digital per kunjungan';
COMMENT ON COLUMN medical_records.diagnosis_enc IS 'Diagnosis terenkripsi AES-256';
COMMENT ON COLUMN medical_records.notes_enc     IS 'Catatan dokter terenkripsi AES-256';
COMMENT ON COLUMN medical_records.icd10_code    IS 'Kode ICD-10 untuk klasifikasi penyakit';
COMMENT ON COLUMN medical_records.created_by    IS 'FK ke users.id (tenaga kesehatan yang input)';

-- ============================================================
-- 7. TABEL: lab_results
-- ============================================================

CREATE TABLE lab_results (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id    UUID         NOT NULL REFERENCES medical_records(id) ON DELETE CASCADE,
    test_name    VARCHAR(255) NOT NULL,
    result_enc   TEXT         NOT NULL,   -- AES-256 encrypted
    unit         VARCHAR(50),
    normal_range VARCHAR(100),
    is_ready     BOOLEAN      NOT NULL DEFAULT FALSE,
    ready_at     TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_lab_results_record_id ON lab_results (record_id);
CREATE INDEX idx_lab_results_is_ready  ON lab_results (is_ready);

COMMENT ON COLUMN lab_results.result_enc IS 'Hasil lab terenkripsi AES-256';
COMMENT ON COLUMN lab_results.is_ready   IS 'TRUE saat hasil siap; trigger event lab.result.ready ke RabbitMQ';

-- ============================================================
-- 8. TABEL: disease_logs
-- ============================================================

CREATE TABLE disease_logs (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id   UUID        NOT NULL REFERENCES medical_records(id) ON DELETE CASCADE,
    icd10_code  VARCHAR(10) NOT NULL,
    facility_id UUID        NOT NULL REFERENCES facilities(id),
    district_id VARCHAR(10) NOT NULL,
    logged_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_disease_logs_icd10_code  ON disease_logs (icd10_code);
CREATE INDEX idx_disease_logs_logged_at   ON disease_logs (logged_at);
CREATE INDEX idx_disease_logs_district_id ON disease_logs (district_id);

COMMENT ON TABLE  disease_logs            IS 'Source data untuk surveilans epidemiologi Dinkes';
COMMENT ON COLUMN disease_logs.icd10_code IS 'Kode penyakit ICD-10';
COMMENT ON COLUMN disease_logs.district_id IS 'Kode wilayah untuk agregasi heatmap';

-- ============================================================
-- 9. TABEL: audit_logs (partisi per bulan)
-- ============================================================

CREATE TABLE audit_logs (
    id          BIGSERIAL,
    user_id     UUID         REFERENCES users(id),
    action      VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id   UUID,
    ip_addr     INET,
    user_agent  TEXT,
    payload     JSONB,
    timestamp   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
) PARTITION BY RANGE (timestamp);

-- Partisi per bulan (contoh: 2 bulan pertama launch)
CREATE TABLE audit_logs_2026_07 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');

CREATE TABLE audit_logs_2026_08 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');

CREATE INDEX idx_audit_logs_user_id   ON audit_logs (user_id);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs (timestamp);

COMMENT ON TABLE  audit_logs         IS 'Immutable audit trail. Insert-only. Partisi by month.';
COMMENT ON COLUMN audit_logs.action  IS 'Contoh: CREATE_RECORD, UPDATE_RECORD, DELETE_APPOINTMENT';
COMMENT ON COLUMN audit_logs.payload IS 'Snapshot data sebelum/sesudah perubahan (JSON)';

-- ============================================================
-- 10. TRIGGER: auto-update updated_at
-- ============================================================

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_facilities_updated_at
    BEFORE UPDATE ON facilities
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_appointments_updated_at
    BEFORE UPDATE ON appointments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_medical_records_updated_at
    BEFORE UPDATE ON medical_records
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ============================================================
-- END OF SCHEMA — 001_init_schema
-- ============================================================
