package dto

import "time"

type CreateSensorRequest struct {
	SensorID   string `json:"sensor_id" binding:"required,min=2,max=50"`
	MachineID  string `json:"machine_id" binding:"required"`
	MetricName string `json:"metric_name" binding:"required,max=50"`
	Unit       string `json:"unit" binding:"required,max=20"`
	SourceAddress *string `json:"source_address" binding:"omitempty,max=100"`
}

type SensorResponse struct {
	SensorID   string    `json:"sensor_id"`
	MachineID  string    `json:"machine_id"`
	MetricName string    `json:"metric_name"`
	Unit       string    `json:"unit"`
	SourceAddress *string   `json:"source_address,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}