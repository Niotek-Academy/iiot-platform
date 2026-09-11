package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/config"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/factoryio"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/ingestion"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.LoadConfig()

	log.Printf("connecting to database...")
	store, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer store.Close()
	log.Printf("database connection established")

	// --- Phase 4: telemetry ingestion ---
	manager := ingestion.NewManager(store, cfg.BufferCapacity, cfg.SnapshotInterval)
	if err := manager.LoadSensorMap(ctx); err != nil {
		log.Fatalf("failed to load sensor map: %v", err)
	}

	var ioClient factoryio.Client
	switch cfg.FactoryIOMode {
	case "opcua":
		nodeIDs := manager.SensorAddresses()
		if len(nodeIDs) == 0 {
			log.Println("warning: no sensors have a source_address set — opcua client will have nothing to subscribe to")
		}
		ioClient = factoryio.NewOPCUAClient(cfg.OPCUAEndpoint, nodeIDs)
		log.Printf("factory IO mode: opcua (%s), %d sensors mapped", cfg.OPCUAEndpoint, len(nodeIDs))
	default:
		ioClient = factoryio.NewSimulatorClient(time.Second)
		log.Printf("factory IO mode: simulator")
	}

	go func() {
		if err := manager.Run(ctx, ioClient); err != nil {
			log.Printf("ingestion manager stopped: %v", err)
		}
	}()
	// --- end Phase 4 wiring ---

	router := server.NewRouter(store, cfg)

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