CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_number VARCHAR(50) NOT NULL UNIQUE,
    driver_name VARCHAR(255) NOT NULL,
    vehicle_type VARCHAR(50) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);