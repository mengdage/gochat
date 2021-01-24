package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"githum.com/mengdage/gochat/lib/redislib"
	"githum.com/mengdage/gochat/route"
	"githum.com/mengdage/gochat/server"
	"githum.com/mengdage/gochat/server/grpc_server"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func init() {
	initConfig()
	initLog()
	initRedis()

}
func main() {
	flag.Parse()
	server.Register()

	router := gin.Default()
	router.Use(cors.Default())
	route.Init(router)

	go grpc_server.Start()
	httpAddr := "0.0.0.0:" + viper.GetString("app.http_port")
	fmt.Printf("Listening on %s\n", httpAddr)
	http.ListenAndServe(httpAddr, router)

}

func initConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	fmt.Println("app configurations:")
	for k, v := range viper.GetStringMapString("app") {
		fmt.Printf("\t%s: %s\n", k, v)
	}

}

func initLog() {
	gin.DisableConsoleColor()

	logFilePath := viper.GetString("app.log_file")
	f, err := os.Create(logFilePath)

	if err != nil {
		panic(fmt.Errorf("Failed to create logFile %s: %s", logFilePath, err.Error()))
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func initRedis() {
	redislib.InitClient()
}
