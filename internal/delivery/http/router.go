package http

import (
	"mediconnect/internal/delivery/http/handler"
	"mediconnect/internal/delivery/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	facilityHandler *handler.FacilityHandler,
	bookingHandler *handler.BookingHandler,
	doctorHandler *handler.DoctorHandler) *gin.Engine {

	router := gin.Default()

	// Setup CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
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

		bookings := api.Group("/bookings")
		bookings.Use(middleware.JWTAuth)
		{
			bookings.POST("", bookingHandler.CreateBooking)
			bookings.GET("", bookingHandler.GetMyBookings)
		}
	}

	return router
}
