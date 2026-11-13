package service

import (
	"context"
	"fmt"
	"strings"

	"mapup-geofencing/internal/model"
	"mapup-geofencing/internal/repository"
)

type VehicleService struct {
	repo *repository.VehicleRepository
}

func NewVehicleService(repo *repository.VehicleRepository) *VehicleService {
	return &VehicleService{
		repo: repo,
	}
}

func (s *VehicleService) GetAllVehicles(ctx context.Context) ([]model.VehicleResponse, error) {
	return s.repo.GetAll(ctx)
}

func (s *VehicleService) CreateVehicle(
	ctx context.Context,
	req model.CreateVehicleRequest,
) (model.VehicleResponse, error) {

	req.VehicleNumber = strings.TrimSpace(req.VehicleNumber)
	req.DriverName = strings.TrimSpace(req.DriverName)
	req.VehicleType = strings.TrimSpace(req.VehicleType)
	req.Phone = strings.TrimSpace(req.Phone)

	if req.VehicleNumber == "" {
		return model.VehicleResponse{}, fmt.Errorf("vehicle_number is required")
	}

	if req.DriverName == "" {
		return model.VehicleResponse{}, fmt.Errorf("driver_name is required")
	}

	if req.VehicleType == "" {
		return model.VehicleResponse{}, fmt.Errorf("vehicle_type is required")
	}

	if req.Phone == "" {
		return model.VehicleResponse{}, fmt.Errorf("phone is required")
	}

	return s.repo.Create(ctx, req)
}