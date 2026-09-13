package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/factoryio"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/utils/pgutil"
)

const (
	CommandEmergencyStop = "EMERGENCY_STOP"
	CommandStart         = "START"
	CommandResetAlerts   = "RESET_ALERTS"
)

type CommandService struct {
	store   *db.Store
	control factoryio.ControlClient
}

func NewCommandService(store *db.Store, control factoryio.ControlClient) *CommandService {
	return &CommandService{store: store, control: control}
}

// Execute runs a command for a machine and always logs it to command_logs,
// whether it came from a real user (issuedBy = username) or the AI
// evaluator (issuedBy = "AI_AUTONOMOUS_ENGINE", see AutoEmergencyStop below).
func (s *CommandService) Execute(ctx context.Context, machineID, command, reason, issuedBy string) (generated.CommandLog, error) {
	if _, err := s.store.GetMachine(ctx, machineID); err != nil {
		return generated.CommandLog{}, apperr.NewNotFound("machine not found")
	}

	switch command {
	case CommandEmergencyStop:
		if err := s.control.EmergencyStop(ctx, machineID); err != nil {
			return generated.CommandLog{}, apperr.NewInternal("failed to execute emergency stop on factory IO")
		}
		s.setMachineStatus(ctx, machineID, "STOPPED")

	case CommandStart:
		if err := s.control.Start(ctx, machineID); err != nil {
			return generated.CommandLog{}, apperr.NewInternal("failed to start machine on factory IO")
		}
		s.setMachineStatus(ctx, machineID, "OPERATIONAL")

	case CommandResetAlerts:
		if err := s.resetAlerts(ctx, machineID); err != nil {
			return generated.CommandLog{}, apperr.NewInternal("failed to reset alerts")
		}

	default:
		return generated.CommandLog{}, apperr.NewBadRequest("unsupported command")
	}

	logEntry, err := s.store.CreateCommandLog(ctx, generated.CreateCommandLogParams{
		MachineID:   machineID,
		CommandType: command,
		IssuedBy:    issuedBy,
		Reason:      reason,
	})
	if err != nil {
		return generated.CommandLog{}, apperr.NewInternal("could not log command")
	}
	return logEntry, nil
}

// AutoEmergencyStop is called directly by evaluation.Evaluator (in-process,
// not over HTTP) — this is what satisfies the evaluation.StopTrigger
// interface without evaluation needing to import this package.
func (s *CommandService) AutoEmergencyStop(ctx context.Context, machineID, reason string) error {
	_, err := s.Execute(ctx, machineID, CommandEmergencyStop, reason, "AI_AUTONOMOUS_ENGINE")
	return err
}

func (s *CommandService) setMachineStatus(ctx context.Context, machineID, status string) {
	_, err := s.store.UpdateMachine(ctx, generated.UpdateMachineParams{
		MachineID: machineID,
		Status:    pgutil.ToText(&status),
	})
	if err != nil {
		log.Printf("command: failed to update machine status to %s: %v", status, err)
	}
}

func (s *CommandService) resetAlerts(ctx context.Context, machineID string) error {
	alerts, err := s.store.ListAlertsByMachineAndStatus(ctx, generated.ListAlertsByMachineAndStatusParams{
		MachineID:  machineID,
		IsResolved: false,
	})
	if err != nil {
		return fmt.Errorf("list unresolved alerts: %w", err)
	}
	for _, a := range alerts {
		if _, err := s.store.ResolveAlert(ctx, a.AlertID); err != nil {
			return fmt.Errorf("resolve alert %d: %w", a.AlertID, err)
		}
	}
	return nil
}