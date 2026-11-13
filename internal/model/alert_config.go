package model

import "time"

type CreateAlertConfigRequest struct {
	VehicleID  string `json:"vehicle_id"`
	GeofenceID string `json:"geofence_id"`
	EventType  string `json:"event_type"`
	Enabled    bool   `json:"enabled"`
}

type AlertConfig struct {
	ID         string    `json:"id"`
	VehicleID  string    `json:"vehicle_id"`
	GeofenceID string    `json:"geofence_id"`
	EventType  string    `json:"event_type"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
}