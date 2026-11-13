package service

import (
	"context"
	"fmt"

	"mapup-geofencing/internal/model"
	"mapup-geofencing/internal/repository"
)

type GeofenceEvent struct {
	VehicleID  string
	GeofenceID string
	EventType  string
}

type GeofenceTransitionService struct {
	stateRepo *repository.VehicleGeofenceStateRepository
}

func NewGeofenceTransitionService(
	stateRepo *repository.VehicleGeofenceStateRepository,
) *GeofenceTransitionService {
	return &GeofenceTransitionService{
		stateRepo: stateRepo,
	}
}

func (s *GeofenceTransitionService) DetectTransitions(
	ctx context.Context,
	vehicleID string,
	currentGeofences []model.GeofenceResponse,
) ([]GeofenceEvent, error) {

	previousStates, err := s.stateRepo.GetByVehicle(ctx, vehicleID)
	if err != nil {
		return nil, err
	}

	// Build lookup of previous state.
	previous := make(map[string]bool)

	for _, state := range previousStates {
		previous[state.GeofenceID] = state.IsInside
	}

	// Build lookup of currently containing geofences.
	current := make(map[string]bool)

	for _, geofence := range currentGeofences {
		current[geofence.ID] = true
	}

	var events []GeofenceEvent

	// Detect ENTRY.
	for geofenceID := range current {
		wasInside := previous[geofenceID]

		if !wasInside {
			events = append(events, GeofenceEvent{
				VehicleID:  vehicleID,
				GeofenceID: geofenceID,
				EventType:  "ENTRY",
			})
		}

		if err := s.stateRepo.Upsert(
			ctx,
			vehicleID,
			geofenceID,
			true,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to update geofence state: %w",
				err,
			)
		}
	}

	// Detect EXIT.
	for geofenceID, wasInside := range previous {
		if wasInside && !current[geofenceID] {
			events = append(events, GeofenceEvent{
				VehicleID:  vehicleID,
				GeofenceID: geofenceID,
				EventType:  "EXIT",
			})

			if err := s.stateRepo.Upsert(
				ctx,
				vehicleID,
				geofenceID,
				false,
			); err != nil {
				return nil, fmt.Errorf(
					"failed to update geofence state: %w",
					err,
				)
			}
		}
	}

	return events, nil
}