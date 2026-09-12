// Package evaluation periodically takes each machine's current sliding
// window (from Phase 4's ingestion.Manager) and runs it through the AI
// client, saving the result and raising an alert on anomalies.
package evaluation

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/aiclient"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/ingestion"
)

type Evaluator struct {
	store    *db.Store
	ingest   *ingestion.Manager // to get the window of each machine
	aiClient aiclient.AIClient
	interval time.Duration   // time between evaluations
	mu     sync.RWMutex
	latest map[string]aiclient.PredictResponse
}

func NewEvaluator(store *db.Store, ingest *ingestion.Manager, client aiclient.AIClient, interval time.Duration) *Evaluator {
	return &Evaluator{
		store:    store,
		ingest:   ingest,
		aiClient: client,
		interval: interval,
		latest:   make(map[string]aiclient.PredictResponse),
	}
}

// Run blocks until ctx is cancelled — call it in a goroutine from main.go.
func (e *Evaluator) Run(ctx context.Context) {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:  // wait for the next interval 
			e.evaluateAll(ctx)
		}
	}
}

func (e *Evaluator) evaluateAll(ctx context.Context) {
	machines, err := e.store.ListMachines(ctx)
	if err != nil {
		log.Printf("evaluation: could not list machines: %v", err)
		return
	}
	for _, m := range machines {
		e.evaluateOne(ctx, m.MachineID)
	}
}

func (e *Evaluator) evaluateOne(ctx context.Context, machineID string) {
	window, ok := e.ingest.Window(machineID)  // get the window of the machine
	if !ok || len(window) == 0 {
		return // no telemetry for this machine yet — nothing to evaluate
	}

	req := aiclient.PredictRequest{
		MachineID:       machineID,
		WindowSize:      len(window),
		TelemetryWindow: toTelemetryPoints(window),
	}

	resp, err := e.aiClient.Predict(ctx, req)
	if err != nil {
		log.Printf("evaluation: predict failed for %s: %v", machineID, err)
		return
	}
	e.mu.Lock()
	e.latest[machineID] = resp
	e.mu.Unlock()

	if _, err := e.store.InsertAIEvaluation(ctx, generated.InsertAIEvaluationParams{
		MachineID:   machineID,
		HealthScore: resp.HealthScore,
		IsAnomaly:   resp.IsAnomaly,
		RulHours:    resp.RULHours,
		EvaluatedAt: pgtype.Timestamptz{Time: resp.EvaluatedAt, Valid: true},
	}); err != nil {
		log.Printf("evaluation: could not save evaluation for %s: %v", machineID, err)
	}

	if resp.IsAnomaly {
		message := resp.RecommendedAction
		if message == "" {
			message = fmt.Sprintf("Anomaly detected in: %s", strings.Join(resp.AnomalyFeatures, ", "))
		}
		if _, err := e.store.CreateAlert(ctx, generated.CreateAlertParams{
			MachineID: machineID,
			Severity:  severityFromHealthScore(resp.HealthScore),
			Message:   message,
		}); err != nil {
			log.Printf("evaluation: could not create alert for %s: %v", machineID, err)
		}
	}
}

func severityFromHealthScore(score float64) string {
	switch {
	case score < 30:
		return "CRITICAL"
	case score < 60:
		return "WARNING"
	default:
		return "INFO"
	}
}


// Take data snapshots and convert them to telemetry points to be sutible for Ai client
func toTelemetryPoints(window []ingestion.Snapshot) []aiclient.TelemetryPoint {
	points := make([]aiclient.TelemetryPoint, 0, len(window))
	for _, snap := range window {
		points = append(points, aiclient.TelemetryPoint{
			Timestamp: snap.Timestamp,
			Metrics:   snap.Metrics,
		})
	}
	return points
}

func (e *Evaluator) Latest(machineID string) (aiclient.PredictResponse, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	resp, ok := e.latest[machineID]
	return resp, ok
}