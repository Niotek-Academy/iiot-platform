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