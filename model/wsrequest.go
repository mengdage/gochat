package model

// WsRequest represents a webscoket request object from client.
type WsRequest struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data,omitempty"`
}

// LoginData represents the data of a login request.
type LoginData struct {
	UserID string `json:"userId,omitempty"`
}
