package handlers

import (
	"net/http"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/mapper"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/response"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/service"
	"github.com/gin-gonic/gin"
)



type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Login godoc
// @Summary      Log in
// @Description  Returns a JWT access token (12h expiry, no refresh)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Credentials"
// @Success      200  {object}  response.Envelope{data=dto.AuthResponse}
// @Failure      401  {object}  response.Envelope
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context)  {
	var req dto.LoginRequest
	if !BindJSON(c, &req) {
		return
	}
	token, user, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		dto.AuthResponse{
			Token: token,
			User:  mapper.ToUserResponse(user),
		},
	)
}

// Register godoc
// @Summary      Register a new user
// @Description  First user in the system becomes ADMIN automatically (no token needed). Every user after that requires an ADMIN token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "New user"
// @Success      201  {object}  response.Envelope{data=dto.UserResponse}
// @Failure      403  {object}  response.Envelope
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !BindJSON(c, &req) {
		return
	}

	// Set by middleware.OptionalAuth if a valid token was presented.
	// Empty string means "no authenticated actor" (the bootstrap case).
	actorRole, _ := c.Get("role")
	actorRoleStr, _ := actorRole.(string)

	user, err := h.service.Register(c.Request.Context(), actorRoleStr, req.Username, req.Password, req.Role)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, mapper.ToUserResponse(user))
}