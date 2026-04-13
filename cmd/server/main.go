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
	"mediconnect/internal/repository/postgres"
	"mediconnect/internal/usecase"
	"mediconnect/pkg/database"
	pkgLogger "mediconnect/pkg/logger"
	"mediconnect/pkg/messaging"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Logger
	pkgLogger.Init(cfg.AppEnv)

	// Connect Database (GORM)
	db, err := database.ConnectPostgres(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect database: %v\n", err)
	}

	// Connect RabbitMQ
	rabbit, err := messaging.ConnectRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect RabbitMQ: %v\n", err)
	}

	// Initialize Repositories
	authRepo := postgres.NewAuthRepository(db)
	facilityRepo := postgres.NewFacilityRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)
  doctorRepo := postgres.NewDoctorRepository(db)

  // Initialize Usecases
  authUsecase := usecase.NewAuthUsecase(authRepo, "mysecret") // Ideally from config
  facilityUsecase := usecase.NewFacilityUsecase(facilityRepo)
  bookingUsecase := usecase.NewBookingUsecase(bookingRepo, rabbit)
  doctorUsecase := usecase.NewDoctorUsecase(doctorRepo)

  // Initialize Handlers
  authHandler := handler.NewAuthHandler(authUsecase)
  facilityHandler := handler.NewFacilityHandler(facilityUsecase)
  bookingHandler := handler.NewBookingHandler(bookingUsecase)
  doctorHandler := handler.NewDoctorHandler(doctorUsecase)

  // Setup Router
  router := httpDelivery.SetupRouter(authHandler, facilityHandler, bookingHandler, doctorHandler)

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
