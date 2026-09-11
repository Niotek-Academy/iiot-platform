package service

import (
	"context"



	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/mapper"
)

type TelemetryService struct {
	store *db.Store
}

func NewTelemetryService(store *db.Store) *TelemetryService {
	return &TelemetryService{store: store}
}

func (s *TelemetryService) GetHistory(ctx context.Context, sensorID string, limit int32) (dto.TelemetryResponse, error) {
	sensor, err := s.store.GetSensor(ctx, sensorID)
	if err != nil {
		return dto.TelemetryResponse{}, apperr.NewNotFound("sensor not found")
	}

	logs, err := s.store.GetRecentTelemetryBySensor(ctx, generated.GetRecentTelemetryBySensorParams{
		SensorID: sensorID,
		Limit:    limit,
	})
	if err != nil {
		return dto.TelemetryResponse{}, apperr.NewInternal("could not fetch telemetry")
	}

	// DB returns newest-first (for an efficient LIMIT); reverse to
	// chronological order for the response, matching the API contract.
	reverse(logs)

	return mapper.ToTelemetryResponse(sensor.SensorID, sensor.Unit, logs), nil
}

func reverse(logs []generated.TelemetryLog) {
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}
}
