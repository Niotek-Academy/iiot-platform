package server

import (
	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/config"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/handlers"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/middleware"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/service"
)

func NewRouter(store *db.Store, cfg config.Config) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler())

	jwtSecret := []byte(cfg.JWTSecret)

	authService := service.NewAuthService(store, cfg.JWTSecret, cfg.JWTExpiryHours)
	authHandler := handlers.NewAuthHandler(authService)

	r.GET("/healthz", handlers.HealthCheck(store))

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			// OptionalAuth, not Auth: must stay reachable with no token at
			// all so the very first admin can be created (see AuthService.Register).
			auth.POST("/register", middleware.OptionalAuth(jwtSecret), authHandler.Register)
		}
	}

	return r
}