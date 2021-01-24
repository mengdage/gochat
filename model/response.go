package model

// Response represents a response to a client.
type Response struct {
	Cmd  string        `json:"cmd"`
	Data *ResponseData `json:"responseData"`
}

// ResponseData represents a response data.
type ResponseData struct {
	Code    int32       `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	Content interface{} `json:"content"`
}

// NewResponse creates a new response.
func NewResponse(cmd string, code int32, codeMsg string, content interface{}) *Response {
	rd := &ResponseData{code, codeMsg, content}

	return &Response{cmd, rd}
}
