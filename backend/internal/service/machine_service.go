package service

import (
	"context"
	"errors"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/utils/pgutil"
	"github.com/jackc/pgx/v5"
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
		MachineID:      req.MachineID,
		Name:           req.Name,
		Location:       req.Location,
		Status:         status,
		ControlAddress: pgutil.ToText(req.ControlAddress),
	})
	if err != nil {
		return generated.Machine{}, apperr.NewInternal("could not create machine")
	}
	return machine, nil
}

func (s *MachineService) Update(ctx context.Context, machineID string, req dto.UpdateMachineRequest) (generated.Machine, error) {
	if _, err := s.store.GetMachine(ctx, machineID); err != nil {
		return generated.Machine{}, apperr.NewNotFound("machine not found")
	}

	machine, err := s.store.UpdateMachine(ctx, generated.UpdateMachineParams{
		MachineID:      machineID,
		Name:           pgutil.ToText(req.Name),
		Location:       pgutil.ToText(req.Location),
		Status:         pgutil.ToText(req.Status),
		ControlAddress: pgutil.ToText(req.ControlAddress),
	})
	if err != nil {
		return generated.Machine{}, apperr.NewInternal("could not update machine")
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

func (s *MachineService) Delete(ctx context.Context, machineID string) error {
	if _, err := s.store.GetMachine(ctx, machineID); err != nil {
		return apperr.NewNotFound("machine not found")
	}
	if err := s.store.DeleteMachine(ctx, machineID); err != nil {
		return apperr.NewInternal("could not delete machine")
	}
	return nil
}

func (s *MachineService) GetOverview(ctx context.Context, machineID string) (dto.MachineOverviewResponse, error) {
	machine, err := s.store.GetMachine(ctx, machineID)
	if err != nil {
		return dto.MachineOverviewResponse{}, apperr.NewNotFound("machine not found")
	}

	overview := dto.MachineOverviewResponse{
		MachineID: machine.MachineID,
		Name:      machine.Name,
		Location:  machine.Location,
		Status:    machine.Status,
		UpdatedAt: machine.CreatedAt.Time, // overwritten below if an AI evaluation exists
	}

	
	eval, err := s.store.GetLatestAIEvaluation(ctx, machineID)
	if err == nil {
		healthScore := eval.HealthScore
		rulHours := eval.RulHours
		overview.CurrentHealthScore = &healthScore
		overview.RULHours = &rulHours
		overview.UpdatedAt = eval.EvaluatedAt.Time
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return dto.MachineOverviewResponse{}, apperr.NewInternal("could not fetch machine health data")
	}

	count, err := s.store.CountUnresolvedAlerts(ctx, machineID)
	if err != nil {
		return dto.MachineOverviewResponse{}, apperr.NewInternal("could not count alerts")
	}
	overview.ActiveAlertsCount = count

	return overview, nil
}