// WebSocket Example

package wse

import (
	// "context"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PING_PERIOD      = (pongWait * 9) / 10
	BROADCAST_PERIOD = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type WebSocketHandler struct {
	Upgrader websocket.Upgrader
}

type PortReceiverInterface[T any] interface {
	Receiver() chan T
}
type PortTransmitterInterface[T any] interface {
	Transmitter() chan T
}

type HandlerPortReceiver[T any] struct {
	rx chan T
}

func NewHandlerPortReceiver[T any](size int) *HandlerPortReceiver[T] {
	return &HandlerPortReceiver[T]{rx: make(chan T, size)}
}

func (h *HandlerPortReceiver[T]) Receiver() chan T {
	return h.rx
}

type AdapterPortTransmitter[T any] struct {
	tx chan T
}

func NewAdapterPortTransmitter[T any](size int) *AdapterPortTransmitter[T] {
	return &AdapterPortTransmitter[T]{tx: make(chan T)}
}
func (h *AdapterPortTransmitter[T]) Transmitter() chan T {
	return h.tx
}

type Register[T any] struct {
	clientId string
	port     PortTransmitterInterface[T]
}
type Message struct {
	clientId string
	text     string
}

type ClientController struct {
	// The websocket connection.
	conn *websocket.Conn

	portUnregister PortReceiverInterface[Register[Message]]
	portRegister   PortReceiverInterface[Register[Message]]
	portInput      PortReceiverInterface[Message]
	portOutput     PortTransmitterInterface[Message]
}

func (c *ClientController) doReading() {
	clientId := fmt.Sprintf("%p", c)
	c.portRegister.Receiver() <- Register[Message]{
		clientId: clientId,
		port:     c.portOutput}

	defer func() {
		c.portUnregister.Receiver() <- Register[Message]{
			clientId: clientId,
			port:     c.portOutput}
		c.conn.Close()

	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Time{})
	//!!! c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//!!! c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v\n", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.portInput.Receiver() <- Message{clientId: clientId, text: string(message)}
	}
}

func (c *ClientController) doWriting() {
	ticker := time.NewTicker(PING_PERIOD)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case msg, ok := <-c.portOutput.Transmitter():
			if !ok {
				log.Printf("Closing websocket\n")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Next writter error: %s\n", err)
				return
			}

			log.Printf("Sending from %s, msg:%s\n", msg.clientId, msg.text)
			if _, err := w.Write([]byte(msg.text)); err != nil {
				log.Println(err)
			}
			w.Close()
			// log.Printf("Sending done\n")
			break
		case <-ticker.C:
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Ping message error\n")

				return
			}
		}
	}

}

type ClientMiddleware struct {
}
type Core struct {
	portUnregister HandlerPortReceiver[Register[Message]]
	portRegister   HandlerPortReceiver[Register[Message]]
	portInput      HandlerPortReceiver[Message]
	portOutput     map[string]PortTransmitterInterface[Message]
}

func (c *Core) Run() {
	ticker := time.NewTicker(BROADCAST_PERIOD)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case msg := <-c.portInput.Receiver():
			log.Printf("Received from %s, msg:%s\n", msg.clientId, msg.text)
			text := fmt.Sprintf("%s, you too", msg.text)
			c.portOutput[msg.clientId].Transmitter() <- Message{clientId: msg.clientId, text: text}
			// _ = msg
			break
		case reg := <-c.portRegister.Receiver():
			log.Printf("Registering %s\n", reg.clientId)
			c.portOutput[reg.clientId] = reg.port
			break
		case reg := <-c.portUnregister.Receiver():
			log.Printf("Unregistering %s\n", reg.clientId)
			if _, ok := c.portOutput[reg.clientId]; ok {
				// Close channel
				close(reg.port.Transmitter())
				// Delete from the map
				delete(c.portOutput, reg.clientId)
			}
			break
		case <-ticker.C:
			log.Printf("Tick\n")
			for clientId, port := range c.portOutput {
				port.Transmitter() <- Message{
					clientId: clientId,
					text:     fmt.Sprintf("Whats up %s", clientId)}
			}
			break
		}
	}
}
func NewCore() *Core {
	return &Core{
		portUnregister: *NewHandlerPortReceiver[Register[Message]](100),
		portRegister:   *NewHandlerPortReceiver[Register[Message]](100),
		portInput:      *NewHandlerPortReceiver[Message](100),
		portOutput:     make(map[string]PortTransmitterInterface[Message]),
	}
}

func ServeWs(core *Core, w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// core := NewCore()

	client := ClientController{
		conn:           conn,
		portUnregister: &core.portUnregister,
		portRegister:   &core.portRegister,
		portInput:      &core.portInput,
		portOutput:     NewAdapterPortTransmitter[Message](10)}

	// portOutput:     &AdapterPortTransmitter[Message]{tx: make(chan Message)}}
	// client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	// client.hub.register <- client

	// // Allow collection of memory referenced by the caller by doing all work in
	// // new goroutines.
	// go client.writePump()
	// go client.readPump()
	go client.doReading()
	go client.doWriting()
	log.Printf("Ready\n")
}
