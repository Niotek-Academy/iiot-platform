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
	endpoint         string
	nodeIDs          map[string]string // sensor_id -> read NodeID
	controlAddresses map[string]string // machine_id -> control NodeID (e.g. EMERGENCY_STOP tag), from machines.control_address
}

func NewOPCUAClient(endpoint string, nodeIDs, controlAddresses map[string]string) *OPCUAClient {
	return &OPCUAClient{endpoint: endpoint, nodeIDs: nodeIDs, controlAddresses: controlAddresses}
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

// EmergencyStop and Start implement ControlClient by writing to the
// EMERGENCY_STOP digital output (see docs/1_project_overview_and_contracts.md,
// Contract A<->B). There's no separate "start" signal in the hardware
// mapping, so Start just clears the same bit.
func (c *OPCUAClient) EmergencyStop(ctx context.Context, machineID string) error {
	return c.writeControl(ctx, machineID, true)
}

func (c *OPCUAClient) Start(ctx context.Context, machineID string) error {
	return c.writeControl(ctx, machineID, false)
}

func (c *OPCUAClient) writeControl(ctx context.Context, machineID string, value bool) error {
	nodeIDStr, ok := c.controlAddresses[machineID]
	if !ok {
		return fmt.Errorf("no control_address configured for machine %s", machineID)
	}

	client, err := opcua.NewClient(c.endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err != nil {
		return fmt.Errorf("opcua: create client: %w", err)
	}
	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("opcua: connect: %w", err)
	}
	defer client.Close(ctx)

	nodeID, err := ua.ParseNodeID(nodeIDStr)
	if err != nil {
		return fmt.Errorf("opcua: invalid node id %q for machine %s: %w", nodeIDStr, machineID, err)
	}

	req := &ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{{
			NodeID:      nodeID,
			AttributeID: ua.AttributeIDValue,
			Value: &ua.DataValue{
				EncodingMask: ua.DataValueValue,
				Value:        ua.MustVariant(value),
			},
		}},
	}

	resp, err := client.Write(ctx, req)
	if err != nil {
		return fmt.Errorf("opcua: write failed: %w", err)
	}
	if len(resp.Results) == 0 || resp.Results[0] != ua.StatusOK {
		return fmt.Errorf("opcua: write returned non-OK status")
	}
	return nil
}