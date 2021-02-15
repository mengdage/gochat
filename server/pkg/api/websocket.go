package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"githum.com/mengdage/gochat/helper"
)

type WSService interface {
	New(ctx context.Context, conn *websocket.Conn, clientManager *ClientManager, messageOperator *MessageOperator) error
}

type WSStorage interface {
	GetUserByName(context.Context, string) (*User, error)
	GetUserServer(ctx context.Context, userName string) (string, error)
	SaveUserServer(ctx context.Context, userName string, serverAddr string) error
}

type wsService struct {
	clientManager *ClientManager
	storage       WSStorage
}

func NewWSService(storage WSStorage, cm *ClientManager) WSService {
	return &wsService{
		storage:       storage,
		clientManager: cm,
	}
}

type LoginHelloBody struct {
	OK bool `json:"ok"`
}
type LoginHello struct {
	Cmd  string         `json:"cmd"`
	Body LoginHelloBody `json:"body"`
}

func (s *wsService) New(ctx context.Context, conn *websocket.Conn, clientManager *ClientManager, messageOperator *MessageOperator) error {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Failed to read the initial message: %v", err)
		return err
	}

	var loginReq WSLoginRequest
	err = json.Unmarshal(msg, &loginReq)
	if err != nil {
		log.Printf("Failed to parse the initial request: %v", err)
		return err
	}
	user, err := s.storage.GetUserByName(ctx, loginReq.Body.UserName)
	if err != nil {
		return err
	}

	client := NewClient(clientManager, user, conn, s.storage, messageOperator)
	if clientManager.Exists(user.Name) {
		return errors.New("already connected")
	}

	s.clientManager.Register <- client
	s.storage.SaveUserServer(ctx, user.Name, helper.GetServerIp())

	go client.Read()
	go client.Write()

	resp := LoginHello{
		Cmd: "login",
		Body: LoginHelloBody{
			OK: true,
		},
	}

	respBs, _ := json.Marshal(resp)
	client.Send <- respBs

	return nil
}
