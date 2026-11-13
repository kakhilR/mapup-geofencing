CREATE TABLE geofences (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    geometry GEOMETRY(POLYGON, 4326) NOT NULL,
    category VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_geofences_geometry
ON geofences
USING GIST (geometry);