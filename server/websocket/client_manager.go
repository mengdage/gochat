package websocket

import (
	"fmt"
	"log"
	"sync"
)

type ClientManager struct {
	Clients     map[*Client]bool
	ClientsLock sync.RWMutex
	Users       map[string]*Client
	UsersLock   sync.RWMutex
	Register    chan *Client
	Login       chan *UserLogin
	Unregister  chan *Client
	Broadcast   chan []byte
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		Login:      make(chan *UserLogin, 1000),
		Unregister: make(chan *Client, 100),
		Broadcast:  make(chan []byte, 1000),
	}

}

func (m *ClientManager) InClient(client *Client) (ok bool) {
	m.ClientsLock.RLock()
	defer m.ClientsLock.RUnlock()

	_, ok = m.Clients[client]

	return
}

func (m *ClientManager) GetClients() (clients []*Client) {
	m.ClientsLock.RLock()
	defer m.ClientsLock.RUnlock()

	clients = make([]*Client, 0)
	for c := range m.Clients {
		clients = append(clients, c)
	}

	return
}

func (m *ClientManager) GetClientsLen() int {
	return len(m.Clients)
}

func (m *ClientManager) AddClient(c *Client) {
	m.ClientsLock.Lock()
	defer m.ClientsLock.Unlock()

	m.Clients[c] = true
}

func (m *ClientManager) DelClient(c *Client) {
	m.ClientsLock.Lock()
	defer m.ClientsLock.Unlock()

	if _, ok := m.Clients[c]; ok {
		delete(m.Clients, c)
	}
}

func (m *ClientManager) GetUserClient(userId string) *Client {
	m.UsersLock.RLock()
	defer m.UsersLock.RUnlock()

	if c, ok := m.Users[userId]; ok {
		return c
	}

	return nil
}

func (m *ClientManager) GetUserLen() int {
	return len(m.Users)
}

func (m *ClientManager) AddUser(userId string, client *Client) {
	m.UsersLock.Lock()
	defer m.UsersLock.Unlock()

	m.Users[userId] = client
}

func (m *ClientManager) DelUsers(client *Client) {
	m.UsersLock.Lock()
	defer m.UsersLock.Unlock()

	if u, ok := m.Users[client.UserID]; ok {
		if u.Addr != client.Addr {
			return
		}

		delete(m.Users, client.UserID)
	}
}

// GetUserIds returns a list of all users' ids
func (m *ClientManager) GetUserIds() []string {
	userKeys := []string{}
	m.UsersLock.RLock()
	defer m.UsersLock.RUnlock()

	for userId, _ := range m.Users {
		userKeys = append(userKeys, userId)
	}

	return userKeys
}

func (m *ClientManager) GetUserClients() []*Client {
	clients := make([]*Client, 0)

	m.ClientsLock.RLock()
	defer m.ClientsLock.RUnlock()

	for c, _ := range m.Clients {
		clients = append(clients, c)
	}

	return clients
}

func (m *ClientManager) RegisterClient(client *Client) {
	m.AddClient(client)
}

func (m *ClientManager) UnRegisterClient(client *Client) {
	log.Printf("Unregister client %s (%s)\n", client.Addr, client.UserID)
	m.DelUsers(client)
	m.DelClient(client)
}

// LoginUser registers a user.
func (m *ClientManager) LoginUser(login *UserLogin) {
	m.AddUser(login.UserID, login.Client)
}

func (m *ClientManager) SendMsgToAll(msg []byte) {
	clients := m.GetClients()
	for _, c := range clients {
		c.Send <- msg
	}
}

func (m *ClientManager) Start() {
	fmt.Println("Start ClientManage")
	for {
		select {
		case client := <-m.Register:
			m.RegisterClient(client)
		case login := <-m.Login:
			m.LoginUser(login)
		case client := <-m.Unregister:
			m.UnRegisterClient(client)
		case bs := <-m.Broadcast:
			m.SendMsgToAll(bs)
		}
	}
}

func GetAllUsers() []string {
	return clientManager.GetUserIds()
}
