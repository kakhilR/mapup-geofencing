package model

import "time"

type Violation struct {
	ID         string    `json:"id"`
	VehicleID  string    `json:"vehicle_id"`
	GeofenceID string    `json:"geofence_id"`
	EventType  string    `json:"event_type"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}