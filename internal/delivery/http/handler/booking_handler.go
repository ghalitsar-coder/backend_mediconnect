package handler

import (
	"mediconnect/internal/domain"
	"mediconnect/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingUsecase domain.BookingUsecase
}

func NewBookingHandler(uc domain.BookingUsecase) *BookingHandler {
	return &BookingHandler{bookingUsecase: uc}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req domain.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	// get the user ID from the Context (set by the Auth Middleware)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDVal.(string)

	booking, err := h.bookingUsecase.CreateBooking(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Booking created successfully", booking)
}

func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDVal.(string)

	bookings, err := h.bookingUsecase.GetMyBookings(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Bookings fetched successfully", bookings)
}
