package tcp

import (
	"net"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

var channelBufSize = 256
var recvBufferSize = 256000

type SessionID string

type Session struct {
	id     SessionID
	conn   net.Conn
	send   chan []byte
	server *Server
}

type SessionMap map[SessionID]*Session

func NewSession(conn net.Conn, server *Server) *Session {
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

func (c *Session) Conn() net.Conn {
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
			slog.Error("[tcp] disconnect error: %s", err.Error())
		}
		c.conn = nil
	}
}

func (c *Session) writePump() {
	defer c.Close()
	for msg := range c.send {
		// var len int
		var err error
		if _, err = c.conn.Write(msg); err != nil {
			slog.Error("[tcp] write message error: ", err)
			return
		}
	}
}

func (c *Session) readPump() {
	defer c.Close()

	buf := make([]byte, recvBufferSize)
	var err error
	var readLen int

	for {
		if readLen, err = c.conn.Read(buf); err != nil {
			slog.Error("[tcp] read message error: %v", err)
			return
		}

		if err = c.server.messageHandler(c.SessionID(), buf[:readLen]); err != nil {
			slog.Error("[tcp] process message error: %v", err)
		}
	}
}
