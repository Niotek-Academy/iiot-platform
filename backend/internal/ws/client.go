package ws

import "github.com/gorilla/websocket"

type websocketConn = websocket.Conn

func newClient(conn *websocket.Conn, machineFilter string) *Client {
	return &Client{
		conn:          conn,
		send:          make(chan []byte, 16),
		machineFilter: machineFilter,
	}
}

// writePump forwards everything the Hub sends this client's way out over
// the socket. Must run in its own goroutine.
func (c *Client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
	_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

// readPump's only job is to detect the client disconnecting (browser tab
// closed, network drop) and clean up. We don't expect the client to send
// us anything meaningful.
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}