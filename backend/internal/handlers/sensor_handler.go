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

// Create godoc
// @Summary      Register a sensor
// @Tags         sensors
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSensorRequest true "Sensor"
// @Success      201  {object}  response.Envelope{data=dto.SensorResponse}
// @Failure      400  {object}  response.Envelope  "machine_id does not exist"
// @Failure      409  {object}  response.Envelope  "sensor_id already exists"
// @Router       /sensors [post]
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

// List godoc
// @Summary      List sensors
// @Description  Returns all sensors if machine_id is omitted, otherwise only that machine's sensors
// @Tags         sensors
// @Security     BearerAuth
// @Produce      json
// @Param        machine_id query string false "Filter by machine"
// @Success      200  {object}  response.Envelope{data=[]dto.SensorResponse}
// @Router       /sensors [get]
func (h *SensorHandler) List(c *gin.Context) {
	machineID := c.Query("machine_id") // optional now
	list, err := h.service.ListSensors(c.Request.Context(), machineID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, mapper.ToSensorResponseList(list))
}

// GetByID godoc
// @Summary      Get one sensor
// @Tags         sensors
// @Security     BearerAuth
// @Produce      json
// @Param        sensor_id path string true "Sensor ID"
// @Success      200  {object}  response.Envelope{data=dto.SensorResponse}
// @Failure      404  {object}  response.Envelope
// @Router       /sensors/{sensor_id} [get]
func (h *SensorHandler) GetByID(c *gin.Context) {
	sensorID := c.Param("sensor_id")
	sensor, err := h.service.GetByID(c.Request.Context(), sensorID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, mapper.ToSensorResponse(sensor))
}

// Delete godoc
// @Summary      Delete a sensor
// @Tags         sensors
// @Security     BearerAuth
// @Produce      json
// @Param        sensor_id path string true "Sensor ID"
// @Success      200  {object}  response.Envelope
// @Router       /sensors/{sensor_id} [delete]
func (h *SensorHandler) Delete(c *gin.Context) {
	sensorID := c.Param("sensor_id")
	if err := h.service.Delete(c.Request.Context(), sensorID); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{"deleted": sensorID})
}
