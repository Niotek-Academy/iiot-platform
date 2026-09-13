package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/response"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/service"
)

const (
	defaultTelemetryLimit = 50
	maxTelemetryLimit     = 1000
)

type TelemetryHandler struct {
	service *service.TelemetryService
}

func NewTelemetryHandler(s *service.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{service: s}
}

// GetHistory godoc
// @Summary      Historical telemetry for one sensor
// @Tags         telemetry
// @Security     BearerAuth
// @Produce      json
// @Param        sensor_id query string true  "Sensor ID"
// @Param        limit     query int    false "Max data points (default 50, max 1000)"
// @Success      200  {object}  response.Envelope{data=dto.TelemetryResponse}
// @Failure      404  {object}  response.Envelope
// @Router       /telemetry [get]
func (h *TelemetryHandler) GetHistory(c *gin.Context) {
	sensorID := c.Query("sensor_id")
	if sensorID == "" {
		c.Error(apperr.NewBadRequest("sensor_id query parameter is required"))
		return
	}

	limit := defaultTelemetryLimit
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			c.Error(apperr.NewBadRequest("limit must be a positive integer"))
			return
		}
		if parsed > maxTelemetryLimit {
			parsed = maxTelemetryLimit
		}
		limit = parsed
	}

	result, err := h.service.GetHistory(c.Request.Context(), sensorID, int32(limit))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, result)
}