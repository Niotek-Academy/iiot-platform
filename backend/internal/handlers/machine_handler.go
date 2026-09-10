package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/mapper"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/response"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/service"
)

type MachineHandler struct {
	service *service.MachineService
}

func NewMachineHandler(s *service.MachineService) *MachineHandler {
	return &MachineHandler{service: s}
}

func (h *MachineHandler) Create(c *gin.Context) {
	var req dto.CreateMachineRequest
	if !BindJSON(c, &req) {
		return
	}
	machine, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusCreated, mapper.ToMachineResponse(machine))
}

func (h *MachineHandler) List(c *gin.Context) {
	list, err := h.service.List(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, mapper.ToMachineResponseList(list))
}

func (h *MachineHandler) Update(c *gin.Context) {
	machineID := c.Param("machine_id")
	var req dto.UpdateMachineRequest
	if !BindJSON(c, &req) {
		return
	}
	machine, err := h.service.Update(c.Request.Context(), machineID, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, mapper.ToMachineResponse(machine))
}

func (h *MachineHandler) Delete(c *gin.Context) {
	machineID := c.Param("machine_id")
	if err := h.service.Delete(c.Request.Context(), machineID); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{"deleted": machineID})
}