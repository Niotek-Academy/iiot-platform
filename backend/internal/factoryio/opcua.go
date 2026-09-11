package factoryio

import (
	"context"
	"fmt"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)


type OPCUAClient struct {
	endpoint string 
	nodeIDs  map[string]string // sensor_id -> OPC UA NodeID string, e.g. "ns=2;s=Channel1.LT_01"
}

func NewOPCUAClient(endpoint string, nodeIDs map[string]string) *OPCUAClient {
	return &OPCUAClient{endpoint: endpoint, nodeIDs: nodeIDs}
}

// Stream implements factoryio.Client.
func (c *OPCUAClient) Stream(ctx context.Context) (<-chan Reading, error) {
	client, err := opcua.NewClient(c.endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err != nil {
		return nil, fmt.Errorf("opcua: create client: %w", err)
	}
	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("opcua: connect to %s: %w", c.endpoint, err)
	}

	m, err := monitor.NewNodeMonitor(client)
	if err != nil {
		return nil, fmt.Errorf("opcua: create node monitor: %w", err)
	}

	out := make(chan Reading)

	// reverse map: NodeID string -> our sensor_id, so we can label incoming
	// values correctly.
	nodeToSensor := make(map[string]string, len(c.nodeIDs))
	nodeIDList := make([]string, 0, len(c.nodeIDs))
	for sensorID, nodeID := range c.nodeIDs {
		nodeToSensor[nodeID] = sensorID
		nodeIDList = append(nodeIDList, nodeID)
	}

	sub, err := m.Subscribe(ctx, &opcua.SubscriptionParameters{Interval: time.Second}, func(s *monitor.Subscription, msg *monitor.DataChangeMessage) {
		if msg.Error != nil {
			return
		}
		sensorID, ok := nodeToSensor[msg.NodeID.String()]
		if !ok {
			return
		}
		value, ok := toFloat64(msg.Value.Value())
		if !ok {
			return
		}
		select {
		case out <- Reading{SensorID: sensorID, Value: value, Timestamp: time.Now()}:
		case <-ctx.Done():
		}
	}, nodeIDList...)
	if err != nil {
		return nil, fmt.Errorf("opcua: subscribe: %w", err)
	}

	go func() {
		<-ctx.Done()
		_ = sub.Unsubscribe(context.Background())
		_ = client.Close(context.Background())
		close(out)
	}()

	return out, nil
}

// toFloat64 converts the handful of OPC UA value types we expect (float,
// double, int, bool) into a plain float64 for our Reading.
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	case bool:
		if n {
			return 1, true
		}
		return 0, true
	default:
		return 0, false
	}
}