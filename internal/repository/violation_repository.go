package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"mapup-geofencing/internal/model"
)

type ViolationRepository struct {
	db *pgxpool.Pool
}

func NewViolationRepository(db *pgxpool.Pool) *ViolationRepository {
	return &ViolationRepository{db: db}
}

func (r *ViolationRepository) Create(
	ctx context.Context,
	violation model.Violation,
) (model.Violation, error) {

	var result model.Violation

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO violations (
			vehicle_id,
			geofence_id,
			event_type,
			message
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			vehicle_id,
			geofence_id,
			event_type,
			message,
			created_at
		`,
		violation.VehicleID,
		violation.GeofenceID,
		violation.EventType,
		violation.Message,
	).Scan(
		&result.ID,
		&result.VehicleID,
		&result.GeofenceID,
		&result.EventType,
		&result.Message,
		&result.CreatedAt,
	)

	if err != nil {
		return model.Violation{}, fmt.Errorf(
			"failed to create violation: %w",
			err,
		)
	}

	return result, nil
}

func (r *ViolationRepository) GetAll(
	ctx context.Context,
) ([]model.Violation, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			vehicle_id,
			geofence_id,
			event_type,
			message,
			created_at
		FROM violations
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch violations: %w", err)
	}
	defer rows.Close()

	var violations []model.Violation

	for rows.Next() {
		var violation model.Violation

		if err := rows.Scan(
			&violation.ID,
			&violation.VehicleID,
			&violation.GeofenceID,
			&violation.EventType,
			&violation.Message,
			&violation.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan violation: %w", err)
		}

		violations = append(violations, violation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading violations: %w", err)
	}

	return violations, nil
}