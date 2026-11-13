package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"mapup-geofencing/internal/model"
	"mapup-geofencing/internal/service"
)

type AlertConfigHandler struct {
	service *service.AlertConfigService
}

func NewAlertConfigHandler(
	service *service.AlertConfigService,
) *AlertConfigHandler {
	return &AlertConfigHandler{
		service: service,
	}
}

func (h *AlertConfigHandler) CreateConfig(
	w http.ResponseWriter,
	r *http.Request,
) {
	start := time.Now()

	var req model.CreateAlertConfigRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	config, err := h.service.CreateConfig(
		r.Context(),
		req,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	response := map[string]interface{}{
		"id":          config.ID,
		"vehicle_id":  config.VehicleID,
		"geofence_id": config.GeofenceID,
		"event_type":  config.EventType,
		"enabled":     config.Enabled,
		"created_at":  config.CreatedAt,
		"time_ns":     time.Since(start).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}


func (h *AlertConfigHandler) GetAlerts(
	w http.ResponseWriter,
	r *http.Request,
) {
	start := time.Now()

	alerts, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"alerts":  alerts,
		"time_ns": time.Since(start).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}