package model

const (
	heartbeatInterval = 3 * 60
)

type UserOnline struct {
	ServerIp      string `json:"serverIp"`
	RpcPort       string `json:"rpcPort"`
	UserID        string `json:"userId"`
	ClientIp      string `json:"clientIp"`
	ClientPort    string `json:"clientPort"`
	LoginTime     int64  `json:"loginTime"`
	HeartbeatTime int64  `json:"heartbeatTime"`
	LogoutTime    int64  `json:"logoutTime"`
	IsLogoff      bool   `json:"isLogoff"`
}

func NewUserOnline(serverIP, rpcPort, userId, addr string, loginTime int64) *UserOnline {
	uo := &UserOnline{
		ServerIp:      serverIP,
		RpcPort:       rpcPort,
		UserID:        userId,
		ClientIp:      addr,
		LoginTime:     loginTime,
		HeartbeatTime: loginTime,
		IsLogoff:      false,
	}

	return uo
}
