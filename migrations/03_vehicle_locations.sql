CREATE TABLE vehicle_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    recorded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vehicle_locations_vehicle_id
ON vehicle_locations(vehicle_id);

CREATE INDEX idx_vehicle_locations_recorded_at
ON vehicle_locations(recorded_at);