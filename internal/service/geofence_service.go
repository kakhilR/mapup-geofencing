package service

import (
	"context"
	"fmt"

	"mapup-geofencing/internal/model"
	"mapup-geofencing/internal/repository"
)

type GeofenceService struct {
	repo *repository.GeofenceRepository
}

func NewGeofenceService(repo *repository.GeofenceRepository) *GeofenceService {
	return &GeofenceService{
		repo: repo,
	}
}

func (s *GeofenceService) CreateGeofence(
	ctx context.Context,
	req model.CreateGeofenceRequest,
) (string, error) {

	if req.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	if !req.Category.IsValid() {
		return "", fmt.Errorf("invalid category")
	}

	if len(req.Coordinates) < 4 {
		return "", fmt.Errorf("coordinates must contain at least 4 points")
	}

	// Polygon must be closed.
	first := req.Coordinates[0]
	last := req.Coordinates[len(req.Coordinates)-1]

	if len(first) != 2 || len(last) != 2 ||
		first[0] != last[0] || first[1] != last[1] {
		return "", fmt.Errorf("first and last coordinates must be identical")
	}


	for i, point := range req.Coordinates {
		if len(point) != 2 {
			return "", fmt.Errorf("coordinate %d must contain latitude and longitude", i)
		}

		lat := point[0]
		lng := point[1]

		if lat < -90 || lat > 90 {
			return "", fmt.Errorf("latitude at coordinate %d must be between -90 and 90", i)
		}

		if lng < -180 || lng > 180 {
			return "", fmt.Errorf("longitude at coordinate %d must be between -180 and 180", i)
		}
	}

	return s.repo.Create(ctx, req)
}



func (s *GeofenceService) GetAllGeofences(
	ctx context.Context,
) ([]model.GeofenceResponse, error) {

	return s.repo.GetAll(ctx)
}

func (s *GeofenceService) FindContainingGeofences(
	ctx context.Context,
	latitude float64,
	longitude float64,
) ([]model.GeofenceResponse, error) {

	if latitude < -90 || latitude > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}

	if longitude < -180 || longitude > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}

	return s.repo.FindContainingGeofences(
		ctx,
		latitude,
		longitude,
	)
}