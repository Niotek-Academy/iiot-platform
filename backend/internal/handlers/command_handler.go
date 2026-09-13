package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/response"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/service"
)

type CommandHandler struct {
	service *service.CommandService
}

func NewCommandHandler(s *service.CommandService) *CommandHandler {
	return &CommandHandler{service: s}
}

func (h *CommandHandler) Execute(c *gin.Context) {
	machineID := c.Param("machine_id")

	var req dto.CommandRequest
	if !BindJSON(c, &req) {
		return
	}

	// issued_by comes from the authenticated token, never from the request
	// body — see middleware.Auth, which sets this in the context.
	usernameVal, _ := c.Get("username")
	username, _ := usernameVal.(string)

	logEntry, err := h.service.Execute(c.Request.Context(), machineID, req.Command, req.Reason, username)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, dto.CommandResponse{
		CommandID:  logEntry.CommandID,
		MachineID:  logEntry.MachineID,
		Command:    logEntry.CommandType,
		Status:     "SUCCESSFULLY_EXECUTED",
		ExecutedAt: logEntry.ExecutedAt.Time,
	})
}