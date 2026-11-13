package repository

import (
	"context"
	"fmt"
	"time"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"mapup-geofencing/internal/model"
)

type GeofenceRepository struct {
	db *pgxpool.Pool
}

func NewGeofenceRepository(db *pgxpool.Pool) *GeofenceRepository {
	return &GeofenceRepository{
		db: db,
	}
}

func (r *GeofenceRepository) Create(
	ctx context.Context,
	req model.CreateGeofenceRequest,
) (string, error) {

	// Convert coordinates to PostGIS WKT polygon.
	var polygon string

	for i, point := range req.Coordinates {
		if i > 0 {
			polygon += ", "
		}

		lat := point[0]
		lng := point[1]

		polygon += fmt.Sprintf("%f %f", lng, lat)
	}

	polygon = fmt.Sprintf("POLYGON((%s))", polygon)

	// id := "geo_" + fmt.Sprintf("%d", 1) // TODO: Replace with a proper unique ID generator, e.g., UUID

	row := r.db.QueryRow(
		ctx,
		`
		INSERT INTO geofences
			(id, name, description, geometry, category)
		VALUES
			(gen_random_uuid(), $1, $2, ST_GeomFromText($3, 4326), $4)
		RETURNING id
		`,
		req.Name,
		req.Description,
		polygon,
		req.Category,
	)

	var id string // doubt point: here doesn't the postgresql genreate the id automatically just like mongodb does
	if err := row.Scan(&id); err != nil {
    return "", fmt.Errorf("failed to create geofence: %w", err)
}

	return id, nil
}


func (r *GeofenceRepository) GetAll(
	ctx context.Context,
) ([]model.GeofenceResponse, error) {

	// without converting wkt to GeoJSON, just fetching the raw geometry as WKT
	// rows, err := r.db.Query(
	// 	ctx,
	// 	`
	// 	SELECT
	// 		id,
	// 		name,
	// 		description,
	// 		ST_AsText(geometry),
	// 		category,
	// 		created_at
	// 	FROM geofences
	// 	ORDER BY created_at DESC
	// 	`,
	// )

	// convert wkt to GeoJSON (currently using ST_AsGeoJSON in the query)
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			name,
			description,
			ST_AsGeoJSON(geometry),
			category,
			created_at
		FROM geofences
		ORDER BY created_at DESC
		`,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch geofences: %w", err)
	}
	defer rows.Close()

	var geofences []model.GeofenceResponse

	for rows.Next() {
		var (
			id          string
			name        string
			description string
			geoJSON         string
			category    model.GeofenceCategory
			createdAt   time.Time
		)

		if err := rows.Scan(
			&id,
			&name,
			&description,
			&geoJSON,
			&category,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan geofence: %w", err)
		}

		var geometry struct {
			Type        string          `json:"type"`
			Coordinates [][][]float64   `json:"coordinates"`
		}

		if err := json.Unmarshal([]byte(geoJSON), &geometry); err != nil {
			return nil, fmt.Errorf("failed to parse geofence geometry: %w", err)
		}

		coordinates := make([][]float64, 0, len(geometry.Coordinates[0]))

		for _, point := range geometry.Coordinates[0] {
			coordinates = append(coordinates, []float64{
				point[1], // latitude
				point[0], // longitude
			})
		}

		geofences = append(geofences, model.GeofenceResponse{
			ID:          id,
			Name:        name,
			Description: description,
			Coordinates: coordinates,
			Category:    category,
			Status:      "active",
			CreatedAt:   createdAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading geofences: %w", err)
	}

	return geofences, nil
}


func (r *GeofenceRepository) FindContainingGeofences(
	ctx context.Context,
	latitude float64,
	longitude float64,
) ([]model.GeofenceResponse, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			name,
			description,
			ST_AsGeoJSON(geometry),
			category,
			created_at
		FROM geofences
		WHERE ST_Contains(
			geometry,
			ST_SetSRID(
				ST_MakePoint($1, $2),
				4326
			)
		)
		`,
		longitude,
		latitude,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to find containing geofences: %w",
			err,
		)
	}

	defer rows.Close()

	geofences := make([]model.GeofenceResponse, 0)

	for rows.Next() {
		var (
			id          string
			name        string
			description string
			geoJSON     string
			category    model.GeofenceCategory
			createdAt   time.Time
		)

		if err := rows.Scan(
			&id,
			&name,
			&description,
			&geoJSON,
			&category,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to scan geofence: %w",
				err,
			)
		}

		var geometry struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		}

		if err := json.Unmarshal(
			[]byte(geoJSON),
			&geometry,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to parse geofence geometry: %w",
				err,
			)
		}

		coordinates := make([][]float64, 0)

		if len(geometry.Coordinates) > 0 {
			for _, point := range geometry.Coordinates[0] {
				coordinates = append(
					coordinates,
					[]float64{
						point[1],
						point[0],
					},
				)
			}
		}

		geofences = append(
			geofences,
			model.GeofenceResponse{
				ID:          id,
				Name:        name,
				Description: description,
				Coordinates: coordinates,
				Category:    category,
				Status:      "active",
				CreatedAt:   createdAt,
			},
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"error while reading geofences: %w",
			err,
		)
	}

	return geofences, nil
}


func (r *GeofenceRepository) GetByID(
	ctx context.Context,
	id string,
) (model.GeofenceResponse, error) {

	var geofence model.GeofenceResponse
	var geoJSON string

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			name,
			description,
			ST_AsGeoJSON(geometry),
			category,
			created_at
		FROM geofences
		WHERE id = $1
	`, id).Scan(
		&geofence.ID,
		&geofence.Name,
		&geofence.Description,
		&geoJSON,
		&geofence.Category,
		&geofence.CreatedAt,
	)

	if err != nil {
		return geofence, err
	}

	// You already have the GeoJSON coordinate conversion
	// logic in GetAll(). Reuse that helper here.

	geofence.Status = "active"

	return geofence, nil
}