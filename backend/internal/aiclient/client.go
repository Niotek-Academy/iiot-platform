package aiclient

import (
	"context"
	"time"
)

type TelemetryPoint struct {
	Timestamp time.Time          `json:"timestamp"`
	Metrics   map[string]float64 `json:"metrics"`
}

type PredictRequest struct {
	MachineID       string           `json:"machine_id"`
	WindowSize      int              `json:"window_size"`
	TelemetryWindow []TelemetryPoint `json:"telemetry_window"`
}

type PredictResponse struct {
	MachineID         string    `json:"machine_id"`
	EvaluatedAt       time.Time `json:"evaluated_at"`
	HealthScore       float64   `json:"health_score"`
	IsAnomaly         bool      `json:"is_anomaly"`
	AnomalyFeatures   []string  `json:"anomaly_features"`
	RULHours          float64   `json:"rul_hours"`
	RecommendedAction string    `json:"recommended_action"`
}

type AIClient interface {
	Predict(ctx context.Context, req PredictRequest) (PredictResponse, error)
}