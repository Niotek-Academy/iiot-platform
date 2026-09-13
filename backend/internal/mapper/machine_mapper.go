package mapper

import (
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
)

func ToMachineResponse(m generated.Machine) dto.MachineResponse {
	resp := dto.MachineResponse{
		MachineID: m.MachineID,
		Name:      m.Name,
		Location:  m.Location,
		Status:    m.Status,
		CreatedAt: m.CreatedAt.Time,
	}
	if m.ControlAddress.Valid {
		resp.ControlAddress = &m.ControlAddress.String
	}
	return resp
}

func ToMachineResponseList(list []generated.Machine) []dto.MachineResponse {
	out := make([]dto.MachineResponse, 0, len(list))
	for _, m := range list {
		out = append(out, ToMachineResponse(m))
	}
	return out
}