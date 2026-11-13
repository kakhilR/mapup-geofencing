package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"mapup-geofencing/internal/service"
)

type ViolationHandler struct {
	service *service.ViolationService
}

func NewViolationHandler(
	service *service.ViolationService,
) *ViolationHandler {
	return &ViolationHandler{
		service: service,
	}
}

func (h *ViolationHandler) GetViolations(
	w http.ResponseWriter,
	r *http.Request,
) {
	start := time.Now()

	violations, err := h.service.GetAllViolations(r.Context())
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := map[string]interface{}{
		"violations": violations,
		"time_ns":    time.Since(start).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}