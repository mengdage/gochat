package api

import (
	"encoding/json"
)

type NewUserRequest struct {
	Name string `json:"name"`
}

type LoginRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WSLoginBody struct {
	UserName string `json:"userName"`
}
type WSLoginRequest struct {
	Cmd  string      `json:"cmd"`
	Body WSLoginBody `json:"body"`
}

type WSMessage struct {
	Cmd  string          `json:"cmd"`
	Body json.RawMessage `json:"body"`
}

type WSSendBody struct {
	UserName string `json:"userName"`
	Content  string `json:"content"`
}
