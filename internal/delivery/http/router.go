package http

import (
	"mediconnect/internal/delivery/http/handler"
	"mediconnect/internal/delivery/http/middleware"
	"mediconnect/pkg/jwt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	facilityHandler *handler.FacilityHandler,
	bookingHandler *handler.BookingHandler,
	doctorHandler *handler.DoctorHandler,
	uploadHandler *handler.UploadHandler,
	jwtManager *jwt.JWTManager) *gin.Engine {

	router := gin.Default()

	// CORS config ...
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://70.153.84.104", "https://mediconnect-ghal.duckdns.org"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	api := router.Group("/api/v1")
	{
		api.GET("/health", handler.HealthHandler)

		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", authHandler.Logout)
		}

		facilities := api.Group("/facilities")
		{
			facilities.GET("", facilityHandler.GetFacilities)
		}

		doctors := api.Group("/doctors")
		{
			doctors.GET("", doctorHandler.GetDoctors)
			doctors.GET("/:id/slots", doctorHandler.GetSlots)
		}

		// --- Grup yang membutuhkan autentikasi ---
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth(jwtManager))
		{
			// Booking endpoints
			protected.POST("/bookings", bookingHandler.CreateBooking)
			protected.GET("/bookings", bookingHandler.GetMyBookings)

			// Upload KTP
			protected.POST("/uploads/ktp", uploadHandler.UploadKTP)
			protected.GET("/auth/me", authHandler.GetMe)
		}
	}

	return router
}
