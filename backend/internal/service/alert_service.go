package service

import (
	"context"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
)

type AlertService struct {
	store *db.Store
}

func NewAlertService(store *db.Store) *AlertService {
	return &AlertService{store: store}
}

// ListByMachine returns all alerts if isResolved is nil, otherwise only
// alerts matching that resolved status.
func (s *AlertService) ListByMachine(ctx context.Context, machineID string, isResolved *bool) ([]generated.Alert, error) {
	if _, err := s.store.GetMachine(ctx, machineID); err != nil {
		return nil, apperr.NewNotFound("machine not found")
	}

	if isResolved == nil {
		alerts, err := s.store.ListAlertsByMachine(ctx, machineID)
		if err != nil {
			return nil, apperr.NewInternal("could not list alerts")
		}
		return alerts, nil
	}

	alerts, err := s.store.ListAlertsByMachineAndStatus(ctx, generated.ListAlertsByMachineAndStatusParams{
		MachineID:  machineID,
		IsResolved: *isResolved,
	})
	if err != nil {
		return nil, apperr.NewInternal("could not list alerts")
	}
	return alerts, nil
}