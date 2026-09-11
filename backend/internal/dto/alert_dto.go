package dto

import "time"

type AlertResponse struct {
	AlertID    int64     `json:"alert_id"`
	Severity   string    `json:"severity"`
	Message    string    `json:"message"`
	IsResolved bool      `json:"is_resolved"`
	CreatedAt  time.Time `json:"created_at"`
}

type AlertsResponse struct {
	MachineID string          `json:"machine_id"`
	Alerts    []AlertResponse `json:"alerts"`
}