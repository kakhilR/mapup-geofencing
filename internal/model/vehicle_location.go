package model

import "time"

type VehicleLocation struct {
	ID         string    
	VehicleID  string   
	Latitude   float64  
	Longitude  float64   
	RecordedAt time.Time 
}

type VehicleLocationRequest struct {
    VehicleID string  `json:"vehicle_id"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

type VehicleLocationResponse struct {
    ID         string    `json:"id"`
    VehicleID  string    `json:"vehicle_id"`
    Latitude   float64   `json:"latitude"`
    Longitude  float64   `json:"longitude"`
    RecordedAt time.Time `json:"recorded_at"`
}