package websocket

var (
	clientManager = NewClientManager()
)

func init() {
	go clientManager.Start()
	registerWsHandler()
}

func registerWsHandler() {
	Register(LoginType, LoginHandler)
	Register(SendType, SendHandler)
}
