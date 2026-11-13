package repository

import (
	"context"
	"fmt"

	"mapup-geofencing/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VehicleLocationRepository struct {
	db *pgxpool.Pool
}

func NewVehicleLocationRepository(db *pgxpool.Pool) *VehicleLocationRepository {
	return &VehicleLocationRepository{
		db: db,
	}
}

func (r *VehicleLocationRepository) Create(
	ctx context.Context,
	req model.VehicleLocationRequest,
) (model.VehicleLocationResponse, error) {

	var location model.VehicleLocationResponse

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO vehicle_locations
			(vehicle_id, latitude, longitude)
		VALUES
			($1, $2, $3)
		RETURNING id, vehicle_id, latitude, longitude, recorded_at
		`,
		req.VehicleID,
		req.Latitude,
		req.Longitude,
	).Scan(
		&location.ID,
		&location.VehicleID,
		&location.Latitude,
		&location.Longitude,
		&location.RecordedAt,
	)

	if err != nil {
		return model.VehicleLocationResponse{}, fmt.Errorf(
			"failed to create vehicle location: %w",
			err,
		)
	}

	return location, nil
}

func (r *VehicleLocationRepository) GetLatest(
	ctx context.Context,
	vehicleID string,
) (model.VehicleLocationResponse, error) {

	var location model.VehicleLocationResponse

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			vehicle_id,
			latitude,
			longitude,
			recorded_at
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY recorded_at DESC
		LIMIT 1
		`,
		vehicleID,
	).Scan(
		&location.ID,
		&location.VehicleID,
		&location.Latitude,
		&location.Longitude,
		&location.RecordedAt,
	)

	if err != nil {
		return model.VehicleLocationResponse{}, fmt.Errorf(
			"failed to fetch latest vehicle location: %w",
			err,
		)
	}

	return location, nil
}