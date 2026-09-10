package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/mapper"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/response"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/service"
)

type SensorHandler struct {
	service *service.SensorService
}

func NewSensorHandler(s *service.SensorService) *SensorHandler {
	return &SensorHandler{service: s}
}

func (h *SensorHandler) Create(c *gin.Context) {
	var req dto.CreateSensorRequest
	if !BindJSON(c, &req) {
		return
	}
	sensor, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusCreated, mapper.ToSensorResponse(sensor))
}

func (h *SensorHandler) List(c *gin.Context) {
	machineID := c.Query("machine_id") // optional now
	list, err := h.service.ListSensors(c.Request.Context(), machineID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, mapper.ToSensorResponseList(list))
}

func (h *SensorHandler) GetByID(c *gin.Context) {
	sensorID := c.Param("sensor_id")
	sensor, err := h.service.GetByID(c.Request.Context(), sensorID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, mapper.ToSensorResponse(sensor))
}

func (h *SensorHandler) Delete(c *gin.Context) {
	sensorID := c.Param("sensor_id")
	if err := h.service.Delete(c.Request.Context(), sensorID); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{"deleted": sensorID})
}
