CREATE TABLE vehicle_geofence_states (
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    geofence_id UUID NOT NULL REFERENCES geofences(id) ON DELETE CASCADE,
    is_inside BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (vehicle_id, geofence_id)
);