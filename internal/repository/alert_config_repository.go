package repository

import (
	"context"
	"errors"
	"fmt"

	"mapup-geofencing/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertConfigRepository struct {
	db *pgxpool.Pool
}

func NewAlertConfigRepository(db *pgxpool.Pool) *AlertConfigRepository {
	return &AlertConfigRepository{db: db}
}

func (r *AlertConfigRepository) Create(
	ctx context.Context,
	config model.CreateAlertConfigRequest,
) (model.AlertConfig, error) {

	var result model.AlertConfig

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO alert_configs (
			vehicle_id,
			geofence_id,
			event_type,
			enabled
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			vehicle_id,
			geofence_id,
			event_type,
			enabled,
			created_at
		`,
		config.VehicleID,
		config.GeofenceID,
		config.EventType,
		config.Enabled,
	).Scan(
		&result.ID,
		&result.VehicleID,
		&result.GeofenceID,
		&result.EventType,
		&result.Enabled,
		&result.CreatedAt,
	)

	if err != nil {
		return model.AlertConfig{}, fmt.Errorf(
			"failed to create alert config: %w",
			err,
		)
	}

	return result, nil
}

func (r *AlertConfigRepository) IsEnabled(
	ctx context.Context,
	vehicleID string,
	geofenceID string,
	eventType string,
) (bool, error) {

	var enabled bool

	err := r.db.QueryRow(
		ctx,
		`
		SELECT enabled
		FROM alert_configs
		WHERE vehicle_id = $1
		  AND geofence_id = $2
		  AND event_type = $3
		`,
		vehicleID,
		geofenceID,
		eventType,
	).Scan(&enabled)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf(
			"failed to check alert config: %w",
			err,
		)
	}

	return enabled, nil
}


func (r *AlertConfigRepository) GetAll(
	ctx context.Context,
) ([]model.AlertConfig, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			vehicle_id,
			geofence_id,
			event_type,
			enabled,
			created_at
		FROM alert_configs
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []model.AlertConfig

	for rows.Next() {
		var config model.AlertConfig

		err := rows.Scan(
			&config.ID,
			&config.VehicleID,
			&config.GeofenceID,
			&config.EventType,
			&config.Enabled,
			&config.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}