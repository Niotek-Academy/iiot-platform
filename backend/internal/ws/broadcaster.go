package ws

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/evaluation"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/ingestion"
)

type aiHealthPayload struct {
	HealthScore     float64  `json:"health_score"`
	IsAnomaly       bool     `json:"is_anomaly"`
	AnomalyFeatures []string `json:"anomaly_features"`
	RULHours        float64  `json:"rul_hours"`
}

type streamPayload struct {
	EventType string             `json:"event_type"`
	Timestamp time.Time          `json:"timestamp"`
	MachineID string             `json:"machine_id"`
	Telemetry map[string]float64 `json:"telemetry"`
	AIHealth  *aiHealthPayload   `json:"ai_health,omitempty"`
}

// Broadcaster ticks every interval and pushes one payload per machine —
// entirely from in-memory sources (ingestion.Manager's window, evaluation.
// Evaluator's latest result), no database round-trip on the hot path.
type Broadcaster struct {
	hub       *Hub
	ingest    *ingestion.Manager         // to get in memory ring buffer
	evaluator *evaluation.Evaluator       
	interval  time.Duration
}

func NewBroadcaster(hub *Hub, ingest *ingestion.Manager, evaluator *evaluation.Evaluator, interval time.Duration) *Broadcaster {
	return &Broadcaster{hub: hub, ingest: ingest, evaluator: evaluator, interval: interval}
}

func (b *Broadcaster) Run(ctx context.Context) {
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, machineID := range b.ingest.KnownMachines() {  // known machines is the only connected machines
				b.broadcastOne(machineID)
			}
		}
	}
}

func (b *Broadcaster) broadcastOne(machineID string) {
	window, ok := b.ingest.Window(machineID)
	if !ok || len(window) == 0 {
		return
	}
	latest := window[len(window)-1]

	payload := streamPayload{
		EventType: "TELEMETRY_AND_AI_UPDATE",
		Timestamp: latest.Timestamp,
		MachineID: machineID,
		Telemetry: latest.Metrics,
	}

	if result, ok := b.evaluator.Latest(machineID); ok {
		payload.AIHealth = &aiHealthPayload{
			HealthScore:     result.HealthScore,
			IsAnomaly:       result.IsAnomaly,
			AnomalyFeatures: result.AnomalyFeatures,
			RULHours:        result.RULHours,
		}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ws broadcaster: marshal failed for %s: %v", machineID, err)
		return
	}

	b.hub.Broadcast(machineID, data)
}