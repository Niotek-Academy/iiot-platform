package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/response"
)

func HealthCheck(store *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := store.Pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, response.Envelope{
				Success: false,
				Error:   &response.ErrorBody{Code: "DATABASE_UNREACHABLE", Message: "database is unreachable"},
			})
			return
		}

		response.Success(c, http.StatusOK, gin.H{"database": "connected"})
	}
}