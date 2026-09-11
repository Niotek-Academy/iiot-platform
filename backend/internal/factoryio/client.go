package factoryio

import (
	"context"
	"time"
)


type Reading struct {
	SensorID  string
	Value     float64
	Timestamp time.Time
}

// Client is anything that can stream readings. Stream should keep sending
// on the returned channel until ctx is cancelled, then close the channel.
type Client interface {
	Stream(ctx context.Context) (<-chan Reading, error)
}