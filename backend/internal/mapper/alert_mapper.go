package mapper

import (
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
)

func ToAlertResponse(a generated.Alert) dto.AlertResponse {
	return dto.AlertResponse{
		AlertID:    a.AlertID,
		Severity:   a.Severity,
		Message:    a.Message,
		IsResolved: a.IsResolved,
		CreatedAt:  a.CreatedAt.Time,
	}
}

func ToAlertResponseList(list []generated.Alert) []dto.AlertResponse {
	out := make([]dto.AlertResponse, 0, len(list))
	for _, a := range list {
		out = append(out, ToAlertResponse(a))
	}
	return out
}