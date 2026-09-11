package ingestion

import "sync"

// MachineState holds the most recently seen value for each sensor on one
// machine. Readings arrive one at a time and asynchronously; this is what
// lets us combine them into a single Snapshot once per tick.
type MachineState struct {
	mu     sync.RWMutex
	values map[string]float64
}

func NewMachineState() *MachineState {
	return &MachineState{values: make(map[string]float64)}
}

func (s *MachineState) Update(sensorID string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[sensorID] = value
}

// Snapshot returns a copy of the current values, safe to store in a
// RingBuffer without worrying about later mutation.
func (s *MachineState) Snapshot() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string]float64, len(s.values))
	for k, v := range s.values {
		out[k] = v
	}
	return out
}