-- 002_booking_schema.sql
-- Create Doctor, Schedule, and Booking tables

CREATE TABLE IF NOT EXISTS doctors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID REFERENCES facilities(id),
    name VARCHAR(150) NOT NULL,
    specialization VARCHAR(100),
    poli_name VARCHAR(100) NOT NULL,
    rating DECIMAL(2, 1) DEFAULT 0.0,
    patients_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS doctor_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    doctor_id UUID REFERENCES doctors(id),
    day_of_week INT NOT NULL, -- 0 (Minggu) sampai 6 (Sabtu)
    start_time TIME NOT NULL, 
    end_time TIME NOT NULL,
    slot_duration_minutes INT DEFAULT 30,
    max_patients INT NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    facility_id UUID REFERENCES facilities(id),
    doctor_id UUID REFERENCES doctors(id),
    
    schedule_date DATE NOT NULL,
    schedule_time TIME NOT NULL,
    
    booking_code VARCHAR(20) UNIQUE NOT NULL,
    queue_number VARCHAR(10),
    
    status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, CONFIRMED, COMPLETED, CANCELLED, NO_SHOW
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing untuk optimasi filter
CREATE INDEX IF NOT EXISTS idx_bookings_doctor_date ON bookings(doctor_id, schedule_date);
CREATE INDEX IF NOT EXISTS idx_doctors_facility_poli ON doctors(facility_id, poli_name);
