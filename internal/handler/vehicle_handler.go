
package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"mapup-geofencing/internal/model"
	"mapup-geofencing/internal/service"
)

type VehicleHandler struct {
	service *service.VehicleService
}

func NewVehicleHandler(service *service.VehicleService) *VehicleHandler {
	return &VehicleHandler{
		service: service,
	}
}

func (h *VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req model.CreateVehicleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	vehicle, err := h.service.CreateVehicle(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := model.CreateVehicleResponse{
		ID:            vehicle.ID,
		VehicleNumber: vehicle.VehicleNumber,
		Status:        "active",
		TimeNS:        time.Since(start).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func (h *VehicleHandler) GetVehicles(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    vehicles, err := h.service.GetAllVehicles(r.Context())
    if err != nil {
        http.Error(w, "failed to fetch vehicles", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "vehicles": vehicles,
        "time_ns":  time.Since(start).Nanoseconds(),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    json.NewEncoder(w).Encode(response)
}