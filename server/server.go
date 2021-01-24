package server

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/spf13/viper"
	"githum.com/mengdage/gochat/helper"
	"githum.com/mengdage/gochat/lib/cache"
	"githum.com/mengdage/gochat/model"
)

var (
	ServerIP string
	RPCPort  string
)

func initConfig() {
	log.Println("Initializaing server...")
	ServerIP = helper.GetServerIp()

	RPCPort = viper.GetString("app.rpc_port")
	log.Printf("Server %s:%s\n", ServerIP, RPCPort)
}

// Register adds the server to the cache.
func Register() (result bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[Register] stop %s\n%s\n", r, string(debug.Stack()))
		}
	}()
	initConfig()
	currentTime := time.Now().Unix()
	server := model.NewServer(ServerIP, RPCPort)
	cache.SetServerInfo(server, currentTime)
	fmt.Printf("Successfully register the server: %v\n", server)
	return true
}
