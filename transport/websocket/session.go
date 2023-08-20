package websocket

import (
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

var channelBufSize = 256

type SessionID string

type Session struct {
	id     SessionID
	conn   *ws.Conn
	send   chan []byte
	server *Server
}

type SessionMap map[SessionID]*Session

func NewSession(conn *ws.Conn, server *Server) *Session {
	if conn == nil {
		panic("conn cannot be nil")
	}

	u1, _ := uuid.NewUUID()

	c := &Session{
		id:     SessionID(u1.String()),
		conn:   conn,
		send:   make(chan []byte, channelBufSize),
		server: server,
	}

	return c
}

func (c *Session) Conn() *ws.Conn {
	return c.conn
}

func (c *Session) SessionID() SessionID {
	return c.id
}

func (c *Session) SendMessage(message []byte) {
	c.send <- message
}

func (c *Session) Close() {
	c.server.unregister <- c
	c.closeConnect()
}

func (c *Session) Listen() {
	go c.writePump()
	go c.readPump()
}

func (c *Session) closeConnect() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			slog.Error("[websocket] disconnect error: %s", err.Error())
		}
		c.conn = nil
	}
}

func (c *Session) SendPingMessage(message string) error {
	return c.conn.WriteMessage(ws.PingMessage, []byte(message))
}

func (c *Session) SendPongMessage(message string) error {
	return c.conn.WriteMessage(ws.PongMessage, []byte(message))
}

func (c *Session) SendTextMessage(message string) error {
	return c.conn.WriteMessage(ws.TextMessage, []byte(message))
}

func (c *Session) SendBinaryMessage(message []byte) error {
	return c.conn.WriteMessage(ws.BinaryMessage, message)
}

func (c *Session) writePump() {
	defer c.Close()
	for msg := range c.send {
		if err := c.SendBinaryMessage(msg); err != nil {
			slog.Error("[websocket] write message error: ", err)
			return
		}
	}
}

func (c *Session) readPump() {
	defer c.Close()

	for {
		messageType, data, err := c.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseNormalClosure, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				slog.Error("[websocket] read message error: %v", err)
			}
			return
		}

		switch messageType {
		case ws.CloseMessage:
			return
		case ws.BinaryMessage:
			_ = c.server.messageHandler(c.SessionID(), data)
			return
		case ws.PingMessage:
			if err := c.SendPongMessage(""); err != nil {
				slog.Error("[websocket] write pong message error: ", err)
				return
			}
			return
		case ws.PongMessage:
			return
		case ws.TextMessage:
			slog.Error("[websocket] not support text message")
			return
		}
	}
}
