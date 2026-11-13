package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"mapup-geofencing/internal/model"
	"mapup-geofencing/internal/service"
)

type GeofenceHandler struct {
	service *service.GeofenceService
}

func NewGeofenceHandler(service *service.GeofenceService) *GeofenceHandler {
	return &GeofenceHandler{
		service: service,
	}
}

func (h *GeofenceHandler) CreateGeofence(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req model.CreateGeofenceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	id, err := h.service.CreateGeofence(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := model.CreateGeofenceResponse{
		ID:     id,
		Name:   req.Name,
		Status: "active",
		TimeNS: time.Since(start).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func (h *GeofenceHandler) GetGeofences(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	geofences, err := h.service.GetAllGeofences(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch geofences", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"geofences": geofences,
		"time_ns":   time.Since(start).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

// temp handler for finding geofences containing a specific point (latitude, longitude)
func (h *GeofenceHandler) FindContainingGeofences(
	w http.ResponseWriter,
	r *http.Request,
) {
	start := time.Now()

	latitude := 37.7799
	longitude := -122.4144

	geofences, err := h.service.FindContainingGeofences(
		r.Context(),
		latitude,
		longitude,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := map[string]interface{}{
		"geofences": geofences,
		"time_ns":   time.Since(start).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}