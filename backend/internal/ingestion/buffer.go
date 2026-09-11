package ingestion

import (
	"sync"
	"time"
)


type Snapshot struct {
	Timestamp time.Time
	Metrics   map[string]float64
}

// RingBuffer holds the last `capacity` snapshots for one machine.
type RingBuffer struct {
	mu        sync.Mutex
	snapshots []Snapshot
	capacity  int
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		snapshots: make([]Snapshot, 0, capacity),
		capacity:  capacity,
	}
}

// Push adds a snapshot, dropping the oldest one if we're at capacity.
func (b *RingBuffer) Push(s Snapshot) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.snapshots = append(b.snapshots, s)
	if len(b.snapshots) > b.capacity {
		b.snapshots = b.snapshots[len(b.snapshots)-b.capacity:]
	}
}

// Window returns a copy of everything currently in the buffer, oldest
// first — safe to hand off to another goroutine (e.g. building an AI request).
func (b *RingBuffer) Window() []Snapshot {
	b.mu.Lock()
	defer b.mu.Unlock()

	out := make([]Snapshot, len(b.snapshots))
	copy(out, b.snapshots)
	return out
}