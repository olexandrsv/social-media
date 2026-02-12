package connection

import "github.com/gorilla/websocket"

type GRPCConnection struct {
	conn *websocket.Conn
}

func NewGRPC(conn *websocket.Conn) *GRPCConnection {
	return &GRPCConnection{
		conn: conn,
	}
}

func (c *GRPCConnection) Receive() ([]byte, error) {
	_, bytes, err := c.conn.ReadMessage()
	return bytes, err
}

func (c *GRPCConnection) Send(b []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, b)
}

func (c *GRPCConnection) Close() error {
	return c.conn.Close()
}
