package factoryio

import (
	"context"
	"math/rand"
	"time"
)


type simSensor struct {
	id       string
	value    float64
	min, max float64
	step     float64 // max change per tick, before clamping to [min, max]
}

// SimulatorClient generates realistic-looking values for all 10 sensors
type SimulatorClient struct {
	sensors  []*simSensor
	interval time.Duration
	rng      *rand.Rand
}

func NewSimulatorClient(interval time.Duration) *SimulatorClient {
	return &SimulatorClient{
		interval: interval,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		sensors: []*simSensor{
			{id: "LT_01", value: 200, min: 0, max: 300, step: 4},
			{id: "LT_02", value: 180, min: 0, max: 300, step: 4},
			{id: "LT_03", value: 210, min: 0, max: 300, step: 4},
			{id: "TE_01", value: 60, min: 0, max: 100, step: 1.5},
			{id: "LS_HIGH", value: 0, min: 0, max: 1, step: 1}, // effectively 0 or 1
			{id: "RPM_01", value: 1200, min: 0, max: 1500, step: 20},
			{id: "CURR_01", value: 14, min: 0, max: 25, step: 1},
			{id: "VIB_01", value: 10, min: 0, max: 50, step: 1.5},
			{id: "VIS_01", value: 0, min: 0, max: 1, step: 1}, // 0=OK, 1=DEFECT — see note below
			{id: "POS_01", value: 2.5, min: 0, max: 5, step: 0.3},
		},
	}
}

// Stream implements factoryio.Client.
func (c *SimulatorClient) Stream(ctx context.Context) (<-chan Reading, error) {
	out := make(chan Reading)

	go func() {
		defer close(out)
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case now := <-ticker.C:
				for _, s := range c.sensors {
					s.tick(c.rng)
					reading := Reading{SensorID: s.id, Value: s.value, Timestamp: now}
					select {
					case out <- reading:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return out, nil
}

func (s *simSensor) tick(rng *rand.Rand) {
	// LS_HIGH and VIS_01 are boolean-ish signals, not continuous ones —
	// mostly stay at 0, rarely flip to 1, instead of random-walking.
	if s.id == "LS_HIGH" || s.id == "VIS_01" {
		if rng.Float64() < 0.02 { // ~2% chance per tick
			s.value = 1
		} else {
			s.value = 0
		}
		return
	}

	delta := (rng.Float64()*2 - 1) * s.step // random value in [-step, +step]
	s.value += delta
	if s.value < s.min {
		s.value = s.min
	}
	if s.value > s.max {
		s.value = s.max
	}
}