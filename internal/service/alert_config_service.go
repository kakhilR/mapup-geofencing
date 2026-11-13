package service

import (
	"context"
	"fmt"
	"strings"

	"mapup-geofencing/internal/model"
	"mapup-geofencing/internal/repository"
)

type AlertConfigService struct {
	repo *repository.AlertConfigRepository
}

func NewAlertConfigService(
	repo *repository.AlertConfigRepository,
) *AlertConfigService {
	return &AlertConfigService{
		repo: repo,
	}
}

func (s *AlertConfigService) CreateConfig(
	ctx context.Context,
	req model.CreateAlertConfigRequest,
) (model.AlertConfig, error) {

	if strings.TrimSpace(req.VehicleID) == "" {
		return model.AlertConfig{}, fmt.Errorf("vehicle_id is required")
	}

	if strings.TrimSpace(req.GeofenceID) == "" {
		return model.AlertConfig{}, fmt.Errorf("geofence_id is required")
	}

	req.EventType = strings.ToUpper(strings.TrimSpace(req.EventType))

	if req.EventType != "ENTRY" && req.EventType != "EXIT" {
		return model.AlertConfig{}, fmt.Errorf(
			"event_type must be ENTRY or EXIT",
		)
	}

	return s.repo.Create(ctx, req)
}

func (s *AlertConfigService) GetAll(
	ctx context.Context,
) ([]model.AlertConfig, error) {
	return s.repo.GetAll(ctx)
}