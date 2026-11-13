package model

import "time"

type Geofence struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Geometry  string    `json:"geometry"`
	Description string    `json:"description"`
	Category    GeofenceCategory    `json:"category"` // convert this to enum
	CreatedAt time.Time `json:"created_at"`
}



type CreateGeofenceRequest struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Coordinates [][]float64   `json:"coordinates"` // convert the API's [][]float64 into a PostGIS polygon and insert it into the geometry column
	Category    GeofenceCategory `json:"category"`
}

type GeofenceResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Coordinates [][]float64 `json:"coordinates"`
	Category    GeofenceCategory `json:"category"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
}


type CreateGeofenceResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	TimeNS int64  `json:"time_ns"`
}