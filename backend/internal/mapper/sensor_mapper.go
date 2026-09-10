package mapper

import (
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
)

func ToSensorResponse(s generated.Sensor) dto.SensorResponse {
	return dto.SensorResponse{
		SensorID:   s.SensorID,
		MachineID:  s.MachineID,
		MetricName: s.MetricName,
		Unit:       s.Unit,
		CreatedAt:  s.CreatedAt.Time,
	}
}

func ToSensorResponseList(list []generated.Sensor) []dto.SensorResponse {
	out := make([]dto.SensorResponse, 0, len(list))
	for _, s := range list {
		out = append(out, ToSensorResponse(s))
	}
	return out
}