package aiclient

import (
	"context"
	"fmt"
	"strings"
	"time"
)


type MockAIClient struct{}

func NewMockAIClient() *MockAIClient {
	return &MockAIClient{}
}

func (m *MockAIClient) Predict(ctx context.Context, req PredictRequest) (PredictResponse, error) {
	if len(req.TelemetryWindow) == 0 {
		return PredictResponse{}, fmt.Errorf("empty telemetry window")
	}

	latest := req.TelemetryWindow[len(req.TelemetryWindow)-1].Metrics

	healthScore := 100.0
	var anomalyFeatures []string

	if v, ok := latest["VIB_01"]; ok && v > 35 {
		healthScore -= (v - 35) * 2
		anomalyFeatures = append(anomalyFeatures, "VIB_01")
	}
	if v, ok := latest["CURR_01"]; ok && v > 20 {
		healthScore -= (v - 20) * 3
		anomalyFeatures = append(anomalyFeatures, "CURR_01")
	}
	if v, ok := latest["TE_01"]; ok && v > 85 {
		healthScore -= (v - 85) * 2
		anomalyFeatures = append(anomalyFeatures, "TE_01")
	}
	if v, ok := latest["LS_HIGH"]; ok && v >= 1 {
		healthScore -= 20
		anomalyFeatures = append(anomalyFeatures, "LS_HIGH")
	}

	if healthScore < 0 {
		healthScore = 0
	}
	if healthScore > 100 {
		healthScore = 100
	}

	isAnomaly := len(anomalyFeatures) > 0

	// RUL is Remaining Useful Life
	rulHours := 500.0 // arbitrary healthy baseline
	if isAnomaly {
		rulHours = healthScore * 5 // crude: lower health -> lower RUL
	}

	action := "No action needed."
	if isAnomaly {
		action = fmt.Sprintf("Inspect: %s", strings.Join(anomalyFeatures, ", "))
	}

	return PredictResponse{
		MachineID:         req.MachineID,
		EvaluatedAt:       time.Now(),
		HealthScore:       healthScore,
		IsAnomaly:         isAnomaly,
		AnomalyFeatures:   anomalyFeatures,
		RULHours:          rulHours,
		RecommendedAction: action,
	}, nil
}