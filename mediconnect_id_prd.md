# MediConnect ID
**Product Requirements Document · Platform Kesehatan Terpadu**

| Versi | Status | Target Rilis | Tahap |
|-------|--------|--------------|-------|
| 2.0 | Draft | Q3 2026 | MVP |

**Tech Stack:** `Next.js 15` `Go 1.23` `PostgreSQL 16` `Redis 7` `RabbitMQ 3.13` `Kubernetes 1.30` `Jenkins` `Docker`

---

## Daftar Isi

1. [Ringkasan Proyek](#1--ringkasan-proyek)
2. [Persona & User Stories](#2--persona--user-stories)
3. [Fitur & Acceptance Criteria](#3--fitur--acceptance-criteria)
4. [Arsitektur & Alur Data](#4--arsitektur--alur-data)
5. [Skema Database](#5--skema-database-postgresql)
6. [API Contract](#6--api-contract-rest-v1)
7. [Kebutuhan Non-Fungsional](#7--kebutuhan-non-fungsional)
8. [CI/CD Pipeline](#8--cicd-pipeline-jenkins)
9. [Risiko & Mitigasi](#9--risiko--mitigasi)
10. [Success Metrics & KPI](#10--success-metrics--kpi)
11. [Open Questions](#11--open-questions)

---

## 1 · Ringkasan Proyek

**MediConnect ID** adalah platform kesehatan publik terpadu yang mendigitalisasi layanan Puskesmas dan klinik masyarakat di Indonesia. Sistem ini menyediakan tiga kapabilitas inti: manajemen antrian berbasis online, rekam medis digital terpusat, dan surveilans epidemiologi real-time bagi Dinas Kesehatan (Dinkes).

Arsitektur yang dipilih adalah **microservices event-driven** dengan Redis sebagai lapisan caching, RabbitMQ sebagai message broker, dan Kubernetes sebagai platform orkestrasi. Pendekatan ini memastikan skalabilitas horizontal saat terjadi lonjakan pendaftaran dan keandalan pengiriman notifikasi.

> **Out of scope MVP:** integrasi BPJS Kesehatan, telemedicine, pembayaran online, dan manajemen inventaris farmasi. Fitur ini masuk roadmap v2.

---

## 2 · Persona & User Stories

### 👤 Warga (Pasien)
Masyarakat umum berusia 18–65 tahun. Akses via smartphone, literasi digital menengah.

**Goal:** Daftar antrian tanpa antre fisik, lihat riwayat periksa kapan saja.

`Mobile-first`

### 🩺 Tenaga Kesehatan
Admin Puskesmas, perawat, atau dokter. Mengelola antrian harian dan rekam medis.

**Goal:** Efisiensi pengelolaan pasien, input diagnosis cepat, pantau antrian real-time.

`Desktop`

### 🏛️ Pemerintah (Dinkes)
Pengambil kebijakan Dinas Kesehatan kota/provinsi. Akses dashboard agregat.

**Goal:** Deteksi dini wabah, alokasi sumber daya tepat waktu berbasis data.

`Read-only`

---

## 3 · Fitur & Acceptance Criteria

### F-01 · Sistem Antrian Smart (Booking) — `P0 · MVP`

> *"Sebagai warga, saya ingin mendaftar antrian puskesmas secara online dari rumah, agar saya tidak perlu datang pagi-pagi untuk antre."*

- Warga dapat memilih fasilitas, poli, dan dokter dari daftar yang tersedia.
- Validasi NIK/KTP via API Dukcapil sebelum pendaftaran dikonfirmasi.
- Slot antrian live disimpan di Redis (TTL 24 jam) untuk pengecekan < 50ms.
- Sistem menghasilkan QR Code unik untuk check-in di lokasi.
- Batas maksimal 2 booking aktif per NIK dalam satu hari.
- Pembatalan antrian otomatis jika tidak ada check-in 30 menit setelah jadwal.

**Tech:** `Redis` `PostgreSQL` `RabbitMQ`

---

### F-02 · Rekam Medis Digital (EMR) — `P0 · MVP`

> *"Sebagai tenaga kesehatan, saya ingin mengakses dan memperbarui rekam medis pasien dengan cepat, agar pelayanan lebih efisien."*

- Tenaga kesehatan dapat membuat, membaca, dan memperbarui rekam medis per kunjungan.
- Data sensitif (diagnosis, obat) dienkripsi AES-256 di level database.
- Warga dapat melihat riwayat medis miliknya sendiri (read-only).
- Setiap perubahan rekam medis dicatat di audit log dengan timestamp dan user ID.
- Hasil laboratorium yang sudah selesai men-trigger event `lab.result.ready` ke RabbitMQ.

**Tech:** `AES-256` `JWT` `Audit Log`

---

### F-03 · Dashboard Surveilans Penyakit — `P1`

> *"Sebagai pejabat Dinkes, saya ingin melihat peta penyebaran penyakit secara real-time, agar saya bisa mengambil keputusan intervensi lebih cepat."*

- Visualisasi heatmap interaktif per kecamatan/kelurahan menggunakan Next.js + library peta.
- Filter berdasarkan jenis penyakit, rentang tanggal, dan level administratif.
- Statistik agregat harian di-cache di Redis (refresh setiap jam), mencegah overload DB saat akses bersamaan.
- Alert otomatis jika kasus penyakit tertentu naik > 20% dalam 48 jam di satu wilayah.
- Export data ke CSV/Excel untuk pelaporan internal Dinkes.

**Tech:** `Next.js` `Redis cache` `Analytics Service`

---

### F-04 · Sistem Notifikasi Asinkron — `P1`

> *"Sebagai warga, saya ingin menerima pengingat jadwal dan notifikasi hasil lab via WhatsApp atau email, tanpa harus mengecek aplikasi terus."*

- Notifikasi dikirim via email dan WhatsApp (Twilio/WA Business API).
- Producer (Booking/EMR Service) mengirim pesan ke RabbitMQ exchange.
- Notification Worker mengonsumsi pesan secara asinkron — tidak memblokir service utama.
- Retry mechanism: 3 kali percobaan ulang dengan exponential backoff jika pengiriman gagal.
- Pesan di queue bersifat persistent (durable) — tidak hilang jika broker restart.

**Tech:** `RabbitMQ` `Twilio` `SMTP` `Dead Letter Queue`

---

## 4 · Arsitektur & Alur Data

### Event-Driven Flow — Booking Appointment

1. **Client request** — Next.js frontend → API Gateway (JWT auth) → Booking Service (Go)
2. **Slot reservation (Redis)** — Booking Service cek ketersediaan slot di Redis. Jika ada, lakukan atomic decrement counter slot. Jika tidak ada, tolak dengan 409 Conflict.
3. **Persist ke PostgreSQL** — Record appointment ditulis ke tabel `appointments` dengan status `CONFIRMED`.
4. **Publish event ke RabbitMQ** — Event `appointment.created` dipublish ke exchange. Routing key menentukan queue tujuan.
5. **Consumer processing (async)** — Notification Worker → kirim WhatsApp/email. Analytics Worker → update statistik penyakit di cache Redis.
6. **Response ke client** — API Gateway mengembalikan booking ID + QR Code payload. Response time target: < 200ms.

### Strategi Caching Redis (Cache-Aside Pattern)

| Cache Key | Data | TTL | Invalidasi |
|-----------|------|-----|------------|
| `slots:{facility_id}:{date}` | Jumlah slot tersedia per poli | 24 jam | Saat ada booking/cancel |
| `facility:list` | Daftar Puskesmas + metadata | 6 jam | Manual via admin |
| `stats:disease:{date}` | Agregat kasus penyakit harian | 1 jam | Cron job setiap jam |
| `doctor:schedule:{doctor_id}` | Jadwal dokter per minggu | 12 jam | Saat jadwal diubah |

---

## 5 · Skema Database (PostgreSQL)

| Tabel | Kolom Kunci | Index | Catatan |
|-------|-------------|-------|---------|
| **users** | `PK` id · nik · email · phone · role · created_at | IDX nik · email | role: PATIENT \| NAKES \| DINKES |
| **appointments** | `PK` id · `FK` user_id · facility_id · doctor_id · scheduled_at · status · qr_token | IDX user_id · scheduled_at · status | status: PENDING \| CONFIRMED \| DONE \| CANCELLED |
| **medical_records** | `PK` id · `FK` user_id · appointment_id · diagnosis_enc · notes_enc · created_by | IDX user_id · created_at | diagnosis_enc & notes_enc: AES-256 |
| **facilities** | `PK` id · name · address · lat · lng · type · district_id | IDX district_id · type | type: PUSKESMAS \| KLINIK |
| **disease_logs** | `PK` id · `FK` record_id · icd10_code · facility_id · logged_at · district_id | IDX icd10_code · logged_at · district_id | Source data untuk surveilans epidemiologi |
| **audit_logs** | `PK` id · user_id · action · entity_type · entity_id · ip_addr · timestamp | IDX user_id · timestamp | Immutable. Insert-only. Partisi by month. |

---

## 6 · API Contract (REST, v1)

### Booking Service

```
GET    /api/v1/facilities?district={id}&type={type}
GET    /api/v1/facilities/{id}/slots?date={YYYY-MM-DD}
POST   /api/v1/appointments — body: {user_id, facility_id, doctor_id, date, poli}
DELETE /api/v1/appointments/{id} — batalkan booking
```

### EMR Service

```
GET    /api/v1/records?user_id={id}&limit=20&cursor={token}
POST   /api/v1/records — create rekam medis baru
PUT    /api/v1/records/{id} — update diagnosis/notes
```

### Analytics Service

```
GET    /api/v1/analytics/disease?icd10={code}&from={date}&to={date}&district={id}
GET    /api/v1/analytics/heatmap?date={YYYY-MM-DD}&icd10={code}
```

> Semua endpoint memerlukan header `Authorization: Bearer {JWT}`. Rate limit: 100 req/menit per IP.

---

## 7 · Kebutuhan Non-Fungsional

### Performa
**Response Time API** — P99 < 500ms. P50 < 200ms. Cache hit rate target ≥ 85%.

`< 200ms P50`

### Skalabilitas
**Horizontal Pod Autoscaler** — Booking Service: HPA trigger CPU > 70%, min 2 pod, max 10 pod.

`min 2 · max 10 pods`

### Ketersediaan
**SLA Target** — 99.5% uptime monthly. Liveness + Readiness probe pada semua pod Go.

`99.5% uptime`

### Keamanan
**Autentikasi & Enkripsi** — JWT RS256 antar service. AES-256 untuk rekam medis. TLS 1.3 mandatory. OWASP Top 10 compliance.

`JWT + AES-256`

### Keandalan Pesan
**RabbitMQ Reliability** — Queue durable + persistent message. Dead Letter Queue untuk pesan gagal setelah 3 retry.

`0 message loss target`

### Observabilitas
**Monitoring Stack** — Prometheus + Grafana untuk metrics. Distributed tracing (OpenTelemetry). Log aggregation (ELK/Loki).

`Full observability`

---

## 8 · CI/CD Pipeline (Jenkins)

| Step | Nama | Deskripsi |
|------|------|-----------|
| S1 | **Checkout & lint** | git checkout → golangci-lint (Go) → ESLint (Next.js) → gagal jika ada error kritis. |
| S2 | **Unit test** | `go test ./...` dengan coverage threshold ≥ 80%. Jest untuk frontend. Gagal pipeline jika di bawah threshold. |
| S3 | **Integration test** | Spin up PostgreSQL + Redis + RabbitMQ via Docker Compose. Jalankan test integrasi antar service. |
| S4 | **Build & push Docker image** | `docker build` multi-stage → tag dengan git SHA pendek → push ke private registry. Image scan (Trivy) untuk vulnerabilitas kritis. |
| S5 | **Deploy staging** | `kubectl set image deployment/{service} {container}={image}:{tag} -n staging`. Tunggu rollout selesai (`kubectl rollout status`). |
| S6 | **Smoke test & approval** | Jalankan smoke test endpoint utama di staging. Manual approval gate sebelum deploy ke production (Dinkes persona minta SLA ketat). |
| S7 | **Deploy production** | Rolling update ke namespace production. Rollback otomatis jika health check gagal dalam 5 menit pasca deploy. |

---

## 9 · Risiko & Mitigasi

| Severity | Risiko | Mitigasi |
|----------|--------|----------|
| 🔴 Tinggi | **Kebocoran data rekam medis** | Enkripsi AES-256 di DB, JWT scope-based, audit log immutable, penetration test sebelum launch, PDPA compliance review. |
| 🔴 Tinggi | **API Dukcapil tidak tersedia** | Implementasi circuit breaker. Fallback: validasi NIK format-only (sementara). Antrian booking tetap berjalan tanpa blokir. |
| 🟡 Sedang | **Redis down — data slot hilang** | Redis Sentinel/Cluster untuk HA. Fallback: query langsung ke PostgreSQL dengan degradasi performa yang dapat ditoleransi. |
| 🟡 Sedang | **Lonjakan traffic saat BPJS enrollment** | HPA Kubernetes sudah dikonfigurasi. Load test dengan k6 sebelum go-live. Rate limiting di API Gateway. |
| 🔵 Rendah | **WhatsApp Business API rate limit** | Queue RabbitMQ menyerap burst. Pengiriman notifikasi di-throttle. Email sebagai fallback channel. |

---

## 10 · Success Metrics & KPI

| Metrik | Target | Keterangan |
|--------|--------|------------|
| API response time (P50) | **< 200ms** | Berkat Redis cache. Diukur via Prometheus. |
| Uptime monthly SLA | **99.5%** | Selama bulan pertama. Alerting via Grafana. |
| Pesan notifikasi hilang | **0 msg** | RabbitMQ persistent. Monitor DLQ size. |
| Redis cache hit rate | **≥ 85%** | Untuk endpoint listing. Redis INFO keyspace. |
| Unit test coverage | **≥ 80%** | Per service Go/Next.js. Jenkins pipeline gate. |
| Mean time to recovery (MTTR) | **< 30 menit** | Saat incident. Rollback otomatis K8s. |

---

## 11 · Open Questions

> ❓ **Integrasi Dukcapil:** Apakah sudah ada MOU dengan Disdukcapil untuk akses API validasi NIK? Apa skema autentikasinya (API key / OAuth)?

> ❓ **Enkripsi rekam medis:** Apakah kunci enkripsi dikelola secara terpusat (KMS) atau per-tenant per-Puskesmas? Siapa key custodian-nya?

> ❓ **Multi-tenancy:** Apakah setiap Puskesmas punya namespace K8s sendiri, atau satu cluster shared dengan RBAC separation?

> ❓ **Regulasi data kesehatan:** Perlu konfirmasi apakah sistem harus memenuhi standar PERMENKES terkait rekam medis elektronik sebelum go-live.

> ❓ **SLA WhatsApp:** Apakah ada fallback SMS jika WhatsApp Business API tidak tersedia? Vendor SMS sudah ditentukan?
