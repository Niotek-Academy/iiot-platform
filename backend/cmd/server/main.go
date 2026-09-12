package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/aiclient"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/config"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/evaluation"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/factoryio"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/ingestion"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/server"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/ws"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.LoadConfig()

	store := mustConnectDB(ctx, cfg)
	defer store.Close()

	manager := startIngestion(ctx, cfg, store)
	evaluator := startEvaluation(ctx, cfg, store, manager)
	hub := startWebSocket(ctx, manager, evaluator)

	router := server.NewRouter(store, cfg, hub)
	runServer(ctx, cfg, router)
}

// mustConnectDB opens the database connection or exits the process — there
// is no meaningful way to run without it.
func mustConnectDB(ctx context.Context, cfg config.Config) *db.Store {
	log.Printf("connecting to database...")
	store, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Printf("database connection established")
	return store
}

// startIngestion wires Phase 4: loads the sensor->machine map, picks the
// factory IO client, and runs ingestion in the background.
func startIngestion(ctx context.Context, cfg config.Config, store *db.Store) *ingestion.Manager {
	manager := ingestion.NewManager(store, cfg.BufferCapacity, cfg.SnapshotInterval)
	if err := manager.LoadSensorMap(ctx); err != nil {
		log.Fatalf("failed to load sensor map: %v", err)
	}

	ioClient := newFactoryIOClient(cfg, manager)

	go func() {
		if err := manager.Run(ctx, ioClient); err != nil {
			log.Printf("ingestion manager stopped: %v", err)
		}
	}()

	return manager
}

func newFactoryIOClient(cfg config.Config, manager *ingestion.Manager) factoryio.Client {
	switch cfg.FactoryIOMode {
	case "opcua":
		nodeIDs := manager.SensorAddresses()
		if len(nodeIDs) == 0 {
			log.Println("warning: no sensors have a source_address set — opcua client will have nothing to subscribe to")
		}
		log.Printf("factory IO mode: opcua (%s), %d sensors mapped", cfg.OPCUAEndpoint, len(nodeIDs))
		return factoryio.NewOPCUAClient(cfg.OPCUAEndpoint, nodeIDs)
	default:
		log.Printf("factory IO mode: simulator")
		return factoryio.NewSimulatorClient(time.Second)
	}
}

// startEvaluation wires Phase 6: picks the AI client and runs evaluation
// in the background.
func startEvaluation(ctx context.Context, cfg config.Config, store *db.Store, manager *ingestion.Manager) *evaluation.Evaluator {
	aiClient := newAIClient(cfg)
	evaluator := evaluation.NewEvaluator(store, manager, aiClient, time.Duration(cfg.AIEvalIntervalSeconds)*time.Second)
	go evaluator.Run(ctx)
	return evaluator
}

func newAIClient(cfg config.Config) aiclient.AIClient {
	switch cfg.AIServiceMode {
	case "http":
		log.Printf("AI service mode: http (%s)", cfg.AIServiceURL)
		return aiclient.NewHTTPAIClient(cfg.AIServiceURL)
	default:
		log.Printf("AI service mode: mock")
		return aiclient.NewMockAIClient()
	}
}

// startWebSocket wires Phase 7: starts the hub and the broadcaster that
// feeds it, both in the background.
func startWebSocket(ctx context.Context, manager *ingestion.Manager, evaluator *evaluation.Evaluator) *ws.Hub {
	hub := ws.NewHub()
	go hub.Run(ctx)

	broadcaster := ws.NewBroadcaster(hub, manager, evaluator, time.Second)
	go broadcaster.Run(ctx)

	return hub
}

// runServer starts the HTTP server and blocks until ctx is cancelled
// (SIGINT/SIGTERM), then shuts down gracefully.
func runServer(ctx context.Context, cfg config.Config, router http.Handler) {
	srv := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Printf("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}