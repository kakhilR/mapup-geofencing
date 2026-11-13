package repository

import (
	"context"
	"fmt"

	"mapup-geofencing/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VehicleRepository struct {
	db *pgxpool.Pool
}

func NewVehicleRepository(db *pgxpool.Pool) *VehicleRepository {
	return &VehicleRepository{
		db: db,
	}
}

func (r *VehicleRepository) Create(
	ctx context.Context,
	req model.CreateVehicleRequest,
) (model.VehicleResponse, error) {

	var vehicle model.VehicleResponse

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO vehicles
			(vehicle_number, driver_name, vehicle_type, phone)
		VALUES
			($1, $2, $3, $4)
		RETURNING id, vehicle_number, driver_name, vehicle_type, phone, created_at
		`,
		req.VehicleNumber,
		req.DriverName,
		req.VehicleType,
		req.Phone,
	).Scan(
		&vehicle.ID,
		&vehicle.VehicleNumber,
		&vehicle.DriverName,
		&vehicle.VehicleType,
		&vehicle.Phone,
		&vehicle.CreatedAt,
	)

	if err != nil {
		return model.VehicleResponse{}, fmt.Errorf(
			"failed to create vehicle: %w",
			err,
		)
	}

	return vehicle, nil
}

func (r *VehicleRepository) GetAll(
	ctx context.Context,
) ([]model.VehicleResponse, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			vehicle_number,
			driver_name,
			vehicle_type,
			phone,
			created_at
		FROM vehicles
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []model.VehicleResponse

	for rows.Next() {
		var vehicle model.VehicleResponse

		if err := rows.Scan(
			&vehicle.ID,
			&vehicle.VehicleNumber,
			&vehicle.DriverName,
			&vehicle.VehicleType,
			&vehicle.Phone,
			&vehicle.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}

		vehicles = append(vehicles, vehicle)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading vehicles: %w", err)
	}

	return vehicles, nil
}


func (r *VehicleRepository) GetByID(
	ctx context.Context,
	id string,
) (model.VehicleResponse, error) {

	var vehicle model.VehicleResponse

	err := r.db.QueryRow(ctx, `
		SELECT id, vehicle_number, driver_name, vehicle_type, phone, created_at
		FROM vehicles
		WHERE id = $1
	`, id).Scan(
		&vehicle.ID,
		&vehicle.VehicleNumber,
		&vehicle.DriverName,
		&vehicle.VehicleType,
		&vehicle.Phone,
		&vehicle.CreatedAt,
	)

	return vehicle, err
}