package handler

import (
	"net/http"

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
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed fetching doctors"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Doctors retrieved successfully", doctors))
}
