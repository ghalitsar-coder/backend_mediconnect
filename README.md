# MediConnect ID — Backend API

> Platform Kesehatan Terpadu | Go · PostgreSQL · Redis · RabbitMQ · Kubernetes

[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## 📐 Arsitektur

Proyek ini menggunakan **Clean Architecture** dengan pemisahan layer yang ketat:

```
cmd/server/main.go          ← Entry point & dependency injection
│
├── config/                 ← Konfigurasi aplikasi (env vars)
│
├── internal/
│   ├── domain/             ← Entities + Repository/Usecase interfaces (CORE)
│   ├── usecase/            ← Business logic
│   ├── delivery/http/      ← HTTP handlers, middleware, router
│   └── repository/postgres/← Data access (PostgreSQL)
│
└── pkg/
    ├── database/           ← PostgreSQL & Redis connection helpers
    ├── messaging/          ← RabbitMQ connection
    ├── logger/             ← Structured logger (slog)
    └── response/           ← Standard JSON response envelope
```

**Dependency Rule:** Outer layers depend on inner layers. Domain tidak bergantung pada apapun.

---

## 🚀 Quick Start

### Prasyarat
- Go 1.23+
- Docker & Docker Compose

### Jalankan Lokal (dengan Docker)

```bash
# 1. Salin environment file
cp .env.example .env

# 2. Jalankan semua services (DB, Redis, RabbitMQ, App)
make docker-up

# 3. Cek health
curl http://localhost:8080/api/v1/health
```

### Jalankan Tanpa Docker

```bash
# 1. Pastikan PostgreSQL, Redis, RabbitMQ sudah berjalan secara lokal
cp .env.example .env

# 2. Jalankan server
make run
```

---

## 📋 Perintah Tersedia

```bash
make help          # Tampilkan semua perintah

make run           # Jalankan server lokal
make build         # Compile binary ke ./bin/
make test          # Jalankan semua unit test + coverage
make lint          # Jalankan golangci-lint
make fmt           # Format source code

make docker-up     # Start semua Docker services
make docker-down   # Stop semua Docker services
make migrate       # Terapkan SQL migrations ke DB container
make clean         # Hapus build artefacts
```

---

## 🌐 API Endpoints

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/api/v1/health` | Liveness check |
| `GET` | `/api/v1/facilities` | List fasilitas kesehatan |

**Query params `/facilities`:**
- `district` — filter berdasarkan kode wilayah BPS
- `type` — filter berdasarkan tipe (`PUSKESMAS` \| `KLINIK`)

---

## 🗂 Struktur Direktori Lengkap

```
backend_mediconnect/
├── cmd/server/             ← main.go (entry point)
├── config/                 ← LoadConfig()
├── internal/
│   ├── domain/             ← Entities & interfaces
│   ├── usecase/            ← Business logic layer
│   ├── delivery/http/
│   │   ├── handler/        ← HTTP handlers
│   │   ├── middleware/     ← JWT auth middleware
│   │   └── router.go       ← Route registration
│   └── repository/
│       └── postgres/       ← PostgreSQL implementations
├── pkg/
│   ├── database/           ← DB connection helpers
│   ├── logger/             ← slog wrapper
│   ├── messaging/          ← RabbitMQ client
│   └── response/           ← JSON response helpers
├── migrations/             ← SQL schema migrations
├── docs/                   ← PRD & API documentation
├── scripts/                ← Utility shell scripts
├── Dockerfile
├── docker-compose.yml
├── Jenkinsfile             ← CI/CD pipeline (7 stages)
├── Makefile
└── .env.example
```

---

## 🛠 Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Language | Go 1.23 |
| HTTP Router | chi v5 |
| Database | PostgreSQL 16 (pgx/v5) |
| Cache | Redis 7 |
| Message Broker | RabbitMQ 3.13 |
| Container | Docker + Kubernetes 1.30 |
| CI/CD | Jenkins |

---

## 📖 Dokumentasi

- [PRD — Product Requirements Document](docs/PRD.md)
- [Database Schema](migrations/001_init_schema.sql)
