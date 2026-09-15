package mapper

import (
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
)

func ToCommandLogResponse(c generated.CommandLog) dto.CommandLogResponse {
	return dto.CommandLogResponse{
		CommandID:   c.CommandID,
		CommandType: c.CommandType,
		IssuedBy:    c.IssuedBy,
		Reason:      c.Reason,
		ExecutedAt:  c.ExecutedAt.Time,
	}
}

func ToCommandLogResponseList(list []generated.CommandLog) []dto.CommandLogResponse {
	out := make([]dto.CommandLogResponse, 0, len(list))
	for _, c := range list {
		out = append(out, ToCommandLogResponse(c))
	}
	return out
}