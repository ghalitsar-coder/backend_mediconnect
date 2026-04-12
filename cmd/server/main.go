package main

import (
	"log/slog"
	"net/http"
	"os"

	"mediconnect/config"
	httpdelivery "mediconnect/internal/delivery/http"
	"mediconnect/internal/delivery/http/handler"
	pgRepo "mediconnect/internal/repository/postgres"
	"mediconnect/internal/usecase"
	"mediconnect/pkg/database"
	"mediconnect/pkg/logger"
	"mediconnect/pkg/messaging"
)

func main() {
	// ── Configuration ────────────────────────────────────────────────────────
	cfg := config.LoadConfig()

	// ── Logger ───────────────────────────────────────────────────────────────
	logger.Init(cfg.AppEnv)
	slog.Info("Booting MediConnect API", "env", cfg.AppEnv, "port", cfg.ServerPort)

	// ── Infrastructure ───────────────────────────────────────────────────────
	db, err := database.ConnectPostgres(cfg.DBURL)
	if err != nil {
		slog.Error("PostgreSQL connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	redisClient, err := database.ConnectRedis(cfg.RedisURL)
	if err != nil {
		slog.Error("Redis connection failed", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	mq, err := messaging.ConnectRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		slog.Error("RabbitMQ connection failed", "error", err)
		os.Exit(1)
	}
	defer mq.Close()

	// ── Dependency Injection ─────────────────────────────────────────────────
	//  Repository layer
	facilityRepo := pgRepo.NewFacilityRepository(db)
	authRepo := pgRepo.NewAuthRepository(db)

	//  Usecase layer
	facilityUC := usecase.NewFacilityUsecase(facilityRepo)
	// Add your own secret key securely in prod (from env variable)
	jwtSecret := "supersecretkey"
	authUC := usecase.NewAuthUsecase(authRepo, jwtSecret)

	//  Handler layer
	facilityHandler := handler.NewFacilityHandler(facilityUC)
	authHandler := handler.NewAuthHandler(authUC)

	// ── HTTP Server ──────────────────────────────────────────────────────────
	router := httpdelivery.NewRouter(facilityHandler, authHandler)

	addr := ":" + cfg.ServerPort
	slog.Info("MediConnect API is ready", "addr", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		slog.Error("Server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
