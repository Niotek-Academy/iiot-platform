package service

import (
	"context"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
)

type SensorService struct {
	store *db.Store
}

func NewSensorService(store *db.Store) *SensorService {
	return &SensorService{store: store}
}

func (s *SensorService) Create(ctx context.Context, req dto.CreateSensorRequest) (generated.Sensor, error) {
	if _, err := s.store.GetMachine(ctx, req.MachineID); err != nil {
		return generated.Sensor{}, apperr.NewBadRequest("machine_id does not exist")
	}
	if _, err := s.store.GetSensor(ctx, req.SensorID); err == nil {
		return generated.Sensor{}, apperr.NewConflict("sensor_id already exists")
	}

	sensor, err := s.store.CreateSensor(ctx, generated.CreateSensorParams{
		SensorID:   req.SensorID,
		MachineID:  req.MachineID,
		MetricName: req.MetricName,
		Unit:       req.Unit,
	})
	if err != nil {
		return generated.Sensor{}, apperr.NewInternal("could not create sensor")
	}
	return sensor, nil
}

// ListSensors returns all sensors when machineID is empty, otherwise only
// the sensors belonging to that machine.
func (s *SensorService) ListSensors(ctx context.Context, machineID string) ([]generated.Sensor, error) {
	if machineID == "" {
		list, err := s.store.ListSensors(ctx)
		if err != nil {
			return nil, apperr.NewInternal("could not list sensors")
		}
		return list, nil
	}

	if _, err := s.store.GetMachine(ctx, machineID); err != nil {
		return nil, apperr.NewNotFound("machine not found")
	}
	list, err := s.store.ListSensorsByMachine(ctx, machineID)
	if err != nil {
		return nil, apperr.NewInternal("could not list sensors")
	}
	
	return list, nil
}

func (s *SensorService) GetByID(ctx context.Context, sensorID string) (generated.Sensor, error) {
	sensor, err := s.store.GetSensor(ctx, sensorID)
	if err != nil {
		return generated.Sensor{}, apperr.NewNotFound("sensor not found")
	}
	return sensor, nil
}

func (s *SensorService) Delete(ctx context.Context, sensorID string) error {
	if _, err := s.store.GetSensor(ctx, sensorID); err != nil {
		return apperr.NewNotFound("sensor not found")
	}
	if err := s.store.DeleteSensor(ctx, sensorID); err != nil {
		return apperr.NewInternal("could not delete sensor")
	}
	return nil
}