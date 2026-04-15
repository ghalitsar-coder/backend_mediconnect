INSERT INTO facilities (name, address, lat, lng, type, district_id, is_active) VALUES
('Puskesmas Kebayoran Baru', 'Jl. Barito II No. 1, Kebayoran Baru', -6.242500, 106.797200, 'PUSKESMAS', '3174020', true),
('Puskesmas Tebet', 'Jl. Tebet Timur Dalam Raya', -6.230800, 106.852200, 'PUSKESMAS', '3174010', true),
('Klinik Sehat Bersama', 'Jl. Sudirman No. 10', -6.221500, 106.804800, 'KLINIK', '3174030', true),
('Klinik Medika Prima', 'Jl. MH Thamrin Kav 9', -6.195600, 106.822800, 'KLINIK', '3171010', true),
('Puskesmas Cilandak', 'Jl. Cilandak Tengah Raya', -6.291700, 106.797900, 'PUSKESMAS', '3174040', true);

INSERT INTO users (nik, email, password_hash, phone, full_name, role, is_active) VALUES
('1234567890123456', 'admin@mediconnect.id', '$2a$10$YrLqnCNNBt9L/WW76EHeeOdEN5y9vdZPhE.fmczTur1bM.sMEaQIi', '081234567890', 'Admin Mediconnect', 'DINKES', true),
('2345678901234567', 'nakes1@mediconnect.id', '$2a$10$YrLqnCNNBt9L/WW76EHeeOdEN5y9vdZPhE.fmczTur1bM.sMEaQIi', '082345678901', 'Dr. Budi Santoso', 'NAKES', true),
('3456789012345678', 'pasien1@mediconnect.id', '$2a$10$YrLqnCNNBt9L/WW76EHeeOdEN5y9vdZPhE.fmczTur1bM.sMEaQIi', '083456789012', 'Ahmad Dhani', 'PATIENT', true);

INSERT INTO doctors (user_id, facility_id, speciality, sip_number, is_active)
SELECT u.id, f.id, 'Umum', '12345/SIP/2026', true
FROM users u, facilities f
WHERE u.email = 'nakes1@mediconnect.id' AND f.name = 'Puskesmas Kebayoran Baru';
