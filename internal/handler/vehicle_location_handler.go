package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"mapup-geofencing/internal/model"
	"mapup-geofencing/internal/service"
)

type VehicleLocationHandler struct {
	service *service.VehicleLocationService
}

func NewVehicleLocationHandler(
	service *service.VehicleLocationService,
) *VehicleLocationHandler {
	return &VehicleLocationHandler{
		service: service,
	}
}

func (h *VehicleLocationHandler) CreateLocation(
	w http.ResponseWriter,
	r *http.Request,
) {
	start := time.Now()

	var req model.VehicleLocationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	location, err := h.service.CreateLocation(
		r.Context(),
		req,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"id":         location.ID,
		"vehicle_id": location.VehicleID,
		"latitude":   location.Latitude,
		"longitude":  location.Longitude,
		"recorded_at": location.RecordedAt,
		"time_ns":    time.Since(start).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}


func (h *VehicleLocationHandler) GetLatestLocation(
	w http.ResponseWriter,
	r *http.Request,
) {
	start := time.Now()

	vehicleID := r.PathValue("vehicle_id")

	location, err := h.service.GetLatestLocation(
		r.Context(),
		vehicleID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"id":          location.ID,
		"vehicle_id":  location.VehicleID,
		"latitude":    location.Latitude,
		"longitude":   location.Longitude,
		"recorded_at": location.RecordedAt,
		"time_ns":     time.Since(start).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}