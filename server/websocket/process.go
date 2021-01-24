package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"githum.com/mengdage/gochat/common"
	"githum.com/mengdage/gochat/model"
)

// DisposeFunc is a handler function to handle websocket requests.
type DisposeFunc func(client *Client, message []byte) (code int32, msg string, data interface{})

var (
	handlersRWMutex sync.RWMutex
	handlers        = make(map[WsType]DisposeFunc)
)

// Register stores a handler by name
func Register(name WsType, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()

	handlers[name] = value
}

func getHandler(key WsType) (handler DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	handler, ok = handlers[key]
	return
}

// ProcessData handles a message from a client using handlers.
func ProcessData(client *Client, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[Process] Error: %s\n", r)
		}
	}()

	req := &model.WsRequest{}

	err := json.Unmarshal(message, req)
	if err != nil {
		fmt.Println("Failed to unmarshal request:", err)
		client.SendMessage([]byte("Invalid request"))
	}

	reqData, err := json.Marshal(req.Data)
	if err != nil {
		fmt.Println("Failed to marshal request data:", err)
		client.SendMessage([]byte("Invalid request data"))
	}

	cmd := WsType(req.Cmd)
	fmt.Printf("[ProcessData]%s: %s\n", client.Addr, cmd)

	var (
		code int32
		msg  string
		data interface{}
	)
	if handler, ok := getHandler(cmd); ok {
		code, msg, data = handler(client, reqData)
	} else {
		code = common.RoutingNotExist
		fmt.Printf("[ProcessData]Unknown CMD %s:%s\n", client.Addr, cmd)
	}

	resp := model.NewResponse(string(cmd), code, msg, data)

	respBs, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("[ProcessData]failed to marshall response")
	}

	client.SendMessage(respBs)
}
