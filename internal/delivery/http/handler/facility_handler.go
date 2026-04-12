package handler

import (
  "net/http"
  "mediconnect/internal/domain"
  "mediconnect/pkg/response"
  "github.com/gin-gonic/gin"
)

type FacilityHandler struct {
  facilityUsecase domain.FacilityUsecase
}

func NewFacilityHandler(uc domain.FacilityUsecase) *FacilityHandler {
  return &FacilityHandler{facilityUsecase: uc}
}

func (h *FacilityHandler) GetFacilities(c *gin.Context) {
  filter := domain.FacilityFilter{
      DistrictID: c.Query("district"),
      Type:       c.Query("type"),
  }

  facilities, err := h.facilityUsecase.GetFacilities(c.Request.Context(), filter)
  if err != nil {
      response.Error(c, http.StatusInternalServerError, "Failed to retrieve facilities")
      return
  }

  if facilities == nil || len(facilities) == 0 {
      facilities = []domain.Facility{}
  }

  response.Success(c, http.StatusOK, "Facilities retrieved successfully", facilities)
}

func HealthHandler(c *gin.Context) {
  response.Success(c, http.StatusOK, "MediConnect API v1 is up and running", map[string]string{
      "version": "1.0.0",
      "status":  "ok",
  })
}
