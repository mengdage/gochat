package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"githum.com/mengdage/gochat/common"
	"githum.com/mengdage/gochat/lib/cache"
	"githum.com/mengdage/gochat/model"
	"githum.com/mengdage/gochat/server"
)

// WsType represents type of a ws message
type WsType string

// Common type of ws message
const (
	LoginType WsType = "login"
	SendType  WsType = "send"
)

const (
	connectWait = 10 * time.Second
)

// WsConnect upgrade a http request to webscoket
func WsConnect(rw http.ResponseWriter, req *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header["Origin"]
			if len(origin) == 0 {
				return true
			}
			u, err := url.Parse(origin[0])
			if err != nil {
				return false
			}
			// return equalASCIIFold(u.Host, r.Host)
			fmt.Printf("origin: %s; host: %s\n", u.Host, r.Host)
			return true
		},
	}

	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}

	fmt.Printf("websocket trying to connect: %s\n", conn.RemoteAddr())

	currentTime := time.Now().Unix()

	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)

	client.Start()

	clientManager.Register <- client
}

// LoginHandler handles a websocket login request
func LoginHandler(client *Client, message []byte) (code int32, msg string, data interface{}) {
	code = common.OK
	currentTime := time.Now().Unix()

	request := &model.LoginData{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.IllegalParameter
		fmt.Printf("Failed to unmarshal message: %s", message)
		return
	}

	log.Printf("[Login] User %s\n", request.UserID)

	if request.UserID == "" || len(request.UserID) >= 20 {
		code = common.UnauthorizedUserId
		msg = fmt.Sprintf("Illegal user %s\n", request.UserID)
		log.Println(msg)
		return
	}

	if client.IsLogin() {
		msg = fmt.Sprintf("User %s has already logged in", client.UserID)
		code = common.OperationFailure
		log.Println(code, msg)

		return
	}

	client.Login(request.UserID, currentTime)

	login := &UserLogin{
		UserID: request.UserID,
		Client: client,
	}
	clientManager.Login <- login

	userOnline := model.NewUserOnline(server.ServerIP, server.RPCPort, request.UserID, client.Addr, currentTime)
	err := cache.SetUserOnlineInfo(userOnline)
	if err != nil {
		code = common.ServerError
		msg = fmt.Sprintf("Failed to set user online info: %s", err.Error())
		log.Println(code, msg)
		return
	}

	fmt.Printf("[Login] Log in successfully addr %s, user %s\n", client.Addr, request.UserID)
	data = map[string]string{
		"userId": request.UserID,
	}

	return

}

// SendMsg is a request for sending a message.
type SendMsg struct {
	UserID  string `json:"userId"`
	Content string `json:"content"`
}

// SendHandler handles a request of sending a message.
func SendHandler(client *Client, message []byte) (code int32, msg string, data interface{}) {
	sendMsg := &SendMsg{}
	if err := json.Unmarshal(message, sendMsg); err != nil {
		return http.StatusBadRequest, "Invalid message", nil
	}

	log.Printf("Send to user %s: %s", sendMsg.UserID, sendMsg.Content)

	go SendUserMessage(client, sendMsg.UserID, sendMsg.Content)

	return http.StatusOK, "Received", nil
}
