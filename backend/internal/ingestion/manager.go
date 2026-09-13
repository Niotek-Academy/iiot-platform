package ingestion

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/factoryio"
)

// Manager owns the per-machine state/buffers and drives both persistence
// (telemetry_logs) and the in-memory sliding window from one stream of
// readings.
type Manager struct {
	store *db.Store

	mu              sync.RWMutex
	sensorToMachine map[string]string
	sensorAddress   map[string]string // sensor_id -> driver address (OPC UA NodeID / Modbus register), non-empty only

	states  map[string]*MachineState
	buffers map[string]*RingBuffer

	bufferCapacity   int
	snapshotInterval time.Duration

	controlAddress map[string]string
}

func NewManager(store *db.Store, bufferCapacity int, snapshotInterval time.Duration) *Manager {
	return &Manager{
		store:            store,
		states:           make(map[string]*MachineState),
		buffers:          make(map[string]*RingBuffer),
		bufferCapacity:   bufferCapacity,
		snapshotInterval: snapshotInterval,
	}
}

// LoadSensorMap reads every sensor from the database once at startup and
// builds the sensor_id -> machine_id lookup used to route incoming readings.
// Call this again (e.g. periodically, or after Phase 3 admin changes) if
// sensors can be added while the service is running.
func (m *Manager) LoadSensorMap(ctx context.Context) error {
	sensors, err := m.store.ListAllSensors(ctx)
	if err != nil {
		return err
	}

	machineMap := make(map[string]string, len(sensors))
	addressMap := make(map[string]string, len(sensors))
	for _, s := range sensors {
		machineMap[s.SensorID] = s.MachineID
		if s.SourceAddress.Valid && s.SourceAddress.String != "" {
			addressMap[s.SensorID] = s.SourceAddress.String
		}
	}

	m.mu.Lock()
	m.sensorToMachine = machineMap
	m.sensorAddress = addressMap
	m.mu.Unlock()
	return nil
}


func (m *Manager) SensorAddresses() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(map[string]string, len(m.sensorAddress))
	for k, v := range m.sensorAddress {
		out[k] = v
	}
	return out
}


// Run starts consuming readings from client until ctx is cancelled. It also
// starts the background snapshot loop that builds the AI-ready sliding window.
func (m *Manager) Run(ctx context.Context, client factoryio.Client) error {
	ch, err := client.Stream(ctx)
	if err != nil {
		return err
	}

	go m.snapshotLoop(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		case reading, ok := <-ch:
			if !ok {
				return nil
			}
			m.handleReading(ctx, reading)
		}
	}
}

func (m *Manager) handleReading(ctx context.Context, r factoryio.Reading) {
	m.mu.RLock()
	machineID, known := m.sensorToMachine[r.SensorID]
	m.mu.RUnlock()

	if !known {
		log.Printf("ingestion: reading for unknown sensor_id %q, dropping", r.SensorID)
		return
	}

	// 1. Persist — permanent audit trail.
	_, err := m.store.InsertTelemetry(ctx, generated.InsertTelemetryParams{
		SensorID:    r.SensorID,
		MetricValue: r.Value,
		RecordedAt:  pgtype.Timestamptz{Time: r.Timestamp, Valid: true},
	})
	if err != nil {
		log.Printf("ingestion: failed to persist reading for %s: %v", r.SensorID, err)
		// don't return — still update in-memory state even if the DB write failed
	}

	// 2. Update in-memory state — feeds the sliding window.
	m.getOrCreateState(machineID).Update(r.SensorID, r.Value)
}

func (m *Manager) getOrCreateState(machineID string) *MachineState {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, ok := m.states[machineID]
	if !ok {
		state = NewMachineState()
		m.states[machineID] = state
		m.buffers[machineID] = NewRingBuffer(m.bufferCapacity)
	}
	return state
}

func (m *Manager) snapshotLoop(ctx context.Context) {
	ticker := time.NewTicker(m.snapshotInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			m.mu.RLock()
			for machineID, state := range m.states {
				m.buffers[machineID].Push(Snapshot{Timestamp: now, Metrics: state.Snapshot()})
			}
			m.mu.RUnlock()
		}
	}
}

// Window returns the current sliding window for one machine — this is what
// Phase 6 (AI integration) will call to build the /ai/v1/predict request.
// Returns (nil, false) if the machine has no data yet.
func (m *Manager) Window(machineID string) ([]Snapshot, bool) {
	m.mu.RLock()
	buf, ok := m.buffers[machineID]
	m.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return buf.Window(), true
}

// KnownMachines returns the machine_ids that currently have at least one
// reading (i.e. have an active state/buffer). Used by the WebSocket
// broadcaster to avoid a DB call on every tick.
func (m *Manager) KnownMachines() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]string, 0, len(m.states))
	for id := range m.states {
		out = append(out, id)
	}
	return out
}

func (m *Manager) LoadControlAddresses(ctx context.Context) error {
	machines, err := m.store.ListMachines(ctx)
	if err != nil {
		return err
	}

	addresses := make(map[string]string, len(machines))
	for _, mch := range machines {
		if mch.ControlAddress.Valid && mch.ControlAddress.String != "" {
			addresses[mch.MachineID] = mch.ControlAddress.String
		}
	}

	m.mu.Lock()
	m.controlAddress = addresses
	m.mu.Unlock()
	return nil
}

func (m *Manager) ControlAddresses() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(map[string]string, len(m.controlAddress))
	for k, v := range m.controlAddress {
		out[k] = v
	}
	return out
}