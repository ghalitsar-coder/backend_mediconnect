package handler

import (
	"net/http"

	"mediconnect/internal/domain"
	"mediconnect/pkg/response"
)

// FacilityHandler handles HTTP requests related to healthcare facilities.
type FacilityHandler struct {
	facilityUsecase domain.FacilityUsecase
}

// NewFacilityHandler wires a FacilityHandler with its usecase dependency.
func NewFacilityHandler(uc domain.FacilityUsecase) *FacilityHandler {
	return &FacilityHandler{facilityUsecase: uc}
}

// GetFacilities godoc
//
//	@Summary      List healthcare facilities
//	@Description  Returns a filterable list of active Puskesmas and Klinik
//	@Tags         facilities
//	@Produce      json
//	@Param        district  query  string  false  "BPS district code"
//	@Param        type      query  string  false  "Facility type (PUSKESMAS | KLINIK)"
//	@Success      200  {array}   domain.Facility
//	@Failure      500  {object}  response.APIResponse
//	@Router       /api/v1/facilities [get]
func (h *FacilityHandler) GetFacilities(w http.ResponseWriter, r *http.Request) {
	filter := domain.FacilityFilter{
		DistrictID: r.URL.Query().Get("district"),
		Type:       r.URL.Query().Get("type"),
	}

	facilities, err := h.facilityUsecase.GetFacilities(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve facilities")
		return
	}

	if facilities == nil {
		facilities = []domain.Facility{} // always return an array, never null
	}

	response.Success(w, http.StatusOK, "Facilities retrieved successfully", facilities)
}

// HealthHandler returns the API liveness status.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "MediConnect API v1 is up and running", map[string]string{
		"version": "1.0.0",
		"status":  "ok",
	})
}
