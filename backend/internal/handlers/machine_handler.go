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

// Create godoc
// @Summary      Register a machine
// @Tags         machines
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateMachineRequest true "Machine"
// @Success      201  {object}  response.Envelope{data=dto.MachineResponse}
// @Failure      409  {object}  response.Envelope
// @Router       /machines [post]
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

// List godoc
// @Summary      List all machines
// @Tags         machines
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Envelope{data=[]dto.MachineResponse}
// @Router       /machines [get]
func (h *MachineHandler) List(c *gin.Context) {
	list, err := h.service.List(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, mapper.ToMachineResponseList(list))
}

// Update godoc
// @Summary      Partially update a machine
// @Tags         machines
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        machine_id path string true "Machine ID"
// @Param        request body dto.UpdateMachineRequest true "Fields to change"
// @Success      200  {object}  response.Envelope{data=dto.MachineResponse}
// @Failure      404  {object}  response.Envelope
// @Router       /machines/{machine_id} [patch]
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

// Delete godoc
// @Summary      Delete a machine
// @Tags         machines
// @Security     BearerAuth
// @Produce      json
// @Param        machine_id path string true "Machine ID"
// @Success      200  {object}  response.Envelope
// @Failure      404  {object}  response.Envelope
// @Router       /machines/{machine_id} [delete]
func (h *MachineHandler) Delete(c *gin.Context) {
	machineID := c.Param("machine_id")
	if err := h.service.Delete(c.Request.Context(), machineID); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{"deleted": machineID})
}

// GetOverview godoc
// @Summary      Machine status + health summary
// @Description  health_score/rul_hours 
// @Tags         machines
// @Security     BearerAuth
// @Produce      json
// @Param        machine_id path string true "Machine ID"
// @Success      200  {object}  response.Envelope{data=dto.MachineOverviewResponse}
// @Failure      404  {object}  response.Envelope
// @Router       /machines/{machine_id} [get]
func (h *MachineHandler) GetOverview(c *gin.Context) {
	machineID := c.Param("machine_id")
	overview, err := h.service.GetOverview(c.Request.Context(), machineID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, overview)
}