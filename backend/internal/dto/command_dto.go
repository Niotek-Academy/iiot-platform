package dto

import "time"

type CommandRequest struct {
	Command string `json:"command" binding:"required,oneof=EMERGENCY_STOP START RESET_ALERTS"`
	Reason  string `json:"reason" binding:"required,max=255"`
}

type CommandResponse struct {
	CommandID  int64     `json:"command_id"`
	MachineID  string    `json:"machine_id"`
	Command    string    `json:"command"`
	Status     string    `json:"status"`
	ExecutedAt time.Time `json:"executed_at"`
}