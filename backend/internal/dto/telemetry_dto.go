package dto

import "time"

type TelemetryDataPoint struct {
	Value      float64   `json:"value"`
	RecordedAt time.Time `json:"recorded_at"`
}

type TelemetryResponse struct {
	SensorID   string                `json:"sensor_id"`
	Unit       string                `json:"unit"`
	Count      int                   `json:"count"`
	DataPoints []TelemetryDataPoint `json:"data_points"`
}