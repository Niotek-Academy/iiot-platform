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

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// List godoc
// @Summary      List all users
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Envelope{data=[]dto.UserResponse}
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.service.List(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, mapper.ToUserResponseList(users))
}

// UpdateRole godoc
// @Summary      Change a user's role
// @Description  An admin cannot change their own role
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user_id path int true "User ID"
// @Param        request body dto.UpdateUserRoleRequest true "New role"
// @Success      200  {object}  response.Envelope{data=dto.UserResponse}
// @Failure      400  {object}  response.Envelope  "cannot change your own role"
// @Failure      404  {object}  response.Envelope
// @Router       /users/{user_id}/role [patch]
func (h *UserHandler) UpdateRole(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.Error(apperr.NewBadRequest("invalid user_id"))
		return
	}

	var req dto.UpdateUserRoleRequest
	if !BindJSON(c, &req) {
		return
	}

	actorIDVal, _ := c.Get("user_id")
	actorID, _ := actorIDVal.(int64)

	user, err := h.service.UpdateRole(c.Request.Context(), actorID, targetUserID, req.Role)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, mapper.ToUserResponse(user))
}