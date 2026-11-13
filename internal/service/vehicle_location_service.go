package service

import (
	"context"
	"fmt"

	"mapup-geofencing/internal/model"
	"mapup-geofencing/internal/repository"
)

type VehicleLocationService struct {
	repo              *repository.VehicleLocationRepository
	geofenceService   *GeofenceService
	transitionService *GeofenceTransitionService
	violationService  *ViolationService
}

func NewVehicleLocationService(
	repo *repository.VehicleLocationRepository,
	geofenceService *GeofenceService,
	transitionService *GeofenceTransitionService,
	violationService *ViolationService,
) *VehicleLocationService {
	return &VehicleLocationService{
		repo:              repo,
		geofenceService:   geofenceService,
		transitionService: transitionService,
		violationService:  violationService,
	}
}

func (s *VehicleLocationService) CreateLocation(
	ctx context.Context,
	req model.VehicleLocationRequest,
) (model.VehicleLocationResponse, error) {

	if req.VehicleID == "" {
		return model.VehicleLocationResponse{}, fmt.Errorf("vehicle_id is required")
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		return model.VehicleLocationResponse{}, fmt.Errorf(
			"latitude must be between -90 and 90",
		)
	}

	if req.Longitude < -180 || req.Longitude > 180 {
		return model.VehicleLocationResponse{}, fmt.Errorf(
			"longitude must be between -180 and 180",
		)
	}

	location, err := s.repo.Create(ctx, req)
	if err != nil {
		return model.VehicleLocationResponse{}, err
	}

	// Find geofences containing the new location.
	geofences, err := s.geofenceService.FindContainingGeofences(
		ctx,
		req.Latitude,
		req.Longitude,
	)
	if err != nil {
		return model.VehicleLocationResponse{}, err
	}

	// Detect ENTRY / EXIT transitions.
	events, err := s.transitionService.DetectTransitions(
		ctx,
		req.VehicleID,
		geofences,
	)
	if err != nil {
		return model.VehicleLocationResponse{}, err
	}

	// Temporary: log detected events.
	for _, event := range events {
		_, err := s.violationService.CreateViolation(ctx, event)
		if err != nil {
			return model.VehicleLocationResponse{}, err
		}

		fmt.Printf(
			"Geofence event: vehicle=%s geofence=%s type=%s\n",
			event.VehicleID,
			event.GeofenceID,
			event.EventType,
		)
	}

	return location, nil
}

func (s *VehicleLocationService) GetLatestLocation(
	ctx context.Context,
	vehicleID string,
) (model.VehicleLocationResponse, error) {

	if vehicleID == "" {
		return model.VehicleLocationResponse{}, fmt.Errorf(
			"vehicle_id is required",
		)
	}

	return s.repo.GetLatest(ctx, vehicleID)
}