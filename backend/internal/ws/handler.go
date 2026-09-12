package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/response"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/utils/jwtutil"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Dev-only: accepts connections from any origin. Restrict this to your
	// actual frontend's origin before going to production (Phase 10).
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWS handles GET /ws/v1/factory-stream. Auth is via ?token=<jwt> since
// browsers can't set an Authorization header on a WebSocket handshake.
// Optional ?machine_id=<id> filters the stream to one machine.
func ServeWS(hub *Hub, jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			response.Fail(c, apperr.NewUnauthorized("missing token query parameter"))
			return
		}
		if _, err := jwtutil.Parse(jwtSecret, token); err != nil {
			response.Fail(c, apperr.NewUnauthorized("invalid or expired token"))
			return
		}

		machineFilter := c.Query("machine_id")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("ws: upgrade failed: %v", err)
			return
		}

		client := newClient(conn, machineFilter)
		hub.register <- client

		go client.writePump()
		go client.readPump(hub)
	}
}