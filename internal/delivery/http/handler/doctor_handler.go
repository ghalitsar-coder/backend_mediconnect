package handler

import (
	"net/http"
	"time"

	"mediconnect/internal/domain"
	"mediconnect/pkg/response"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	usecase domain.DoctorUsecase
}

func NewDoctorHandler(usecase domain.DoctorUsecase) *DoctorHandler {
	return &DoctorHandler{usecase: usecase}
}

func (h *DoctorHandler) GetDoctors(c *gin.Context) {
	facilityID := c.Query("facility_id")
	poliName := c.Query("poli")

	doctors, err := h.usecase.GetDoctors(c.Request.Context(), facilityID, poliName)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed fetching doctors")
		return
	}

	response.Success(c, http.StatusOK, "Doctors retrieved successfully", doctors)
}

func (h *DoctorHandler) GetSlots(c *gin.Context) {
	doctorID := c.Param("id")
	dateStr := c.Query("date")

	if doctorID == "" || dateStr == "" {
		response.Error(c, http.StatusBadRequest, "Doctor ID and date are required")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
		return
	}

	slots, err := h.usecase.GetAvailableSlots(c.Request.Context(), doctorID, date)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed fetching available slots")
		return
	}

	response.Success(c, http.StatusOK, "Slots retrieved successfully", slots)
}
