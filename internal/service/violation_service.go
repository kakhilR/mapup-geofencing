package service

import (
	"context"
	"encoding/json"
	"fmt"

	"mapup-geofencing/internal/model"
	"mapup-geofencing/internal/repository"
	ws "mapup-geofencing/internal/websocket"
)

type ViolationService struct {
	repo          *repository.ViolationRepository
	alertRepo     *repository.AlertConfigRepository
	vehicleRepo   *repository.VehicleRepository
	geofenceRepo  *repository.GeofenceRepository
	hub           *ws.Hub
}

func NewViolationService(
	repo *repository.ViolationRepository,
	alertRepo *repository.AlertConfigRepository,
	vehicleRepo *repository.VehicleRepository,
	geofenceRepo *repository.GeofenceRepository,
	hub *ws.Hub,
) *ViolationService {
	return &ViolationService{
		repo:         repo,
		alertRepo:    alertRepo,
		vehicleRepo:  vehicleRepo,
		geofenceRepo: geofenceRepo,
		hub:          hub,
	}
}

func (s *ViolationService) CreateViolation(
	ctx context.Context,
	event GeofenceEvent,
) (model.Violation, error) {

	enabled, err := s.alertRepo.IsEnabled(
		ctx,
		event.VehicleID,
		event.GeofenceID,
		event.EventType,
	)

	if err != nil {
		return model.Violation{}, err
	}

	if !enabled {
		return model.Violation{}, nil
	}

	vehicle, err := s.vehicleRepo.GetByID(ctx, event.VehicleID)
	if err != nil {
		return model.Violation{}, err
	}

	geofence, err := s.geofenceRepo.GetByID(ctx, event.GeofenceID)
	if err != nil {
		return model.Violation{}, err
	}

	message := fmt.Sprintf(
		"Vehicle %s %s geofence %s",
		vehicle.VehicleNumber,
		event.EventType,
		geofence.Name,
	)

	violation := model.Violation{
		VehicleID:  event.VehicleID,
		GeofenceID: event.GeofenceID,
		EventType:  event.EventType,
		Message:    message,
	}

	result, err := s.repo.Create(ctx, violation)
	if err != nil {
		return model.Violation{}, err
	}

	// Broadcast the persisted violation.
	payload, err := json.Marshal(result)
	if err != nil {
		return model.Violation{}, fmt.Errorf(
			"failed to marshal violation: %w",
			err,
		)
	}

	s.hub.Broadcast(payload)

	return result, nil
}

func (s *ViolationService) GetAllViolations(
	ctx context.Context,
) ([]model.Violation, error) {
	return s.repo.GetAll(ctx)
}