package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VehicleGeofenceStateRepository struct {
	db *pgxpool.Pool
}

func NewVehicleGeofenceStateRepository(db *pgxpool.Pool) *VehicleGeofenceStateRepository {
	return &VehicleGeofenceStateRepository{
		db: db,
	}
}

type VehicleGeofenceState struct {
	VehicleID  string
	GeofenceID string
	IsInside   bool
	UpdatedAt  time.Time
}

func (r *VehicleGeofenceStateRepository) Get(
	ctx context.Context,
	vehicleID string,
	geofenceID string,
) (VehicleGeofenceState, error) {

	var state VehicleGeofenceState

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			vehicle_id,
			geofence_id,
			is_inside,
			updated_at
		FROM vehicle_geofence_states
		WHERE vehicle_id = $1
		  AND geofence_id = $2
		`,
		vehicleID,
		geofenceID,
	).Scan(
		&state.VehicleID,
		&state.GeofenceID,
		&state.IsInside,
		&state.UpdatedAt,
	)

	if err != nil {
		return VehicleGeofenceState{}, fmt.Errorf(
			"failed to get vehicle geofence state: %w",
			err,
		)
	}

	return state, nil
}



func (r *VehicleGeofenceStateRepository) Upsert(
	ctx context.Context,
	vehicleID string,
	geofenceID string,
	isInside bool,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO vehicle_geofence_states
			(vehicle_id, geofence_id, is_inside, updated_at)
		VALUES
			($1, $2, $3, NOW())
		ON CONFLICT (vehicle_id, geofence_id)
		DO UPDATE SET
			is_inside = EXCLUDED.is_inside,
			updated_at = NOW()
		`,
		vehicleID,
		geofenceID,
		isInside,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to update vehicle geofence state: %w",
			err,
		)
	}

	return nil
}


func (r *VehicleGeofenceStateRepository) GetByVehicle(
	ctx context.Context,
	vehicleID string,
) ([]VehicleGeofenceState, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			vehicle_id,
			geofence_id,
			is_inside,
			updated_at
		FROM vehicle_geofence_states
		WHERE vehicle_id = $1
		`,
		vehicleID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get vehicle geofence states: %w",
			err,
		)
	}

	defer rows.Close()

	var states []VehicleGeofenceState

	for rows.Next() {
		var state VehicleGeofenceState

		if err := rows.Scan(
			&state.VehicleID,
			&state.GeofenceID,
			&state.IsInside,
			&state.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to scan vehicle geofence state: %w",
				err,
			)
		}

		states = append(states, state)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"error while reading vehicle geofence states: %w",
			err,
		)
	}

	return states, nil
}