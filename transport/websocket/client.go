package websocket

import (
	"errors"
	"net/url"
	"time"

	"github.com/fzf-labs/fpkg/encoding"
	ws "github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

type ClientMessageHandler func(MessagePayload) error

type ClientHandlerData struct {
	Handler ClientMessageHandler
	Binder  Binder
}
type ClientMessageHandlerMap map[MessageType]ClientHandlerData

type Client struct {
	conn *ws.Conn

	url      string
	endpoint *url.URL

	codec           encoding.Codec
	messageHandlers ClientMessageHandlerMap

	timeout time.Duration
}

func NewClient(opts ...ClientOption) *Client {
	cli := &Client{
		url:             "",
		timeout:         1 * time.Second,
		codec:           encoding.GetCodec("json"),
		messageHandlers: make(ClientMessageHandlerMap),
	}

	cli.init(opts...)

	return cli
}

func (c *Client) init(opts ...ClientOption) {
	for _, o := range opts {
		o(c)
	}

	c.endpoint, _ = url.Parse(c.url)
}

func (c *Client) Connect() error {
	if c.endpoint == nil {
		return errors.New("endpoint is nil")
	}

	slog.Info("[websocket] connecting to %s", c.endpoint.String())

	conn, resp, err := ws.DefaultDialer.Dial(c.endpoint.String(), nil)
	if err != nil {
		slog.Error("%s [%v]", err.Error(), resp)
		return err
	}
	c.conn = conn

	go c.run()

	return nil
}

func (c *Client) Disconnect() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			slog.Error("[websocket] disconnect error: %s", err.Error())
		}
		c.conn = nil
	}
}

func (c *Client) RegisterMessageHandler(messageType MessageType, handler ClientMessageHandler, binder Binder) {
	if _, ok := c.messageHandlers[messageType]; ok {
		return
	}

	c.messageHandlers[messageType] = ClientHandlerData{handler, binder}
}

func (c *Client) DeregisterMessageHandler(messageType MessageType) {
	delete(c.messageHandlers, messageType)
}

func (c *Client) SendMessage(messageType int, message interface{}) error {
	var msg Message
	msg.Type = MessageType(messageType)
	msg.Body, _ = encoding.Marshal(c.codec, message)

	buff, err := msg.Marshal()
	if err != nil {
		return err
	}

	if err := c.SendBinaryMessage(buff); err != nil {
		return err
	}

	return nil
}

func (c *Client) SendPingMessage(message string) error {
	return c.conn.WriteMessage(ws.PingMessage, []byte(message))
}

func (c *Client) SendPongMessage(message string) error {
	return c.conn.WriteMessage(ws.PongMessage, []byte(message))
}

func (c *Client) SendTextMessage(message string) error {
	return c.conn.WriteMessage(ws.TextMessage, []byte(message))
}

func (c *Client) SendBinaryMessage(message []byte) error {
	return c.conn.WriteMessage(ws.BinaryMessage, message)
}

func (c *Client) run() {
	defer c.Disconnect()

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
			_ = c.messageHandler(data)
			return
		case ws.TextMessage:
			slog.Error("[websocket] not support text message")
			return
		case ws.PingMessage:
			if err := c.SendPongMessage(""); err != nil {
				slog.Error("[websocket] write pong message error: ", err)
				return
			}
			return
		case ws.PongMessage:
			return
		}
	}
}

func (c *Client) messageHandler(buf []byte) error {
	var msg Message
	if err := msg.Unmarshal(buf); err != nil {
		slog.Error("[websocket] decode message exception: %s", err)
		return err
	}

	handlerData, ok := c.messageHandlers[msg.Type]
	if !ok {
		slog.Error("[websocket] message type not found:", msg.Type)
		return errors.New("message handler not found")
	}

	var payload MessagePayload

	if handlerData.Binder != nil {
		payload = handlerData.Binder()
	} else {
		payload = msg.Body
	}

	if err := encoding.Unmarshal(c.codec, msg.Body, &payload); err != nil {
		slog.Error("[websocket] unmarshal message exception: %s", err)
		return err
	}

	if err := handlerData.Handler(payload); err != nil {
		slog.Error("[websocket] message handler exception: %s", err)
		return err
	}

	return nil
}
