package dto

import "time"

type CreateMachineRequest struct {
	MachineID string `json:"machine_id" binding:"required,min=3,max=50"`
	Name      string `json:"name" binding:"required,max=100"`
	Location  string `json:"location" binding:"required,max=100"`
	Status    string `json:"status" binding:"omitempty,oneof=OPERATIONAL WARNING CRITICAL STOPPED"`
}

// Pointer fields distinguish "not sent" (nil) from "sent as empty string" —
// required for a PATCH-style partial update.
type UpdateMachineRequest struct {
	Name     *string `json:"name" binding:"omitempty,max=100"`
	Location *string `json:"location" binding:"omitempty,max=100"`
	Status   *string `json:"status" binding:"omitempty,oneof=OPERATIONAL WARNING CRITICAL STOPPED"`
}

type MachineResponse struct {
	MachineID string    `json:"machine_id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type MachineOverviewResponse struct {
	MachineID          string    `json:"machine_id"`
	Name               string    `json:"name"`
	Location           string    `json:"location"`
	Status             string    `json:"status"`
	CurrentHealthScore *float64  `json:"current_health_score"` 
	RULHours           *float64  `json:"rul_hours"`            
	ActiveAlertsCount  int64     `json:"active_alerts_count"`
	UpdatedAt          time.Time `json:"updated_at"`
}