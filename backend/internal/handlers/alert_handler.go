package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/mapper"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/response"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/service"
)

type AlertHandler struct {
	service *service.AlertService
}

func NewAlertHandler(s *service.AlertService) *AlertHandler {
	return &AlertHandler{service: s}
}

func (h *AlertHandler) List(c *gin.Context) {
	machineID := c.Query("machine_id")
	if machineID == "" {
		c.Error(apperr.NewBadRequest("machine_id query parameter is required"))
		return
	}

	var isResolved *bool
	if raw := c.Query("is_resolved"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			c.Error(apperr.NewBadRequest("is_resolved must be true or false"))
			return
		}
		isResolved = &parsed
	}

	alerts, err := h.service.ListByMachine(c.Request.Context(), machineID, isResolved)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, dto.AlertsResponse{
		MachineID: machineID,
		Alerts:    mapper.ToAlertResponseList(alerts),
	})
}