package mapper

import (
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
)

// ToTelemetryResponse expects logs already in chronological (oldest-first)
// order — the service is responsible for reversing the DESC-ordered DB
// result before calling this.
func ToTelemetryResponse(sensorID, unit string, logs []generated.TelemetryLog) dto.TelemetryResponse {
	points := make([]dto.TelemetryDataPoint, 0, len(logs))
	for _, l := range logs {
		points = append(points, dto.TelemetryDataPoint{
			Value:      l.MetricValue,
			RecordedAt: l.RecordedAt.Time,
		})
	}
	return dto.TelemetryResponse{
		SensorID:   sensorID,
		Unit:       unit,
		Count:      len(points),
		DataPoints: points,
	}
}