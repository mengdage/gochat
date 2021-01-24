package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type UserLogin struct {
	UserID string
	Client *Client
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Addr          string
	Conn          *websocket.Conn
	Send          chan []byte
	UserID        string
	FirstTime     int64
	HeartbeatTime int64
	LoginTime     int64
}

// NewClient create a new client.
func NewClient(addr string, conn *websocket.Conn, t int64) *Client {
	return &Client{
		Addr:          addr,
		Conn:          conn,
		Send:          make(chan []byte, 100),
		FirstTime:     t,
		HeartbeatTime: t,
	}
}

// SendMessage sends a message to the client.
func (c *Client) SendMessage(message []byte) {
	if c == nil {
		return
	}

	c.Send <- message
}

// read listens on the websocket and process each message.
func (c *Client) read() {
	defer func() {
		log.Println("Closing send channel")
		close(c.Send)

		clientManager.Unregister <- c
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(d string) error {
		log.Printf("[Pong Handler] %s", d)
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		msgType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Unexpected error: %v", err)
			}
			return
		}

		log.Printf("New message received msgType: %d, content: %s", msgType, string(msg))

		ProcessData(c, msg)
	}
}

// write reads message from the channel and send it to the client.
func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error while writing message: %v", err)
				return
			}
		case <-ticker.C:
			log.Println("Send ping...")
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("meng")); err != nil {
				log.Printf("Error sending ping. Close the connect: %v", err)
				return
			}
		}
	}

}

// Start starts the read and write goroutine.
func (c *Client) Start() {
	go c.read()
	go c.write()
}

func (c *Client) GetKey() (key string) {
	key = c.UserID
	return
}

func (c *Client) Login(userId string, loginTime int64) {
	c.UserID = userId
	c.LoginTime = loginTime
}

func (c *Client) IsLogin() bool {
	if c.UserID != "" {
		return true
	}

	return false
}
