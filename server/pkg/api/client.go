package api

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"githum.com/mengdage/gochat/helper"
)

var (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type ClientStorage interface {
	SaveUserServer(ctx context.Context, userName string, serverAddr string) error
}
type Client struct {
	clientManager   *ClientManager
	messageOperator *MessageOperator
	UserID          int
	UserName        string
	Socket          *websocket.Conn
	Send            chan []byte
	storage         ClientStorage
}

func NewClient(clientManager *ClientManager, user *User, socket *websocket.Conn, storage ClientStorage, messageOperator *MessageOperator) *Client {
	return &Client{
		clientManager:   clientManager,
		messageOperator: messageOperator,
		UserID:          user.ID,
		UserName:        user.Name,
		Socket:          socket,
		Send:            make(chan []byte, 32),
		storage:         storage,
	}
}

func (c *Client) Read() {
	defer func() {
		log.Println("Closing send channel")
		close(c.Send)
		c.clientManager.Unregister <- c
	}()

	c.Socket.SetReadDeadline(time.Now().Add(pongWait))
	c.Socket.SetPongHandler(func(string) error {
		c.Socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unspected error: %v", err)
			}

			break
		}

		log.Printf("Rcv: %s", string(msg))
		if err := c.messageOperator.ProcessMessage(msg, c); err != nil {
			log.Printf("Error sending message: %v", err)

		}
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, nil)
				return
			}

			err := c.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error while writing message: %v", err)
				return
			}

		case <-ticker.C:
			log.Printf(`Send ping to %d`, c.UserID)
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping. Close the connection: %v", err)
				return
			}

			c.storage.SaveUserServer(context.Background(), c.UserName, helper.GetServerIp())
		}
	}
}
