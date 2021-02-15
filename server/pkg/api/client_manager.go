package api

import (
	"log"
)

type ClientManager struct {
	Register   chan *Client
	Unregister chan *Client
	clients    map[string]*Client
}

func NewClientManager() *ClientManager {

	return &ClientManager{
		clients:    make(map[string]*Client),
		Register:   make(chan *Client, 1024),
		Unregister: make(chan *Client, 1024),
	}
}

func (c *ClientManager) Exists(userName string) bool {
	_, exist := c.clients[userName]
	return exist
}

func (c *ClientManager) GetClient(userName string) *Client {
	client, exist := c.clients[userName]
	if !exist {
		return nil
	}
	return client
}

func (c *ClientManager) Start() {
	for {
		select {
		case client := <-c.Register:
			c.RegisterClient(client)
		case client := <-c.Unregister:
			c.UnregisterClient(client)

		}

	}
}

func (c *ClientManager) sendMessageTo(userName string, msg []byte) {
	client, exist := c.clients[userName]
	if !exist {
		log.Printf("User %s not connected", userName)
		return
	}

	log.Printf("Send message to %s", userName)

	client.Send <- msg
}

// func (c *ClientManager) ProcessMessage(msg []byte, fromClient *Client) error {
// 	wsMsg := WSMessage{}
// 	if err := json.Unmarshal(msg, &wsMsg); err != nil {
// 		return err
// 	}

// 	if wsMsg.Cmd == "send" {
// 		body := WSSendBody{}
// 		if err := json.Unmarshal(wsMsg.Body, &body); err != nil {
// 			return err
// 		}
// 		recvBody := WSRecvBody{
// 			FromUserName: fromClient.UserName,
// 			ToUserName:   body.UserName,
// 			Content:      body.Content,
// 		}
// 		recvBodyBs, _ := json.Marshal(recvBody)

// 		recvMsg := WSMessage{
// 			Cmd:  "recv",
// 			Body: recvBodyBs,
// 		}

// 		recvBs, _ := json.Marshal(recvMsg)

// 		c.sendMessageTo(body.UserName, recvBs)

// 	}

// 	return nil
// }

func (c *ClientManager) RegisterClient(client *Client) {
	log.Printf("Registering user: %d %s", client.UserID, client.UserName)
	c.clients[client.UserName] = client
}

func (c *ClientManager) UnregisterClient(client *Client) {
	log.Printf("Unregistering user: %d %s", client.UserID, client.UserName)
	if _, ok := c.clients[client.UserName]; ok {
		delete(c.clients, client.UserName)
	}
}
