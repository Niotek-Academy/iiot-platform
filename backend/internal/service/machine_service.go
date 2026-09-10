package service

import (
	"context"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/utils/pgutil"
)

type MachineService struct {
	store *db.Store
}

func NewMachineService(store *db.Store) *MachineService {
	return &MachineService{store: store}
}

func (s *MachineService) Create(ctx context.Context, req dto.CreateMachineRequest) (generated.Machine, error) {
	if _, err := s.store.GetMachine(ctx, req.MachineID); err == nil {
		return generated.Machine{}, apperr.NewConflict("machine_id already exists")
	}

	status := req.Status
	if status == "" {
		status = "OPERATIONAL"
	}

	machine, err := s.store.CreateMachine(ctx, generated.CreateMachineParams{
		MachineID: req.MachineID,
		Name:      req.Name,
		Location:  req.Location,
		Status:    status,
	})
	if err != nil {
		return generated.Machine{}, apperr.NewInternal("could not create machine")
	}
	return machine, nil
}

func (s *MachineService) List(ctx context.Context) ([]generated.Machine, error) {
	list, err := s.store.ListMachines(ctx)
	if err != nil {
		return nil, apperr.NewInternal("could not list machines")
	}
	return list, nil
}

func (s *MachineService) Update(ctx context.Context, machineID string, req dto.UpdateMachineRequest) (generated.Machine, error) {
	if _, err := s.store.GetMachine(ctx, machineID); err != nil {
		return generated.Machine{}, apperr.NewNotFound("machine not found")
	}

	machine, err := s.store.UpdateMachine(ctx, generated.UpdateMachineParams{
		MachineID: machineID,
		Name:      pgutil.ToText(req.Name),
		Location:  pgutil.ToText(req.Location),
		Status:    pgutil.ToText(req.Status),
	})
	if err != nil {
		return generated.Machine{}, apperr.NewInternal("could not update machine")
	}
	return machine, nil
}

func (s *MachineService) Delete(ctx context.Context, machineID string) error {
	if _, err := s.store.GetMachine(ctx, machineID); err != nil {
		return apperr.NewNotFound("machine not found")
	}
	if err := s.store.DeleteMachine(ctx, machineID); err != nil {
		return apperr.NewInternal("could not delete machine")
	}
	return nil
}