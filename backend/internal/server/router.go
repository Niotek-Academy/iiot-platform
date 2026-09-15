package server

import (
	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/config"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/factoryio"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/handlers"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/middleware"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/service"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/ws"

	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Niotek-Academy/iiot-platform/backend/docs"
)

func NewRouter(store *db.Store, cfg config.Config, hub *ws.Hub, ioClient factoryio.ControlClient) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler())

	jwtSecret := []byte(cfg.JWTSecret)

	authService := service.NewAuthService(store, cfg.JWTSecret, cfg.JWTExpiryHours)
	authHandler := handlers.NewAuthHandler(authService)

	machineService := service.NewMachineService(store)
	machineHandler := handlers.NewMachineHandler(machineService)

	sensorService := service.NewSensorService(store)
	sensorHandler := handlers.NewSensorHandler(sensorService)

	alertService := service.NewAlertService(store)
	alertHandler := handlers.NewAlertHandler(alertService)

	telemetryService := service.NewTelemetryService(store)
	telemetryHandler := handlers.NewTelemetryHandler(telemetryService)

	commandService := service.NewCommandService(store, ioClient) 
	commandHandler := handlers.NewCommandHandler(commandService)

	userService := service.NewUserService(store)
	userHandler := handlers.NewUserHandler(userService)

	r.GET("/healthz", handlers.HealthCheck(store))
	r.GET("/ws/v1/factory-stream", ws.ServeWS(hub, jwtSecret))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			// OptionalAuth, not Auth: must stay reachable with no token at
			// all so the very first admin can be created (see AuthService.Register).
			auth.POST("/register", middleware.OptionalAuth(jwtSecret), authHandler.Register)

			admin := api.Group("")
			admin.Use(middleware.Auth(jwtSecret))
			admin.Use(middleware.RequireRole("ADMIN")) // Authrize Admin role
			{
				machines := admin.Group("/machines")
				machines.POST("", machineHandler.Create)
				machines.GET("", machineHandler.List)
				machines.PATCH("/:machine_id", machineHandler.Update)
				machines.DELETE("/:machine_id", machineHandler.Delete)

				sensors := admin.Group("/sensors")
				sensors.POST("", sensorHandler.Create)
				sensors.GET("/:sensor_id", sensorHandler.GetByID)
				sensors.GET("", sensorHandler.List) // ?machine_id=...
				sensors.DELETE("/:sensor_id", sensorHandler.Delete)

				users := admin.Group("/users")
				users.GET("", userHandler.List)
				users.PATCH("/:user_id/role", userHandler.UpdateRole)
			}
			
			authed := api.Group("")
			authed.Use(middleware.Auth(jwtSecret))
			{
				authed.GET("/machines/:machine_id", machineHandler.GetOverview)
				authed.GET("/telemetry", telemetryHandler.GetHistory)
				authed.GET("/alerts", alertHandler.List)
				authed.POST("/machines/:machine_id/command", commandHandler.Execute)
				authed.GET("/machines/:machine_id/commands", commandHandler.GetHistory)
			}
		}
	}

	return r
}