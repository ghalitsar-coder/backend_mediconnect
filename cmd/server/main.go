package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"mediconnect/config"
	httpDelivery "mediconnect/internal/delivery/http"
	"mediconnect/internal/delivery/http/handler"
	"mediconnect/internal/domain"
	"mediconnect/internal/repository/postgres"
	"mediconnect/internal/seed"
	"mediconnect/internal/usecase"
	"mediconnect/pkg/database"
	"mediconnect/pkg/jwt"
	pkgLogger "mediconnect/pkg/logger"
	"mediconnect/pkg/storage"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Logger
	pkgLogger.Init(cfg.AppEnv)

	// Connect Database (PostgreSQL)
	db, err := database.ConnectPostgres(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect database: %v\n", err)
	}

	// AutoMigrate missing schema
	log.Println("Migrating PostgreSQL schema...")
	db.AutoMigrate(&domain.User{}, &domain.Facility{}, &domain.Doctor{}, &domain.Booking{})

	// Jalankan Seeder (idempotent: ON CONFLICT DO NOTHING)
	if err := seed.SeedAll(db); err != nil {
		log.Printf("⚠️  Seeding failed (non-fatal): %v", err)
	}
	jwtManager := jwt.NewJWTManager(
		"mysecret",     // sebaiknya dari config
		24*time.Hour,   // access expiry
		7*24*time.Hour, // refresh expiry
	)
	// Initialize Repositories
	authRepo := postgres.NewAuthRepository(db)
	facilityRepo := postgres.NewFacilityRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)
	doctorRepo := postgres.NewDoctorRepository(db)

	// Initialize Usecases
	authUsecase := usecase.NewAuthUsecase(authRepo, "mysecret") // Ideally from config
	facilityUsecase := usecase.NewFacilityUsecase(facilityRepo)
	// Pass nil for rabbit MQ since we are dropping it
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo, authRepo, nil)
	doctorUsecase := usecase.NewDoctorUsecase(doctorRepo)

	// Initialize Storage Service
	blobService, err := storage.NewAzureBlobService(cfg.AzureStorageConnectionString)
	if err != nil {
		log.Printf("⚠️  Failed to initialize blob service: %v", err)
	}

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(authUsecase, jwtManager)
	facilityHandler := handler.NewFacilityHandler(facilityUsecase)
	bookingHandler := handler.NewBookingHandler(bookingUsecase)
	doctorHandler := handler.NewDoctorHandler(doctorUsecase)
	uploadHandler := handler.NewUploadHandler(blobService, authRepo)

	// Setup Router
	router := httpDelivery.SetupRouter(authHandler, facilityHandler, bookingHandler, doctorHandler, uploadHandler, jwtManager)
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("Starting Mediconnect server on port %s\n", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen and Serve error: %v\n", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
